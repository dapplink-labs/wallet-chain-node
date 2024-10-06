package flows

import (
	"encoding/json"
	"fmt"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	"github.com/savour-labs/wallet-chain-node/wallet/flows/subgraph/transactions"
	"math/big"

	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const (
	ChainName        = "Flow"
	Coin             = "FLOW "
	FlowSubGraphUrl  = "https://api.findlabs.io/flowdiver/v1/graphql"
	FlowComputeLimit = 100
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

func (w *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, sequenceNumber, err := w.getClient().GetBalance(req.Address, req.ProposerKeyIndex)
	if err != nil {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get sui balance fail",
			Balance: "0",
		}, err
	}
	return &wallet2.BalanceResponse{
		Code:           common.ReturnCode_SUCCESS,
		Msg:            "success",
		Balance:        big.NewInt(balance).String(),
		SequenceNumber: sequenceNumber,
	}, nil
}

func (w *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {

	txListByAddress, err := w.getClient().GetTxListByAddress(req.Address, uint64(req.Pagesize), uint64(req.Page))

	if err != nil {
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transactions fail",
		}, err
	}
	participations := txListByAddress.Data.Participations
	var tx_list []*wallet2.TxMessage

	for _, participation := range participations {
		txMessage, msgErr := w.getTxMessage(participation.Transaction)
		if msgErr != nil {
			return &wallet2.TxAddressResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "get transactions fail",
			}, msgErr
		}
		tx_list = append(tx_list, txMessage)
	}

	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transactions success",
		Tx:   tx_list,
	}, nil
}

func (w *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	txDetailByHash, err := w.getClient().GetTxDetailByHash(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction detail fail",
		}, err
	}
	if len(txDetailByHash.Data.Transactions) == 0 {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction detail fail",
		}, err
	}
	transation := txDetailByHash.Data.Transactions[0]
	fmt.Println(transation)
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	from_addrs = append(from_addrs, &wallet2.Address{Address: transation.Payer})
	to_addrs = append(to_addrs, &wallet2.Address{Address: transation.Proposer})
	txStr, jsonErr := json.Marshal(transation)
	if jsonErr != nil {
		return nil, jsonErr
	}
	value_list = append(value_list, &wallet2.Value{
		Value: string(txStr),
	})
	tx_fee := big.NewFloat(transation.Fee).String()
	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transactions fail",
		Tx: &wallet2.TxMessage{
			Hash:            transation.ID,
			Froms:           from_addrs,
			Tos:             to_addrs,
			Values:          value_list,
			Fee:             tx_fee,
			Status:          wallet2.TxStatus_Success,
			Type:            0,
			Height:          big.NewInt(int64(transation.BlockHeight)).String(),
			ContractAddress: "0x00",
			Index:           uint32(transation.ProposerIndex),
			Datetime:        transation.Timestamp,
		},
	}, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	tx := req.RawTx
	err := w.getClient().SendTx(tx)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: "",
	}, nil
}

func (w *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {

	panic("implement me")
}

func (w *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
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

func (w *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	height, err := w.getClient().GetLatestBlockHeight()
	if err != nil {
		return &wallet2.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get sui balance fail",
		}, err
	}
	return &wallet2.BlockResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "success",
		Height: height,
	}, nil
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := NewFlowClient(conf)
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

func newWalletAdaptor(client *flowClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (w *WalletAdaptor) getClient() *flowClient {
	return w.clients.BestClient().(*flowClient)
}

func (w *WalletAdaptor) getTxMessage(transaction transactions.Transaction) (*wallet2.TxMessage, error) {
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	from_addrs = append(from_addrs, &wallet2.Address{Address: transaction.Payer})
	to_addrs = append(to_addrs, &wallet2.Address{Address: transaction.Proposer})
	txStr, jsonErr := json.Marshal(transaction)
	if jsonErr != nil {
		return nil, jsonErr
	}
	value_list = append(value_list, &wallet2.Value{
		Value: string(txStr),
	})
	tx_fee := big.NewFloat(transaction.Fee).String()
	txMsg := &wallet2.TxMessage{
		Hash:            transaction.ID,
		Froms:           from_addrs,
		Tos:             to_addrs,
		Values:          value_list,
		Fee:             tx_fee,
		Status:          wallet2.TxStatus_Success,
		Type:            0,
		Height:          big.NewInt(int64(transaction.BlockHeight)).String(),
		ContractAddress: "0x00",
		Index:           uint32(transaction.ProposerIndex),
		Datetime:        transaction.Timestamp.String(),
	}
	return txMsg, nil
}

func stringToInt(amount string) *big.Int {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}
