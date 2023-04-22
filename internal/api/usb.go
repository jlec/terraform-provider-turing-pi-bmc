package turingpi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUsb - Returns Turing Pi BMC SD card status.
func (c *Client) GetUsb() (Usb, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?opt=get&type=usb", c.ApiURI), nil)
	if err != nil {
		return Usb{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return Usb{}, err
	}

	sdCardResponse := UsbResponse{}
	err = json.Unmarshal(body, &sdCardResponse)

	if err != nil {
		return Usb{}, err
	}

	return sdCardResponse.Response[0], nil
}
