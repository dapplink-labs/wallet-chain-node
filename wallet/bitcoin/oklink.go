package bitcoin

import (
	"fmt"
	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/savour-labs/wallet-hd-chain/wallet/bitcoin/types"
)

/*
 * api docs:
 *	https://www.oklink.com/docs/zh/#explorer-api-chain-info-query-the-best-handling-fee-or-gas-fee
 */
var errOkLinkHTTPError = errors.New("ok link http error")

type OkLinkClient struct {
	client *gresty.Client
}

func NewOkLinkClient(url string) (*OkLinkClient, error) {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errOkLinkHTTPError)
		}
		return nil
	})
	return &OkLinkClient{
		client: client,
	}, nil
}

func (c *OkLinkClient) GetGasFee(coinName string) (string, error) {
	var GasFeeData types.GasFeeData
	okAccessKey := make(map[string]string)
	okAccessKey["Ok-Access-Key"] = "ac6562a6-825e-4a16-8200-de3879da9b73"
	response, err := c.client.R().
		SetHeaders(okAccessKey).
		SetResult(&GasFeeData).
		Get("/api/v5/explorer/blockchain/fee?chainShortName=" + coinName)
	if err != nil {
		return "0", fmt.Errorf("cannot get account balance: %w", err)
	}
	if response.StatusCode() != 200 {
		return "0", errors.New("get account balance fail")
	}
	return GasFeeData.Data[0].BestTransactionFee, nil
}
