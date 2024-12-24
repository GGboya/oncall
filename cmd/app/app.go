package app

import (
	"oncall/config"
	"oncall/internal/apiservice"
	"oncall/internal/httpclient"

	"github.com/elastic/go-elasticsearch/v8"
)

type AppDependencies struct {
	Config     *config.Config
	HTTPClient *httpclient.HTTPClient
	APIService *apiservice.APIService
	ESClient   *elasticsearch.Client
}

func NewAppDependencies() (*AppDependencies, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	httpClient := httpclient.NewHTTPClient()
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			cfg.Es.Address, // 配置 Elasticsearch 地址
		},
		Username: cfg.Es.Username,
		Password: cfg.Es.Password,
	})
	apiService := apiservice.NewAPIService(httpClient, cfg, es)
	return &AppDependencies{
		Config:     cfg,
		HTTPClient: &httpClient,
		APIService: apiService,
	}, nil
}
