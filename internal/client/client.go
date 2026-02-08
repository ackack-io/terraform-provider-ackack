// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://api.ackack.io"
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
	retryBaseDelay = time.Second
)

// Client is the ackack.io API client.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	UserAgent  string
}

// NewClient creates a new ackack.io API client.
func NewClient(apiKey, endpoint, version string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	baseURL := endpoint
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	userAgent := "terraform-provider-ackack"
	if version != "" {
		userAgent = fmt.Sprintf("terraform-provider-ackack/%s", version)
	}

	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
		UserAgent: userAgent,
	}, nil
}

// doRequest performs an HTTP request with retries and error handling.
func (c *Client) doRequest(ctx context.Context, method, path string, body, result any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	var lastErr error
	for attempt := range maxRetries {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryBaseDelay * time.Duration(attempt)):
			}
		}

		// Reset body reader if retrying
		if body != nil && attempt > 0 {
			jsonBody, _ := json.Marshal(body)
			bodyReader = bytes.NewReader(jsonBody)
		}

		req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.UserAgent)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		// Handle rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := 60 // default to 60 seconds
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if parsed, err := strconv.Atoi(ra); err == nil {
					retryAfter = parsed
				}
			}
			lastErr = &APIError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("rate limited, retry after %d seconds", retryAfter),
			}
			// Wait for the retry-after duration
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(retryAfter) * time.Second):
			}
			continue
		}

		// Handle successful responses
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if result != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, result); err != nil {
					return fmt.Errorf("failed to unmarshal response: %w", err)
				}
			}
			return nil
		}

		// Handle error responses
		var errorResp ErrorResponse
		if len(respBody) > 0 {
			_ = json.Unmarshal(respBody, &errorResp)
		}

		message := errorResp.Message
		if message == "" {
			message = errorResp.Error
		}
		if message == "" {
			message = http.StatusText(resp.StatusCode)
		}

		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Message:    message,
			ErrorField: errorResp.Error,
		}

		// Don't retry client errors (except rate limiting which is handled above)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return apiErr
		}

		lastErr = apiErr
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("max retries exceeded")
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, result any) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body, result any) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

// put performs a PUT request.
func (c *Client) put(ctx context.Context, path string, body, result any) error {
	return c.doRequest(ctx, http.MethodPut, path, body, result)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
