package eosio

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
)

const (
	ChainName = "EOS"
	Coin      = "EOS"
)

var (
	Startblock = 0
	Endblock   = 999999999
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client *EosClient
}

func (w *WalletAdaptor) GetLatestSafeBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetLatestFinalizedBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBlockHeaderByHash(req *wallet2.BlockHeaderByHashRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBlockByRange(req *wallet2.BlockByRangeRequest) (*wallet2.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetTxReceiptByHash(req *wallet2.TxReceiptByHashRequest) (*wallet2.TxReceiptByHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetStorageHash(req *wallet2.StorageHashRequest) (*wallet2.StorageHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetFilterLogs(req *wallet2.FilterLogsRequest) (*wallet2.FilterLogsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetTxCountByAddress(req *wallet2.TxCountByAddressRequest) (*wallet2.TxCountByAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetSuggestGasPrice(req *wallet2.SuggestGasPriceRequest) (*wallet2.SuggestGasPriceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetSuggestGasTipCap(req *wallet2.SuggestGasPriceRequest) (*wallet2.SuggestGasPriceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBlockByNumber(req *wallet2.BlockInfoRequest) (*wallet2.BlockInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBlockHeaderByNumber(req *wallet2.BlockHeaderRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	client, err := newClient(conf.Fullnode.Eos.RPCs[0].RPCURL)
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
		return nil, err
	}
	return &WalletAdaptor{
		client: client,
	}, nil
}

func (w *WalletAdaptor) getClient() *EosClient {
	return w.client
}

func (w *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	res, err := w.getClient().GetAccount(req.Address)
	if err != nil {
		log.Error("GetBalance error", "err", err)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "GetBalance error",
			Balance: "0",
		}, err
	}
	return &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "GetBalance success",
		Balance: res,
	}, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	res, err := w.getClient().PushTransaction(req.ConsumerToken, req.RawTx)
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
		TxHash: res.TransactionID,
	}, nil
}

func (w *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	txMsgList := []*wallet2.TxMessage{}
	pos := (req.Page - 1) * req.Pagesize
	offset := req.Pagesize
	res, err := w.getClient().GetActions(req.ContractAddress, int64(pos), int64(offset))
	if err != nil {
		log.Error("GetTxByAddress error", "err", err)
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetTxByAddress error",
			Tx:   txMsgList,
		}, err
	}
	for _, action := range res.Actions {
		txMsgList = append(txMsgList, &wallet2.TxMessage{
			Hash:  string(action.Trace.TransactionID),
			Index: uint32(action.Trace.ActionOrdinal),
		})
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "GetTxByAddress success",
		Tx:   txMsgList,
	}, nil
}

func (w *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	res, err := w.getClient().GetTransaction(req.Hash)
	if err != nil {
		log.Error("GetTxByHash error", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetTxByHash error",
		}, err
	}
	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "GetTxByHash success",
		Tx: &wallet2.TxMessage{
			Hash:   string(res.ID),
			Status: wallet2.TxStatus(res.Receipt.Status),
		},
	}, nil
}

func (w *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal([]byte(req.JsonString), &dataMap)
	if err != nil {
		log.Error("ABIJSONToBin Unmarshal error", "err", err)
		return &wallet2.ABIJSONToBinResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "ABIJSONToBin Unmarshal error",
		}, err
	}
	res, err := w.getClient().ABIJSONToBin(
		req.Code,
		req.Action,
		dataMap,
	)
	if err != nil {
		log.Error("ABIJSONToBin error", "err", err)
		return &wallet2.ABIJSONToBinResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "ABIJSONToBin error",
		}, err
	}
	return &wallet2.ABIJSONToBinResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "ABIJSONToBin success",
		Bin:  res,
	}, nil
}

func (w *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	res, err := w.getClient().ABIBinToJSON(
		req.Code,
		req.Action,
		req.Bin,
	)
	if err != nil {
		log.Error("ABIBinToJSON error", "err", err)
		return &wallet2.ABIBinToJSONResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "ABIBinToJSON error",
		}, err
	}
	resByte, err := json.Marshal(res)
	if err != nil {
		log.Error("ABIBinToJSON Marshal error", "err", err)
		return &wallet2.ABIBinToJSONResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "ABIBinToJSON Marshal error",
		}, err
	}
	return &wallet2.ABIBinToJSONResponse{
		Code:       common.ReturnCode_SUCCESS,
		Msg:        "ABIBinToJSON success",
		JsonString: string(resByte),
	}, nil
}

func (w *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	//TODO implement me
	panic("implement me")
}
