package sui

import (
	"encoding/json"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const (
	ChainName   = "Sui"
	Coin        = "SUI"
	SuiCoinType = "0x2::sui::SUI"
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

func (a *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := NewSuiClient(conf)
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

func newWalletAdaptor(client *suiClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}
func (a *WalletAdaptor) getClient() *suiClient {
	return a.clients.BestClient().(*suiClient)
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}
func (a *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	price, err := a.getClient().GetGasPrice()
	if err != nil {
		log.Info("QueryGasPrice", "err", err)
		return &wallet2.GasPriceResponse{Code: common.ReturnCode_ERROR,
			Msg: err.Error()}, err
	}
	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get gas price success",
		Gas:  big.NewInt(int64(price)).String(),
	}, nil
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, err := a.getClient().GetAccountBalance(req.Address, req.Coin)
	if err != nil {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get sui balance fail",
			Balance: "0",
		}, err
	}
	return &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
		Balance: balance.TotalBalance,
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {

	cursor := req.Cursor
	txList, err := a.getClient().GetTxListByAddress(req.Address, cursor, req.Pagesize)
	if err != nil {
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transactions fail",
		}, err
	}
	// todo sui 专有的交易结构，直接放到value中，前端自定义获取解析
	var tx_list []*wallet2.TxMessage
	for _, tx := range txList.Data {
		message, _ := a.getTxMessage(tx)
		tx_list = append(tx_list, message)
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transactions success",
		Tx:   tx_list,
	}, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	txDetail, err := a.getClient().GetTxDetailByDigest(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction fail",
		}, err
	}

	message, _ := a.getTxMessage(txDetail)

	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   message,
	}, nil
}

func (a *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	_, err := a.getClient().SendTx(req.RawTx)
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

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	address := req.Address
	cursor := req.Cursor
	limit := req.Limit
	coinType := req.CoinType

	coins, err := a.getClient().GetCoins(address, coinType, cursor, limit)
	if err != nil {
		return &wallet2.UnspentOutputsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get unspentOutputs fail",
		}, err
	}
	var unspentOutputList []*wallet2.UnspentOutput
	for _, value := range coins.Data {
		//value.
		marshal, jsonErr := json.Marshal(value)
		if jsonErr != nil {
			return nil, jsonErr
		}
		unspentOutput := &wallet2.UnspentOutput{
			Script:          value.Digest,
			TxHashBigEndian: string(marshal),
			TxHash:          value.CoinObjectId,
			Value:           stringToInt(value.Balance).Uint64(),
		}
		unspentOutputList = append(unspentOutputList, unspentOutput)
	}
	return &wallet2.UnspentOutputsResponse{
		Code:           common.ReturnCode_SUCCESS,
		Msg:            "get unspentOutputs success",
		UnspentOutputs: unspentOutputList,
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) getTxMessage(suiTransaction models.SuiTransactionBlockResponse) (*wallet2.TxMessage, error) {
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	totalAmount := big.NewInt(0)
	toAmount := big.NewInt(0)
	for _, bc := range suiTransaction.BalanceChanges {
		if bc.Owner.AddressOwner != "" {
			from_addrs = append(from_addrs, &wallet2.Address{Address: bc.Owner.AddressOwner})
			totalAmount = new(big.Int).Add(totalAmount, stringToInt(bc.Amount))
		} else {
			to_addrs = append(to_addrs, &wallet2.Address{Address: bc.Owner.ObjectOwner})
			toAmount = new(big.Int).Add(toAmount, stringToInt(bc.Amount))
			value_list = append(value_list, &wallet2.Value{Value: bc.Amount})
		}
	}
	totalAmount = new(big.Int).Abs(totalAmount)
	fee := new(big.Int).Sub(totalAmount, toAmount).String()
	return &wallet2.TxMessage{
		Hash:     suiTransaction.Digest,
		Height:   suiTransaction.Checkpoint,
		Status:   wallet2.TxStatus_Success,
		Type:     0,
		Datetime: suiTransaction.TimestampMs,
		Froms:    from_addrs,
		Tos:      to_addrs,
		Values:   value_list,
		Fee:      fee,
	}, nil
}

func stringToInt(amount string) *big.Int {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}
