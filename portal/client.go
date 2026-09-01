package portal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient    *http.Client
	portalBaseUrl *url.URL
	portalToken   string
	defaultSource string
}

// Option is a function that configures the Client. It can be used to set optional parameters such as authentication tokens.
type Option func(c *Client) error

// NewClient creates a new Client with the given portalBaseUrl and applies any provided options.
// It returns an error if any of the options fail to apply.
func NewClient(portalBaseUrl string, opts ...Option) (*Client, error) {
	baseUrl, err := url.Parse(portalBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse portal base URL: %w", err)
	}

	c := &Client{
		httpClient:    &http.Client{},
		portalBaseUrl: baseUrl,
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
		if c.portalBaseUrl.Scheme != "https" {
			slog.Warn("Using a token with a non-HTTPS URL is not secure. Consider using HTTPS for secure communication.")
		}

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

func WithDefaultSource(source string) Option {
	return func(c *Client) error {
		c.defaultSource = source
		return nil
	}
}

func (c *Client) createGetRequest(ctx context.Context, portalApiEndpoint string, query url.Values) (*http.Request, error) {
	apiEndpoint := c.portalBaseUrl.JoinPath(portalApiEndpoint)
	apiEndpoint.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiEndpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if c.portalToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.portalToken)
	}

	return httpReq, nil
}

func (c *Client) createPostRequest(ctx context.Context, portalApiEndpoint string, query url.Values, payload io.Reader) (*http.Request, error) {
	apiEndpoint := c.portalBaseUrl.JoinPath(portalApiEndpoint)
	apiEndpoint.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/octet-stream")

	if c.portalToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.portalToken)
	}

	return httpReq, nil
}

func (c *Client) createPostRequestWithJsonBody(ctx context.Context, portalApiEndpoint string, query url.Values, payload any) (*http.Request, error) {
	apiEndpoint := c.portalBaseUrl.JoinPath(portalApiEndpoint)
	apiEndpoint.RawQuery = query.Encode()

	reqData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint.String(), bytes.NewBuffer(reqData))
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
		errText, _ := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024)) // limit to 1MB
		return fmt.Errorf("failed to call API, code: %d expected: %d, response: %s", resp.StatusCode, expectedStatusCode, errText)
	}

	return nil
}
