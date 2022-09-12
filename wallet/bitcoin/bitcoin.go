package bitcoin

import (
	"bytes"
	"github.com/SavourDao/savour-hd/config"
	"github.com/SavourDao/savour-hd/rpc/common"
	wallet2 "github.com/SavourDao/savour-hd/rpc/wallet"
	"github.com/SavourDao/savour-hd/wallet"
	"github.com/SavourDao/savour-hd/wallet/fallback"
	"github.com/SavourDao/savour-hd/wallet/multiclient"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"strconv"
	"strings"
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

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
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

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	tx, err := a.getClient().GetRawTransactionVerbose((*chainhash.Hash)([]byte((req.Hash))))
	if err != nil {
		if rpcErr, ok := err.(*btcjson.RPCError); ok && rpcErr.Code == btcjson.ErrRPCBlockNotFound {
			return &wallet2.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "Tx status not found",
			}, nil
		}
		log.Error("queryTransaction GetRawTransactionVerbose", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if tx == nil || req.Hash != tx.Txid {
		log.Error("queryTransaction txid mismatch")
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	if tx.Confirmations < confirms {
		log.Error("queryTransaction confirmes too low", "tx confirms", tx.Confirmations, "need confirms", confirms)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	blockHash, _ := chainhash.NewHashFromStr(tx.BlockHash)
	block, err := a.getClient().GetBlockVerbose(blockHash)
	if err != nil {
		log.Error("queryTransaction GetBlockVerbose", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	reply, err := a.assembleUtxoTransactionReply(tx, block.Height, block.Time, func(txid string, index uint32) (int64, string, error) {
		preHash, err2 := chainhash.NewHashFromStr(txid)
		if err2 != nil {
			return 0, "", err2
		}
		preTx, err2 := a.getClient().GetRawTransactionVerbose(preHash)
		if err2 != nil {
			return 0, "", err2
		}
		amount := btcToSatoshi(preTx.Vout[index].Value).Int64()

		return amount, preTx.Vout[index].ScriptPubKey.Addresses[0], nil
	})
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return reply, nil
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
	return &wallet2.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support this interface",
	}, nil
}

func (w *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return &wallet2.SupportCoinsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
		Support: true,
	}, nil
}

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	return &wallet2.NonceResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support this interface",
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	reply, err := w.getClient().EstimateSmartFee(btcFeeBlocks)
	if err != nil {
		log.Info("QueryGasPrice", "err", err)
		return &wallet2.GasPriceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	price := btcToSatoshi(reply.Feerate)
	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get gas price success",
		Gas:  price.String(),
	}, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	r := bytes.NewReader([]byte(req.RawTx))
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	txHash, err := w.getClient().SendRawTransaction(&msgTx)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if strings.Compare(msgTx.TxHash().String(), txHash.String()) != 0 {
		log.Error("BroadcastTransaction, txhash mismatch", "local hash", msgTx.TxHash().String(), "hash from net", txHash.String(), "signedTx", req.RawTx)
	}

	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: txHash.String(),
	}, nil
}

func btcToSatoshi(btcCount float64) *big.Int {
	amount := strconv.FormatFloat(btcCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(btcDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}

func (a *WalletAdaptor) assembleUtxoTransactionReply(tx *btcjson.TxRawResult, blockHeight, blockTime int64, getPrevTxInfo func(txid string, index uint32) (int64, string, error)) (*wallet2.TxHashResponse, error) {
	var totalAmountIn, totalAmountOut int64
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	var direction int32
	for _, in := range tx.Vin {
		amount, address, err := getPrevTxInfo(in.Txid, in.Vout)
		if err != nil {
			return nil, err
		}
		totalAmountIn += amount
		from_addrs = append(from_addrs, &wallet2.Address{Address: address})
		value_list = append(value_list, &wallet2.Value{Value: strconv.FormatInt(totalAmountIn, 10)})
	}
	for _, out := range tx.Vout {
		amount := btcToSatoshi(out.Value).Int64()
		totalAmountOut += amount
		addr := ""
		if len(out.ScriptPubKey.Addresses) > 0 {
			addr = out.ScriptPubKey.Addresses[0]
		}
		to_addrs = append(to_addrs, &wallet2.Address{Address: addr})
		value_list = append(value_list, &wallet2.Value{Value: strconv.FormatInt(totalAmountOut, 10)})
	}
	gasUsed := totalAmountIn - totalAmountOut
	reply := &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get tx success",
		Tx: &wallet2.TxMessage{
			Hash:   tx.Hash,
			Status: wallet2.TxStatus_Success,
			From:   from_addrs,
			To:     to_addrs,
			Fee:    strconv.FormatInt(gasUsed, 10),
			Value:  value_list,
			Height: strconv.FormatInt(blockHeight, 10),
			Type:   direction,
		},
	}
	return reply, nil
}
