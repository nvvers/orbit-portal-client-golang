package portal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient    *http.Client
	portalBaseUrl string
	portalToken   string
}

// Option is a function that configures the Client. It can be used to set optional parameters such as authentication tokens.
type Option func(c *Client) error

// NewClient creates a new Client with the given portalBaseUrl and applies any provided options.
// It returns an error if any of the options fail to apply.
func NewClient(portalBaseUrl string, opts ...Option) (*Client, error) {
	c := &Client{
		httpClient:    &http.Client{},
		portalBaseUrl: portalBaseUrl,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithToken(token string) Option {
	return func(c *Client) error {
		c.portalToken = token
		return nil
	}
}

func WithHttpClient(hClient *http.Client) Option {
	return func(c *Client) error {
		if hClient == nil {
			return errors.New("http client cannot be nil")
		}

		c.httpClient = hClient
		return nil
	}
}

func (c *Client) createGetRequest(ctx context.Context, portalApiEndpoint, query string) (*http.Request, error) {
	apiEndpoint, err := url.JoinPath(c.portalBaseUrl, portalApiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct API endpoint: %w", err)
	}

	if query != "" {
		apiEndpoint = apiEndpoint + "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if c.portalToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.portalToken)
	}

	return httpReq, nil
}

func (c *Client) createPostRequest(ctx context.Context, portalApiEndpoint, query string, payload any) (*http.Request, error) {
	apiEndpoint, err := url.JoinPath(c.portalBaseUrl, portalApiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct API endpoint: %w", err)
	}

	if query != "" {
		apiEndpoint = apiEndpoint + "?" + query
	}

	reqData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if c.portalToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.portalToken)
	}

	return httpReq, nil
}

func expect(resp *http.Response, expectedStatusCode int) error {
	if resp.StatusCode != expectedStatusCode {
		errText, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to call API, code: %d expected: %d, response: %s", resp.StatusCode, expectedStatusCode, errText)
	}

	return nil
}
