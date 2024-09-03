package myip

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-resty/resty/v2"
)

func NewIPIFYClient() Client {
	return &ipifyClient{
		client: resty.New().
			SetRetryCount(5).
			SetJSONUnmarshaler(json.Unmarshal),
	}
}

const ipifyEndpoint = "https://api.ipify.org"

type ipifyClient struct {
	client *resty.Client
}

type ipifyResult struct {
	IP string `json:"ip"`
}

func (c *ipifyClient) GetIPAddress(ctx context.Context) (string, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&ipifyResult{}).
		SetQueryParam("format", "json").
		Get(ipifyEndpoint)
	if err != nil {
		return "", err
	}
	result, ok := resp.Result().(*ipifyResult)
	if !ok {
		return "", fmt.Errorf("invalid result type: %s", reflect.TypeOf(resp.Result()))
	}
	return result.IP, nil
}
