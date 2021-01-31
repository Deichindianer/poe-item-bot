package api

import (
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
}

func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) callAPI(host string, path string, query string) (*http.Response, error) {
	callURL := url.URL{
		Scheme:   "http",
		Host:     host,
		Path:     path,
		RawQuery: query,
	}
	req, err := http.NewRequest(http.MethodGet, callURL.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
