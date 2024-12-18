package apiservice

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"oncall/config"
	"oncall/internal/httpclient"
	"oncall/model"
)

type OncallService struct {
	client httpclient.HTTPClient
	config *config.Config
}

func NewOncallService(client httpclient.HTTPClient, config *config.Config) *OncallService {
	return &OncallService{
		client: client,
		config: config,
	}
}

func (s *OncallService) FetchAllData(kind string) ([]string, error) {
	topics, err := s.getTopics(kind)
	if err != nil {
		return nil, err
	}
	resp := make([]string, len(topics))

	// 现在拿到所有的 topic 了，根据 day 还是 week，进行不同的大模型处理
	var wg sync.WaitGroup
	wg.Add(len(topics))
	var errMsg error
	for idx, topic := range topics {
		go func(idx int, topic model.Topic) {
			if err := func() error {
				defer wg.Done()
				str := fmt.Sprintf("%s, %s/%d\n", topic.Title, s.config.TopicListURL, topic.ID)
				if kind == "day" {
					// 还需要额外获取详细信息
					topicDetail, err := s.getTopicDetails(kind, topic.ID)
					if err != nil {
						return err
					}
					str += topicDetail
				}
				resp[idx] = str
				return nil
			}(); err != nil {
				errMsg = fmt.Errorf("%v; %w", errMsg, err)
			}

		}(idx, topic)
	}
	wg.Wait()

	if errMsg != nil {
		return nil, errMsg
	}

	return resp, nil
}

func (s *OncallService) getTopics(kind string) ([]model.Topic, error) {
	timeWindow := map[string]time.Duration{"day": 24 * time.Hour, "week": 7 * 24 * time.Hour}[kind]
	now := time.Now()
	topics := []model.Topic{}

	for page := 0; page < 100; page++ {
		responseData, err := s.FetchPageData(page)
		if err != nil {
			return nil, err
		}
		for idx, topic := range responseData.TopicList.Topics {
			if page == 0 && idx <= 2 {
				continue
			}
			//  先拿到所有的范围内的 topic
			if now.Sub(topic.CreatedAt) > timeWindow && now.Sub(topic.LastPostedAt) > timeWindow {
				return topics, nil
			}
			topics = append(topics, topic)
		}
	}
	return topics, nil
}

// getTopicDetails 根据需要获取 Topic 详细信息
func (s *OncallService) getTopicDetails(kind string, topicID int) (string, error) {
	url := fmt.Sprintf("%s/%d/", s.config.TopicBaseURL, topicID) + "1.json?track_visit=true&forceLoad=true"
	data, err := s.FetchTopicDetail(url)
	if err != nil {
		return "", err
	}
	postArr := []string{}
	for _, post := range data.PostStream.Posts {
		str, _ := decodeUnicode(post.Cooked)
		postArr = append(postArr, post.Name+": "+str)
	}
	str := strings.Join(postArr, "\n")
	return str, nil
}

func (s *OncallService) FetchPageData(page int) (*model.Response, error) {
	url := os.Getenv("PAGE_DATA_BASE_URL") + strconv.Itoa(page)
	respBody, err := doRequest(s.client, "GET", url, nil, *s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page data: %w", err)
	}

	var responseData model.Response
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse page data: %w", err)
	}

	return &responseData, nil
}

func (s *OncallService) FetchTopicDetail(url string) (*model.PostResponse, error) {
	respBody, err := doRequest(s.client, "GET", url, nil, *s.config)
	if err != nil {
		return nil, err
	}
	var responseData model.PostResponse
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}
