package app

import (
	"oncall/config"
	"oncall/internal/apiservice"
	"oncall/internal/httpclient"
)

type AppDependencies struct {
	Config     *config.Config
	HTTPClient *httpclient.HTTPClient
	APIService *apiservice.APIService
}

func NewAppDependencies() (*AppDependencies, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	httpClient := httpclient.NewHTTPClient()
	apiService := apiservice.NewAPIService(httpClient, cfg)

	return &AppDependencies{
		Config:     cfg,
		HTTPClient: &httpClient,
		APIService: apiService,
	}, nil
}
