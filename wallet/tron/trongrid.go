package tron

import (
	"fmt"
	"strings"

	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var errTronGridHTTPError = errors.New("tron grid http error")

type TronGridClient struct {
	client *gresty.Client
}

func NewTronGridClient(url string) (*TronGridClient, error) {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errTronGridHTTPError)
		}
		return nil
	})
	return &TronGridClient{
		client: client,
	}, nil
}

func (tgc *TronGridClient) GetBalance(address string) error {
	var balanceData interface{}
	okAccessKey := make(map[string]string)
	okAccessKey["Ok-Access-Key"] = "ac6562a6-825e-4a16-8200-de3879da9b73"
	payload := strings.NewReader("{\"account_identifier\":{\"address\":\"TLLM21wteSPs4hKjbxgmH1L6poyMjeTbHm\"},\"block_identifier\":{\"hash\":\"0000000000010c4a732d1e215e87466271e425c86945783c3d3f122bfa5affd9\",\"number\":68682},\"visible\":true}")
	response, err := tgc.client.R().
		SetHeaders(okAccessKey).
		SetBody(payload).
		SetResult(&balanceData).
		Post("wallet/getaccountbalance")
	if err != nil {
		return nil
	}
	if response.StatusCode() != 200 {
		return nil
	}
	return nil
}
