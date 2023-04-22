package turingpi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPower - Returns Turing Pi BMC SD card status.
func (c *Client) GetPower() (Power, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?opt=get&type=power", c.ApiURI), nil)
	if err != nil {
		return Power{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return Power{}, err
	}

	sdCardResponse := PowerResponse{}
	err = json.Unmarshal(body, &sdCardResponse)

	if err != nil {
		return Power{}, err
	}

	return sdCardResponse.Response[0], nil
}
