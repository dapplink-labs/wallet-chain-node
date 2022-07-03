package wallet

import (
	wallet2 "github.com/SavourDao/savour-core/rpc/savourrpc/go-savourrpc/wallet"
)

type WalletAdaptor interface {
	GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error)
	GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error)
	GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error)
	SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error)
}
