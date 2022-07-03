package bitcoin

import (
	"github.com/SavourDao/savour-core/config"
	"github.com/SavourDao/savour-core/rpc/common"
	wallet2 "github.com/SavourDao/savour-core/rpc/wallet"
	"github.com/SavourDao/savour-core/wallet"
	"github.com/SavourDao/savour-core/wallet/fallback"
	"github.com/SavourDao/savour-core/wallet/multiclient"
)

const (
	confirms     = 1
	btcDecimals  = 8
	btcFeeBlocks = 3
	ChainName    = "btc"
	Symbol       = "btc"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *multiclient.MultiClient
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newBtcClients(conf)
	if err != nil {
		return nil, err
	}
	return newChainAdaptorWithClients(clients), nil
}

func NewLocalChainAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newChainAdaptorWithClients([]*btcClient{newLocalBtcClient(network)})
}

func newChainAdaptorWithClients(clients []*btcClient) *WalletAdaptor {
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}
	return &WalletAdaptor{
		clients: multiclient.New(clis),
	}
}

func (a *WalletAdaptor) getClient() *btcClient {
	return a.clients.BestClient().(*btcClient)
}

func (w *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return nil, nil
}

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	return &wallet2.NonceResponse{
		Error: &common.Error{Code: 404},
		Nonce: "",
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	return nil, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	return nil, nil
}
