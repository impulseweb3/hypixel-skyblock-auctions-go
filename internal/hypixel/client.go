package hypixel

import "net/http"

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}
