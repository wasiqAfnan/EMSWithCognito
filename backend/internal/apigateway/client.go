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
func (c *Client) Get(path string, payload []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var req *http.Request
	var err error
	if len(payload) > 0 {
		req, err = http.NewRequest(http.MethodGet, url, strings.NewReader(string(payload)))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to create GET request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
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

// Post sends a POST request to the given path with the provided JSON payload.
// It returns the raw response body, the HTTP status code, and any error.
func (c *Client) Post(path string, payload []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
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

// Patch sends a PATCH request to the given path with the provided JSON payload.
// It returns the raw response body, the HTTP status code, and any error.
func (c *Client) Patch(path string, payload []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	req, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create PATCH request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
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

// Delete sends a DELETE request to the given path.
// It returns the raw response body, the HTTP status code, and any error.
func (c *Client) Delete(path string, payload []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var req *http.Request
	var err error
	if len(payload) > 0 {
		req, err = http.NewRequest(http.MethodDelete, url, strings.NewReader(string(payload)))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest(http.MethodDelete, url, nil)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
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
