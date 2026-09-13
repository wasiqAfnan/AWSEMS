package apigateway

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"awsems/internal/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API Gateway client. So we can make requests to the API Gateway.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(cfg.APIGatewayBaseURL, "/"),
		HTTPClient: &http.Client{},
	}
}

// Get sends a GET request to the given path on the API Gateway base URL.
// It returns the raw response body, the HTTP status code, and any error.
func (c *Client) Get(path string) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("request to API Gateway failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, nil
}
