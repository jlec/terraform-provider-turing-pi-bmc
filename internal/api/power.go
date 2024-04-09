package turingpi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPower - Returns Turing Pi BMC SD card status.
func (c *Client) GetPower(ctx context.Context) (Power, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.APIURI+"?opt=get&type=power",
		nil,
	)
	if err != nil {
		return Power{}, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return Power{}, fmt.Errorf("failed to call API: %w", err)
	}

	sdCardResponse := PowerResponse{}

	err = json.Unmarshal(body, &sdCardResponse)
	if err != nil {
		return Power{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return sdCardResponse.Response[0], nil
}
