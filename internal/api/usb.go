package turingpi

import (
	"encoding/json"
	"fmt"
)

// GetUsb - Returns Turing Pi BMC SD card status.
func (c *Client) GetUsb() (Usb, error) {
	response := UsbResponse{}

	body, err := c.Get("opt=get&type=usb")
	if err != nil {
		return Usb{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return Usb{}, err
	}

	return response.Response[0], nil
}

// SetUsb - Returns Turing Pi BMC SD card status.
func (c *Client) SetUsb(mode, node int64) (Result, error) {
	response := ResultResponse{}

	body, err := c.Set(fmt.Sprintf("opt=set&type=usb&mode=%d&node=%d", mode, node))
	if err != nil {
		return Result{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return Result{}, err
	}

	if response.Response[0].Result != "ok" {
		return response.Response[0], &ResultError{Reason: response.Response[0].Result}
	}

	return response.Response[0], nil
}
