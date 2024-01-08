package turingpi

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetUsb - Returns Turing Pi BMC SD card status.
func (c *Client) GetUsb(ctx context.Context) (Usb, error) {
	response := UsbResponse{}

	body, err := c.Get(ctx, "opt=get&type=usb")
	if err != nil {
		return Usb{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return Usb{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response[0], nil
}

// SetUsb - Returns Turing Pi BMC SD card status.
func (c *Client) SetUsb(ctx context.Context, mode, node int64) (Result, error) {
	response := ResultResponse{}

	body, err := c.Set(ctx, fmt.Sprintf("opt=set&type=usb&mode=%d&node=%d", mode, node))
	if err != nil {
		return Result{}, fmt.Errorf("failed to set mode: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return Result{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Response[0].Result != "ok" {
		return response.Response[0], &ResultError{Reason: response.Response[0].Result}
	}

	return response.Response[0], nil
}

func APIToMode(mode int64) (string, error) {
	switch mode {
	case 0:
		return "host", nil
	case 1:
		return "device", nil
	default:
		return "", fmt.Errorf("%w: '%d'", ErrInvalidMode, mode)
	}
}

func ModeToAPI(mode string) (int64, error) {
	switch mode {
	case "host":
		return 0, nil
	case "device":
		return 1, nil
	default:
		return -1, fmt.Errorf("%w: '%s'", ErrInvalidMode, mode)
	}
}
