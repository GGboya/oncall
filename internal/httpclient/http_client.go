package httpclient

import (
	"net/http"
)

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type httpClient struct{}

func (c *httpClient) Do(request *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(request)
}

func NewHTTPClient() HTTPClient {
	return &httpClient{}
}
