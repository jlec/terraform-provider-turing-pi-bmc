package turingpi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetNodeInfo - Returns Turing Pi BMC SD card status.
func (c *Client) GetNodeInfo(ctx context.Context) (NodeInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.APIURI+"?opt=get&type=nodeinfo",
		nil,
	)
	if err != nil {
		return NodeInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return NodeInfo{}, fmt.Errorf("failed to call API: %w", err)
	}

	sdCardResponse := NodeInfoResponse{}

	err = json.Unmarshal(body, &sdCardResponse)
	if err != nil {
		return NodeInfo{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return sdCardResponse.Response[0], nil
}
