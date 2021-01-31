package api

import (
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
	Host       string
	Scheme     string
}

func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		Host:       "api.pathofexile.com",
		Scheme:     "http",
	}
}

func (c *Client) CallAPI(path string, query string) (*http.Response, error) {
	callURL := url.URL{
		Scheme:   c.Scheme,
		Host:     c.Host,
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
