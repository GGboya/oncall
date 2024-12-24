package apiservice

import (
	"oncall/config"
	"oncall/internal/httpclient"

	"github.com/elastic/go-elasticsearch/v8"
)

type APIService struct {
	DeepSeek *DeepSeekService
	Oncall   *OncallService
}

func NewAPIService(client httpclient.HTTPClient, config *config.Config, es *elasticsearch.Client) *APIService {
	return &APIService{
		DeepSeek: NewDeepSeekService(client, config),
		Oncall:   NewOncallService(client, config, es),
	}
}
