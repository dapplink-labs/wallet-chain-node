package xrp

import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const (
	ChainName = "XRP"
	Symbol    = "XRP"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *multiclient.MultiClient
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

func (w *WalletAdaptor) getClient() *Client {
	return w.clients.BestClient().(*Client)
}

func (w *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, err := w.getClient().GetBalance(req.Address)
	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance error",
			Balance: "0",
		}, err
	}
	return &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "get balance success",
		Balance: balance,
	}, nil
}

func (w *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	txs, err := w.getClient().GetTxsByAddress(req.Address)

	if err != nil {
		log.Error("get GetTxByAddress error", "err", err)
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "send tx fail",
		}, err
	}

	list := make([]*wallet2.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		list = append(list, &wallet2.TxMessage{
			Hash:     txs[i].Hash,
			Tos:      []*wallet2.Address{{Address: txs[i].To}},
			Froms:    []*wallet2.Address{{Address: txs[i].From}},
			Fee:      txs[i].Fee,
			Status:   wallet2.TxStatus_Success,
			Values:   []*wallet2.Value{{Value: txs[i].Amount}},
			Type:     1,
			Height:   txs[i].BlockHeight,
			Datetime: "",
		})
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		Tx:   list,
	}, nil
}

func (w *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	tx, err := w.getClient().GetTxByHash(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
			Tx:   nil,
		}, err
	}
	return &wallet2.TxHashResponse{
		Tx: &wallet2.TxMessage{
			Hash:     tx.Hash,
			Tos:      []*wallet2.Address{{Address: tx.To}},
			Froms:    []*wallet2.Address{{Address: tx.From}},
			Fee:      tx.Fee,
			Status:   wallet2.TxStatus_Success,
			Values:   []*wallet2.Value{{Value: tx.Amount}},
			Type:     1,
			Height:   tx.BlockHeight,
			Datetime: "",
		},
	}, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	value, err := w.getClient().SendTx(req.RawTx)
	if err != nil {
		log.Error("get SendTx error", "err", err)
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "send tx fail",
			TxHash: "",
		}, err
	}
	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: value,
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

func (w *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newClients(conf)
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
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
