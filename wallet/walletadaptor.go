package wallet

import (
	wallet2 "github.com/SavourDao/savour-core/rpc/wallet"
)

type WalletAdaptor interface {
	GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error)
	GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error)
	GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error)
	SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error)
	GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error)
	GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error)
	GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error)
	GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error)
	GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error)
	GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error)
}
