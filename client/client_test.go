package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mollie/go-apicurio-registry/client"
	"github.com/stretchr/testify/assert"
)

func TestNewClient_Defaults(t *testing.T) {
	c := client.NewClient("https://example.com")

	assert.Equal(t, "https://example.com", c.BaseURL)
	assert.NotNil(t, c.HTTPClient)
	assert.Equal(t, 30*time.Second, c.HTTPClient.Timeout)
}

func TestNewClient_WithCustomHTTPClient(t *testing.T) {
	customHTTPClient := &http.Client{Timeout: 10 * time.Second}

	c := client.NewClient("https://example.com", client.WithHTTPClient(customHTTPClient))

	assert.Equal(t, "https://example.com", c.BaseURL)
	assert.Equal(t, customHTTPClient, c.HTTPClient)
}

func TestNewClient_WithAuthHeader(t *testing.T) {
	authHeader := "Bearer test-token"

	c := client.NewClient("https://example.com", client.WithAuthHeader(authHeader))

	assert.Equal(t, "Bearer test-token", c.AuthHeader)
}

func TestClient_Do_WithAuthHeader(t *testing.T) {
	// Create a test HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	c := client.NewClient(server.URL, client.WithAuthHeader("Bearer test-token"))

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.NoError(t, err)

	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClient_Do_WithoutAuthHeader(t *testing.T) {
	// Create a test HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	c := client.NewClient(server.URL)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.NoError(t, err)

	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
