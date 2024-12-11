package client

import (
	"net"
	"net/http"
	"time"
)

// Client is a reusable HTTP client for the SDK.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthHeader string
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient is an option for setting a custom http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithAuthHeader is an option for setting an authentication header.
func WithAuthHeader(authHeader string) Option {
	return func(c *Client) {
		c.AuthHeader = authHeader
	}
}

// defaultHTTPClient provides a preconfigured HTTP client for the SDK.
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}
}

func NewClient(baseURL string, options ...Option) *Client {
	client := &Client{
		BaseURL:    baseURL,
		HTTPClient: defaultHTTPClient(),
	}

	// Apply functional options
	for _, opt := range options {
		opt(client)
	}

	return client
}

// Do perform an HTTP request with optional authentication.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.AuthHeader != "" {
		req.Header.Set("Authorization", c.AuthHeader)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.HTTPClient.Do(req)
}
