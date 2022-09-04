package ethereum

import (
	"context"
	"fmt"
	"github.com/SavourDao/savour-hd/cache"
	"github.com/SavourDao/savour-hd/config"
	"github.com/SavourDao/savour-hd/rpc/common"
	wallet2 "github.com/SavourDao/savour-hd/rpc/wallet"
	"github.com/SavourDao/savour-hd/wallet"
	"github.com/SavourDao/savour-hd/wallet/fallback"
	"github.com/SavourDao/savour-hd/wallet/multiclient"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/nanmu42/etherscan-api"
	"math/big"
	"strconv"
	"strings"
)

const (
	ChainName = "eth"
	Coin      = "eth"
)

var (
	Startblock = 0
	Endblock   = 999999999
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients      *multiclient.MultiClient
	etherscanCli *etherscan.Client
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
	fmt.Println("1111")
	fmt.Println(conf.Fullnode.Eth.TpApiUrl)
	fmt.Println("1111")
	return &WalletAdaptor{
		clients:      multiclient.New(clis),
		etherscanCli: NewEtherscanClient(conf.Fullnode.Eth.TpApiUrl, conf.Fullnode.Eth.TpApiKey),
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
	key := strings.Join([]string{req.Chain, req.Coin, req.Address, req.ContractAddress}, ":")
	balanceCache := cache.GetBalanceCache()
	if r, exist := balanceCache.Get(key); exist {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
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
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}, err
	} else {
		balanceCache.Add(key, result)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
			Balance: result.String(),
		}, nil
	}
}

func (wa *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	key := strings.Join([]string{req.Coin, req.Address, strconv.Itoa(int(req.Page)), strconv.Itoa(int(req.Pagesize))}, ":")
	txCache := cache.GetTxCache()
	if r, exist := txCache.Get(key); exist {
		return r.(*wallet2.TxAddressResponse), nil
	}
	txs, err := wa.etherscanCli.NormalTxByAddress(req.Address, &Startblock, &Endblock, int(req.Page), int(req.Pagesize), false)
	if err != nil {
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	var tx_list []*wallet2.TxMessage
	for _, ktx := range txs {
		var from_addrs []*wallet2.Address
		var to_addrs []*wallet2.Address
		var value_list []*wallet2.Value
		from_addrs = append(from_addrs, &wallet2.Address{Address: ktx.From})
		to_addrs = append(to_addrs, &wallet2.Address{Address: ktx.To})
		value_list = append(value_list, &wallet2.Value{Value: ktx.Value.Int().String()})
		tx := &wallet2.TxMessage{
			Hash:            ktx.Hash,
			From:            from_addrs,
			To:              to_addrs,
			Value:           value_list,
			Fee:             "0",
			Status:          false,
			Type:            1,
			Height:          string(rune(ktx.BlockNumber)),
			ContractAddress: ktx.ContractAddress,
		}
		tx_list = append(tx_list, tx)
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get address tx success",
		Tx:   tx_list,
	}, nil
}

func (wa *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	key := strings.Join([]string{req.Coin, req.Hash}, ":")
	txCache := cache.GetTxCache()
	if r, exist := txCache.Get(key); exist {
		return r.(*wallet2.TxHashResponse), nil
	}

	tx, _, err := wa.getClient().TransactionByHash(context.TODO(), ethcommon.HexToHash(req.Hash))
	if err != nil {
		if err == ethereum.NotFound {
			return &wallet2.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "Ethereum Tx NotFoun",
			}, nil
		}
		log.Error("get transaction error", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Ethereum Tx NotFoun",
		}, nil
	}

	receipt, err := wa.getClient().TransactionReceipt(context.TODO(), ethcommon.HexToHash(req.Hash))
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Get transaction receipt error",
		}, nil
	}
	ok := true
	if receipt.Status != 0 {
		ok = false
	}
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	from_addrs = append(from_addrs, &wallet2.Address{Address: ""})
	to_addrs = append(to_addrs, &wallet2.Address{Address: tx.To().Hex()})
	value_list = append(value_list, &wallet2.Value{Value: tx.Value().String()})
	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &wallet2.TxMessage{
			Hash:            tx.Hash().Hex(),
			Index:           1,
			From:            from_addrs,
			To:              to_addrs,
			Value:           value_list,
			Fee:             tx.GasFeeCap().String(),
			Status:          ok,
			Type:            1,
			Height:          receipt.BlockNumber.String(),
			ContractAddress: tx.To().String(),
		},
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

func (wa *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return nil, nil
}

func (wa *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	var bockHeight *big.Int
	nonce, err := wa.getClient().NonceAt(context.TODO(), ethcommon.HexToAddress(req.Address), bockHeight)
	if err != nil {
		log.Error("get nonce failed", "err", err)
		return &wallet2.NonceResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   "get nonce failed",
			Nonce: strconv.FormatUint(nonce, 10),
		}, nil
	}
	return &wallet2.NonceResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "get nonce success",
		Nonce: strconv.FormatUint(nonce, 10),
	}, nil
}

func (wa *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	price, err := wa.getClient().SuggestGasPrice(context.TODO())
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &wallet2.GasPriceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get gas price fail",
			Gas:  "",
		}, nil
	}
	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get gas price success",
		Gas:  price.String(),
	}, nil
}

func (wa *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	signedTx := new(types.Transaction)
	if err := rlp.DecodeBytes([]byte(req.RawTx), signedTx); err != nil {
		log.Error("signedTx DecodeBytes failed", "err", err)
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "Send tx fail",
			TxHash: "",
		}, err
	}
	log.Info("broadcast tx", "tx", hexutil.Encode([]byte(req.RawTx)))
	txHash := fmt.Sprintf("0x%x", signedTx.Hash())
	if err := wa.getClient().SendTransaction(context.TODO(), signedTx); err != nil {
		log.Error("braoadcast tx failed", "tx_hash", txHash, "err", err)
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "Send tx fail",
			TxHash: "",
		}, err
	}
	log.Info("braoadcast tx success", "tx_hash", txHash)
	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "Send and  braoadcast tx success",
		TxHash: txHash,
	}, nil
}
