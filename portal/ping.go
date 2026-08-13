package portal

import (
	"context"
	"fmt"
)

func (c *Client) Ping() error {
	hReq, err := c.createGetRequest(context.Background(), "/api/v1/ping", nil)
	if err != nil {
		return fmt.Errorf("failed to create ping request: %w", err)
	}

	hRes, err := c.httpClient.Do(hReq)
	if err != nil {
		return fmt.Errorf("failed to send ping request: %w", err)
	}
	defer hRes.Body.Close()

	err = expect(hRes, 200)
	if err != nil {
		return fmt.Errorf("failed to ping portal: %w", err)
	}

	return nil
}
