package tezos

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/wallet-hd-chain/cache"
	"github.com/savour-labs/wallet-hd-chain/config"
	"github.com/savour-labs/wallet-hd-chain/rpc/common"
	wallet2 "github.com/savour-labs/wallet-hd-chain/rpc/wallet"
	"github.com/savour-labs/wallet-hd-chain/wallet"
	"github.com/savour-labs/wallet-hd-chain/wallet/fallback"
	"math/big"
	"strings"
)

const (
	ChainName = "Tezos"
	Coin      = "XTZ"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client       TezosClient
	oasisscanCli *ScanClient
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	client, err := NewTezosClient(conf.Fullnode.Tezos.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	return &WalletAdaptor{
		client:       client,
		oasisscanCli: nil,
	}, nil
}

func (w *WalletAdaptor) getClient() TezosClient {
	return w.client
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
	result, err := wa.getClient().GetAccountBalance(req.Address)
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
			Balance: result.Balance,
		}, nil
	}
}

func (wa *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	return nil, nil
}

func (wa *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	return nil, nil
}

func (wa *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return &wallet2.SupportCoinsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "this coin support",
		Support: true,
	}, nil
}

func (wa *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	nonce, err := wa.getClient().GetAccountCounter(req.Address)
	if err != nil {
		log.Error("get nonce failed", "err", err)
		return &wallet2.NonceResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   "get nonce failed",
			Nonce: nonce.Counter,
		}, nil
	}
	return &wallet2.NonceResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "get nonce success",
		Nonce: nonce.Counter,
	}, nil
}

func (wa *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	return nil, nil
}

func (wa *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	res, err := wa.getClient().SendRawTransaction(req.RawTx)
	if err != nil {
		log.Error("SendTx error", "err", err)
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SendTx error",
		}, err
	}
	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "SendTx success",
		TxHash: res.TxHash,
	}, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return nil, nil
}

// queryTransaction retrieve transaction information from a signed data.
func (a *WalletAdaptor) queryTransaction(isERC20 bool, tx *types.Transaction, receipt *types.Receipt, blockNumber uint64, signer types.Signer) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

// queryRawTransaction retrieve transaction information from a raw(unsigned) data.
func (a *WalletAdaptor) queryRawTransaction(isERC20 bool, rawTx *types.Transaction) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	return &wallet2.UtxoInsResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return &wallet2.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	return &wallet2.CreateUtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return &wallet2.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (wa *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	return &wallet2.AccountResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (wa *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	return &wallet2.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (wa *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	return &wallet2.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

type semaphore chan struct{}

func (s semaphore) Acquire() {
	s <- struct{}{}
}

func (s semaphore) Release() {
	<-s
}
