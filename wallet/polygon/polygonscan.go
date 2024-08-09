package polygon

import (
	"time"

	etherscan "github.com/the-web3/etherscan-api"
)

func NewEtherscanClient(apiURL string, Key string) *etherscan.Client {
	return etherscan.NewCustomized(etherscan.Customization{
		Timeout: 30 * time.Second,
		Key:     Key,
		BaseURL: apiURL,
	})
}
