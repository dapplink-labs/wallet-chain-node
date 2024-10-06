package ada

import (
	"math/big"
	"time"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"

	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const (
	ChainName = "Ada"
	Coin      = "Ada"
	// agent is the user-agent on requests to the
	// Rosetta Server.
	agent = "rosetta-sdk-go"

	// defaultTimeout is the default timeout for
	// HTTP requests.
	defaultTimeout  = 10 * time.Second
	Input           = "input"
	Output          = "output"
	relativeTtl     = int64(1000)
	transactionSize = int64(5000)
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *multiclient.MultiClient
}

func (a *WalletAdaptor) GetLatestSafeBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetLatestFinalizedBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetBlockHeaderByHash(req *wallet2.BlockHeaderByHashRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetBlockByRange(req *wallet2.BlockByRangeRequest) (*wallet2.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxReceiptByHash(req *wallet2.TxReceiptByHashRequest) (*wallet2.TxReceiptByHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetStorageHash(req *wallet2.StorageHashRequest) (*wallet2.StorageHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetFilterLogs(req *wallet2.FilterLogsRequest) (*wallet2.FilterLogsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxCountByAddress(req *wallet2.TxCountByAddressRequest) (*wallet2.TxCountByAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetSuggestGasPrice(req *wallet2.SuggestGasPriceRequest) (*wallet2.SuggestGasPriceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetSuggestGasTipCap(req *wallet2.SuggestGasPriceRequest) (*wallet2.SuggestGasPriceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetBlockByNumber(req *wallet2.BlockInfoRequest) (*wallet2.BlockInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetBlockHeaderByNumber(req *wallet2.BlockHeaderRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := NewAdaClient(conf)
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

func newWalletAdaptor(client *adaClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *WalletAdaptor) getClient() *adaClient {
	return a.clients.BestClient().(*adaClient)
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, err := a.getClient().GetAccountBalance(req.Address)
	if err != nil {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get ada balance fail",
			Balance: "0",
		}, err
	}
	return &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "get btc balance success",
		Balance: balance[0].Value,
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	page := req.Page
	pagesize := req.Pagesize
	// api limit maxmum pagesize is 25
	limit := int64(pagesize)
	offset := int64((page - 1) * pagesize)
	transactionsResponse, err := a.getClient().GetTransactionsByAddress(req.Address, &limit, &offset)
	if err != nil {
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transactions fail",
		}, err
	}
	transactions := transactionsResponse.Transactions
	var tx_list []*wallet2.TxMessage
	for _, transactionInfo := range transactions {
		txMessage, blockRespErr := a.getTxMessage(transactionInfo)
		if blockRespErr != nil {
			return &wallet2.TxAddressResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "get transactions fail",
			}, blockRespErr
		}
		tx_list = append(tx_list, txMessage)
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transactions fail",
		Tx:   tx_list,
	}, err
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	transactionsByAddress, err := a.getClient().GetTransactionsByHash(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction fail",
		}, err
	}
	blockTransaction := transactionsByAddress[0]
	txMsg, blockRespErr := a.getTxMessage(blockTransaction)
	if blockRespErr != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction fail",
		}, blockRespErr
	}
	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   txMsg,
	}, nil
}

func (a *WalletAdaptor) getTxMessage(blockTransaction *types.BlockTransaction) (*wallet2.TxMessage, error) {
	transaction := blockTransaction.Transaction
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	input_value_sum := big.NewInt(0)
	output_value_sum := big.NewInt(0)
	for _, operation := range transaction.Operations {
		amount_value := stringToInt(operation.Amount.Value)
		address := operation.Account.Address
		if operation.Type == Input {
			from_addrs = append(from_addrs, &wallet2.Address{Address: address})
			input_value_sum.Add(input_value_sum, amount_value)
		}
		if operation.Type == Output {
			to_addrs = append(to_addrs, &wallet2.Address{Address: address})
			value_list = append(value_list, &wallet2.Value{Value: operation.Amount.Value})
			output_value_sum.Add(output_value_sum, amount_value)
		}
	}
	// compute tx fee
	tx_fee := new(big.Int).Sub(new(big.Int).Abs(input_value_sum), new(big.Int).Abs(output_value_sum))
	blockIdentifier := blockTransaction.BlockIdentifier
	// get block info ,use block time as tx time
	blockResponse, blockRespErr := a.getClient().GetBlockByHash(blockIdentifier.Hash, blockIdentifier.Index)
	if blockRespErr != nil {
		return nil, blockRespErr
	}
	datetime := big.NewInt(blockResponse.Block.Timestamp).String()
	txMsg := &wallet2.TxMessage{
		Hash:            transaction.TransactionIdentifier.Hash,
		Froms:           from_addrs,
		Tos:             to_addrs,
		Values:          value_list,
		Fee:             tx_fee.String(),
		Status:          wallet2.TxStatus_Success,
		Type:            0,
		Height:          big.NewInt(blockIdentifier.Index).String(),
		ContractAddress: "0x00",
		Datetime:        datetime,
	}
	return txMsg, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	txHash, err := w.getClient().SendRawTransaction(req.RawTx)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: txHash,
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	// todo 需要等待参数结构体修改
	//oprationStr := ""
	// get gas price preprocess
	//preprocess, preErr := w.getClient().Preprocess(oprationStr)
	//if preErr != nil {
	//	log.Info("preprocess", "err", preErr)
	//	return &wallet2.GasPriceResponse{
	//		Code: common.ReturnCode_ERROR,
	//		Msg:  preErr.Error(),
	//	}, preErr
	//}
	//relativeTtl := preprocess["relative_ttl"].(int64)
	//transactionSize := preprocess["transaction_size"].(int64)

	gasFee, err := w.getClient().GetTxFee(relativeTtl, transactionSize)
	if err != nil {
		log.Info("QueryGasPrice", "err", err)
		return &wallet2.GasPriceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get gas price success",
		Gas:  gasFee[0].Value,
	}, nil
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

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	coinsResponse, err := a.getClient().GetUtxoByAddress(req.Address)
	if err != nil {
		return &wallet2.UnspentOutputsResponse{
			Code:           common.ReturnCode_ERROR,
			Msg:            err.Error(),
			UnspentOutputs: nil,
		}, err
	}
	coins := coinsResponse.Coins
	var unspentOutputList []*wallet2.UnspentOutput
	for _, value := range coins {
		unspentOutput := &wallet2.UnspentOutput{
			Script: value.CoinIdentifier.Identifier,
			Value:  stringToInt(value.Amount.Value).Uint64(),
		}
		unspentOutputList = append(unspentOutputList, unspentOutput)
	}
	return &wallet2.UnspentOutputsResponse{
		Code:           common.ReturnCode_SUCCESS,
		Msg:            "get unspent outputs success",
		UnspentOutputs: unspentOutputList,
	}, err
}

func stringToInt(amount string) *big.Int {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}
