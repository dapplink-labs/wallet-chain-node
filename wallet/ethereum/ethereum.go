package ethereum

import (
	"context"
	"fmt"
	"github.com/SavourDao/savour-core/cache"
	"github.com/SavourDao/savour-core/config"
	"github.com/SavourDao/savour-core/rpc/common"
	wallet2 "github.com/SavourDao/savour-core/rpc/wallet"
	"github.com/SavourDao/savour-core/wallet"
	"github.com/SavourDao/savour-core/wallet/fallback"
	"github.com/SavourDao/savour-core/wallet/multiclient"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strconv"
	"strings"
)

const (
	ChainName = "eth"
	Symbol    = "eth"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *multiclient.MultiClient
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newEthClients(conf)
	if err != nil {
		return nil, err
	}
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}
	return &WalletAdaptor{
		clients: multiclient.New(clis),
	}, nil
}

func NewLocalWalletAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newWalletAdaptor(newLocalEthClient(network))
}

func newWalletAdaptor(client *ethClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *WalletAdaptor) getClient() *ethClient {
	return a.clients.BestClient().(*ethClient)
}

func stringToInt(amount string) *big.Int {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}

func (wa *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	key := strings.Join([]string{req.ChainId, req.CoinId, req.Address, req.ContractAddress}, ":")
	balanceCache := cache.GetBalanceCache()
	if r, exist := balanceCache.Get(key); exist {
		return &wallet2.BalanceResponse{
			Error:   &common.Error{Code: 2000},
			Balance: r.(*big.Int).String(),
		}, nil
	}
	var result *big.Int
	var err error
	if len(req.ContractAddress) > 0 {
		result, err = wa.getClient().Erc20BalanceOf(req.ContractAddress, req.Address, nil)
	} else {
		result, err = wa.getClient().BalanceAt(context.TODO(), ethcommon.HexToAddress(req.Address), nil)
	}
	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Error:   &common.Error{Code: 404},
			Balance: "0",
		}, err
	} else {
		balanceCache.Add(key, result)
		return &wallet2.BalanceResponse{
			Error:   &common.Error{Code: 2000},
			Balance: result.String(),
		}, nil
	}
}

func (wa *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (wa *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (wa *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (wa *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (wa *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (wa *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return nil, nil
}

func (wa *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	var bockHeight *big.Int
	nonce, err := wa.getClient().NonceAt(context.TODO(), ethcommon.HexToAddress(req.Address), bockHeight)
	if err != nil {
		log.Error("get nonce failed", "err", err)
		return &wallet2.NonceResponse{
			Error: &common.Error{Code: 404},
			Nonce: strconv.FormatUint(nonce, 10),
		}, nil
	}
	return &wallet2.NonceResponse{
		Error: &common.Error{Code: 2000},
		Nonce: strconv.FormatUint(nonce, 10),
	}, nil
}

func (wa *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	price, err := wa.getClient().SuggestGasPrice(context.TODO())
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &wallet2.GasPriceResponse{
			Error: &common.Error{Code: 404},
			Gas:   "",
		}, nil
	}
	return &wallet2.GasPriceResponse{
		Error: &common.Error{Code: 2000},
		Gas:   price.String(),
	}, nil
}

func (wa *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	signedTx := new(types.Transaction)
	if err := rlp.DecodeBytes([]byte(req.RawTx), signedTx); err != nil {
		log.Error("signedTx DecodeBytes failed", "err", err)
		return &wallet2.SendTxResponse{
			Error:  &common.Error{Code: 404},
			TxHash: "",
		}, err
	}
	log.Info("broadcast tx", "tx", hexutil.Encode([]byte(req.RawTx)))

	txHash := fmt.Sprintf("0x%x", signedTx.Hash())
	if err := wa.getClient().SendTransaction(context.TODO(), signedTx); err != nil {
		log.Error("braoadcast tx failed", "tx_hash", txHash, "err", err)
		return &wallet2.SendTxResponse{
			Error:  &common.Error{Code: 404},
			TxHash: "",
		}, err
	}
	log.Info("braoadcast tx success", "tx_hash", txHash)
	return &wallet2.SendTxResponse{
		Error:  &common.Error{Code: 404},
		TxHash: txHash,
	}, nil
}
