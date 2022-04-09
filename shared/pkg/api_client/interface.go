package api_client

import "net/http"

type ApiClient interface {
	Get(url string) (resp *http.Response, err error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
