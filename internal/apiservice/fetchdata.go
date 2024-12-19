package apiservice

import (
	"encoding/json"
	"fmt"
	"oncall/config"
	"oncall/internal/httpclient"
	"oncall/model"
	"strings"
	"sync"
)

type FetchDataStrategy interface {
	FetchData(topics []model.Topic, config *config.Config) ([]string, error)
}

type DayFetchDataStrategy struct {
	client httpclient.HTTPClient
}

func NewDayFetchDataStrategy(client httpclient.HTTPClient) *DayFetchDataStrategy {
	return &DayFetchDataStrategy{
		client: client,
	}
}

func (d *DayFetchDataStrategy) FetchData(topics []model.Topic, config *config.Config) ([]string, error) {
	// day 相关逻辑
	// 每天的数据，还需要额外把每个帖子中的讨论也进行整合，放入到大模型中处理
	resp := make([]string, len(topics))
	fmt.Println("day fetch data", topics)

	// 现在拿到所有的 topic 了，根据 day 还是 week，进行不同的大模型处理
	var wg sync.WaitGroup
	wg.Add(len(topics))
	var errMsg error
	for idx, topic := range topics {
		go func(idx int, topic model.Topic) {
			if err := func() error {
				defer wg.Done()
				str := fmt.Sprintf("%s, %s/%d\n", topic.Title, config.TopicListURL, topic.ID)
				url := fmt.Sprintf("%s/%d/", config.TopicBaseURL, topic.ID) + "1.json?track_visit=true&forceLoad=true"
				respBody, err := doRequest(d.client, "GET", url, nil, config)
				if err != nil {
					return err
				}
				var responseData model.PostResponse
				err = json.Unmarshal(respBody, &responseData)
				if err != nil {
					return err
				}

				postArr := []string{}
				for _, post := range responseData.PostStream.Posts {
					str, _ := decodeUnicode(post.Cooked)
					postArr = append(postArr, post.Name+": "+str)
				}
				str += strings.Join(postArr, "\n")
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

type WeekFetchDataStrategy struct {
	client httpclient.HTTPClient
}

func NewWeekFetchDataStrategy(client httpclient.HTTPClient) *DayFetchDataStrategy {
	return &DayFetchDataStrategy{
		client: client,
	}
}

func (w *WeekFetchDataStrategy) FetchData(topics []model.Topic, config *config.Config) ([]string, error) {
	// week 相关逻辑
	resp := make([]string, len(topics))

	for idx, topic := range topics {
		str := fmt.Sprintf("%s, %s/%d\n", topic.Title, config.TopicListURL, topic.ID)
		resp[idx] = str

	}
	return resp, nil
}
