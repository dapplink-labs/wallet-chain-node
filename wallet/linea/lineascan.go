package linea

import (
	"time"

	etherscan "github.com/nanmu42/etherscan-api"
)

func NewEtherscanClient(apiURL string, Key string) *etherscan.Client {
	return etherscan.NewCustomized(etherscan.Customization{
		Timeout: 30 * time.Second,
		Key:     Key,
		BaseURL: apiURL,
	})
}
