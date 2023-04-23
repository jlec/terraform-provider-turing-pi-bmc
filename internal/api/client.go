package turingpi

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL
// const HostURL string = "http://localhost:19090"
const TimeOut int = 10

type Client struct {
	ApiURI     string
	HTTPClient *http.Client
}

func NewClient(endpoint string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: time.Duration(TimeOut) * time.Second},
		ApiURI:     fmt.Sprintf("http://%s/api/bmc", endpoint),
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) Get(params string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", c.ApiURI, params), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) Set(params string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", c.ApiURI, params), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}
