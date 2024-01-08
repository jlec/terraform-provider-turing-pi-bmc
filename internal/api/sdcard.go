package turingpi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSDCard - Returns Turing Pi BMC SD card status.
func (c *Client) GetSDCard(ctx context.Context) (SDCard, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?opt=get&type=sdcard", c.APIURI),
		nil,
	)
	if err != nil {
		return SDCard{}, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return SDCard{}, fmt.Errorf("failed to call API: %w", err)
	}

	sdCardResponse := SDCardResponse{}
	err = json.Unmarshal(body, &sdCardResponse)

	if err != nil {
		return SDCard{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return sdCardResponse.Response[0], nil
}
