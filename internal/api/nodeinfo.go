package turingpi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetNodeInfo - Returns Turing Pi BMC SD card status.
func (c *Client) GetNodeInfo() (NodeInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?opt=get&type=nodeinfo", c.ApiURI), nil)
	if err != nil {
		return NodeInfo{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return NodeInfo{}, err
	}

	sdCardResponse := NodeInfoResponse{}
	err = json.Unmarshal(body, &sdCardResponse)

	if err != nil {
		return NodeInfo{}, err
	}

	return sdCardResponse.Response[0], nil
}
