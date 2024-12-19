package apiservice

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"oncall/config"
	"oncall/internal/httpclient"
	"oncall/model"
)

type OncallService struct {
	client   httpclient.HTTPClient
	config   *config.Config
	strategy FetchDataStrategy
}

func NewOncallService(client httpclient.HTTPClient, config *config.Config) *OncallService {
	return &OncallService{
		client: client,
		config: config,
	}
}

func (s *OncallService) SetFetchDataStrategy(kind string) {
	switch kind {
	case "day":
		s.strategy = &DayFetchDataStrategy{s.client}
	case "week":
		s.strategy = &WeekFetchDataStrategy{s.client}
	default:
		s.strategy = &DayFetchDataStrategy{s.client} // 或者设置一个默认策略
	}
}

func (s *OncallService) FetchAllData(kind string) ([]string, error) {
	topics, err := s.getTopics(kind)
	if err != nil {
		return nil, err
	}
	s.SetFetchDataStrategy(kind)
	return s.strategy.FetchData(topics, s.config)
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
	respBody, err := doRequest(s.client, "GET", url, nil, s.config)
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
	respBody, err := doRequest(s.client, "GET", url, nil, s.config)
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
