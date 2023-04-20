package turingpi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSDCard - Returns Turing Pi BMC SD card status.
func (c *Client) GetSDCard() (SDCard, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?opt=get&type=sdcard", c.ApiURI), nil)
	if err != nil {
		return SDCard{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return SDCard{}, err
	}

	sdCardResponse := SDCardResponse{}
	err = json.Unmarshal(body, &sdCardResponse)

	if err != nil {
		return SDCard{}, err
	}

	return sdCardResponse.Response[0], nil
}
