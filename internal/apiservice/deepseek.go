package apiservice

import (
	"encoding/json"
	"fmt"
	"strings"

	"oncall/config"
	"oncall/internal/httpclient"
	"oncall/model"
)

type DeepSeekService struct {
	client httpclient.HTTPClient
	config *config.Config
}

func NewDeepSeekService(client httpclient.HTTPClient, config *config.Config) *DeepSeekService {
	return &DeepSeekService{
		client: client,
		config: config,
	}
}

func (s *DeepSeekService) SendDeepSeekRequest(data []string, action string) (*model.DeepSeekResponse, error) {
	url := "https://api.deepseek.com/chat/completions"

	// 构造请求体
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "system", "content": action},
			{"role": "user", "content": strings.Join(data, "\n")},
		},
		"stream": false,
	}

	respBody, err := doRequest(s.client, "POST", url, requestBody, *s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to DeepSeek: %v", err)
	}

	var deepSeekResponse model.DeepSeekResponse
	if err := json.Unmarshal(respBody, &deepSeekResponse); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}

	return &deepSeekResponse, nil
}
