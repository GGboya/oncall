package apiservice

import (
	"oncall/config"
	"oncall/internal/httpclient"
)

type APIService struct {
	DeepSeek *DeepSeekService
	Oncall   *OncallService
}

func NewAPIService(client httpclient.HTTPClient, config *config.Config) *APIService {
	return &APIService{
		DeepSeek: NewDeepSeekService(client, config),
		Oncall:   NewOncallService(client, config),
	}
}
