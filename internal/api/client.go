package turingpi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL
// const HostURL string = "http://localhost:19090"
const TimeOut int = 10

type Client struct {
	APIURI     string
	HTTPClient *http.Client
}

func NewClient(endpoint string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: time.Duration(TimeOut) * time.Second},
		APIURI:     fmt.Sprintf("http://%s/api/bmc", endpoint),
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failure: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failure: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: '%d', body: %s", ErrInvalidStatus, res.StatusCode, body)
	}

	return body, fmt.Errorf("failure: %w", err)
}

func (c *Client) Get(ctx context.Context, params string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?%s", c.APIURI, params),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failure: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failure: %w", err)
	}

	return body, nil
}

func (c *Client) Set(ctx context.Context, params string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s?%s", c.APIURI, params),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failure: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failure: %w", err)
	}

	return body, nil
}
