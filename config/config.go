// config.go
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Fenlei, Fenxi string

type Config struct {
	APIKey          string
	AccessKey       string
	UserName        string
	TopicBaseURL    string
	TopicListURL    string
	PageDataBaseURL string
	Es              *Elasticsearch
}

type Elasticsearch struct {
	Address  string
	Username string
	Password string
}

func init() {
	// 加载 .env 文件
	err := godotenv.Load("/Users/songyuhang/study/oncall/.env")
	if err != nil {
		logrus.WithError(err).Error("Error loading .env file")
	}
	Fenlei = os.Getenv("FENLEI_TEMPLATE")
	Fenxi = os.Getenv("FENXI_TEMPLATE")
}

func NewConfig() (*Config, error) {
	apiKey := os.Getenv("API_KEY")
	accessKey := os.Getenv("ACCESS_KEY")
	userName := os.Getenv("USER_NAME")
	topicBaseURL := os.Getenv("TOPIC_BASE_URL")
	topicListURL := os.Getenv("TOPIC_LIST_URL")
	pageDataBaseURL := os.Getenv("PAGE_DATA_BASE_URL")
	esAddress := os.Getenv("ES_ADDRESS")
	esUsername := os.Getenv("ES_USERNAME")
	esPassword := os.Getenv("ES_PASSWORD")

	// 校验关键配置是否缺失
	if apiKey == "" || accessKey == "" || userName == "" || topicBaseURL == "" || topicListURL == "" || pageDataBaseURL == "" {
		return nil, errors.New("missing required environment variables (API_KEY, ACCESS_KEY, USER_NAME, TOPIC_BASE_URL, TOPIC_LIST_URL)")
	}

	if esAddress == "" || esUsername == "" || esPassword == "" {
		return nil, errors.New("missing required environment variables (ES_ADDRESS, ES_USERNAME, ES_PASSWORD)")
	}

	return &Config{
		APIKey:          apiKey,
		AccessKey:       accessKey,
		UserName:        userName,
		TopicBaseURL:    topicBaseURL,
		TopicListURL:    topicListURL,
		PageDataBaseURL: pageDataBaseURL,
		Es: &Elasticsearch{
			Address:  esAddress,
			Username: esUsername,
			Password: esPassword,
		},
	}, nil
}
