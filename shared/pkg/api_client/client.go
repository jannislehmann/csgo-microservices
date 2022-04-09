package api_client

import "net/http"

type HttpApiClient struct{}

func (s *HttpApiClient) Get(u string) (*http.Response, error) {
	return http.Get(u)
}
