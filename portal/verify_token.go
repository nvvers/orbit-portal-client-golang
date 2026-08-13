package portal

import (
	"context"
	"fmt"
)

func (c *Client) VerifyToken() error {
	hReq, err := c.createGetRequest(context.Background(), "/api/v1/verify-token", nil)
	if err != nil {
		return fmt.Errorf("failed to create verify-token request: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to send verify-token request: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, 200)
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	return nil
}
