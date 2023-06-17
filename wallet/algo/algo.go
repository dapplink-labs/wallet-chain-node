package algo

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/wallet-hd-chain/cache"
	"github.com/savour-labs/wallet-hd-chain/config"
	"github.com/savour-labs/wallet-hd-chain/rpc/common"
	wallet2 "github.com/savour-labs/wallet-hd-chain/rpc/wallet"
	"github.com/savour-labs/wallet-hd-chain/wallet"
	"github.com/savour-labs/wallet-hd-chain/wallet/fallback"
	"github.com/savour-labs/wallet-hd-chain/wallet/multiclient"
	"math/big"
	"strconv"
	"strings"
)

const (
	ChainName = "Algo"
	Coin      = "Algo"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients     *multiclient.MultiClient
	algoscanCli *AlgoClient
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newAlgoClients(conf)
	if err != nil {
		return nil, err
	}
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}
	return &WalletAdaptor{
		clients:     multiclient.New(clis),
		algoscanCli: NewAlgoScanClient(conf),
	}, nil
}

func NewLocalWalletAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newWalletAdaptor(newLocalAlgoClient(network))
}

func newWalletAdaptor(client *algoClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *WalletAdaptor) getClient() *algoClient {
	return a.clients.BestClient().(*algoClient)
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (wa *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	key := strings.Join([]string{req.Chain, req.Coin, req.Address, req.ContractAddress}, ":")
	balanceCache := cache.GetBalanceCache()
	if r, exist := balanceCache.Get(key); exist {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
			Balance: r.(*big.Int).String(),
		}, nil
	}
	result, err := wa.getClient().GetAccountBalance(req.Address).Do(context.Background())
	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}, err
	} else {
		balanceCache.Add(key, result)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
			Balance: strconv.FormatUint(result.Amount, 10),
		}, nil
	}
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

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}
