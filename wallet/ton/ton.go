package ton

import (
	"context"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
	"github.com/xssnick/tonutils-go/address"
	"math/big"
	"strconv"
)

const (
	ChainName = "Ton"
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

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newTonClients(conf)
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

func (w *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
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

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	ret, err := w.getClient().WalletInfo(req.Address)
	if err != nil {
		return nil, err
	}
	return &wallet2.NonceResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "get nonce success",
		Nonce: strconv.FormatUint(ret.Seqno, 10),
	}, nil
}

func (w *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	api := w.getClient().api
	block, err := api.CurrentMasterchainInfo(context.Background())
	acc, err := api.GetAccount(context.Background(), block, address.MustParseAddr(req.Address))

	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance error",
			Balance: "0",
		}, err
	}

	if !acc.IsActive {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "account is not active",
			Balance: "0",
		}, nil
	} else {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
			Balance: acc.State.Balance.String(),
		}, nil
	}
}

func (w *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {

	ret, err := w.getClient().GetTxByAddr(req.Address)
	if err != nil {
		return nil, err
	}

	var tx_list []*wallet2.TxMessage
	for _, transactionInfo := range ret.Transactions {
		txMessage, blockRespErr := GetTonTxMessage(ret, &transactionInfo)
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

func (w *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	ret, err := w.getClient().GetTxByTxHash(req.Hash)
	if err != nil {
		return nil, err
	}

	if len(ret.Transactions) == 0 {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transactions fail",
		}, nil
	}

	tx := ret.Transactions[0]

	txMsg, _ := GetTonTxMessage(ret, &tx)

	res := &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		Tx:   txMsg,
	}

	return res, nil
}

func (w *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {

	api := w.getClient().api
	block, err := api.CurrentMasterchainInfo(context.Background())

	if err != nil {
		log.Error("get account error", "err", err)
		return &wallet2.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get account error",
		}, err
	}

	acc, err := api.GetAccount(context.Background(), block, address.MustParseAddr(req.Address))
	if err != nil {
		log.Error("get account error", "err", err)
		return &wallet2.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get account error",
		}, err
	} else {
		utils.PrettyPrint(acc)
		lastLt := strconv.FormatUint(acc.State.LastTransactionLT, 10)
		return &wallet2.AccountResponse{
			Code:     common.ReturnCode_SUCCESS,
			Msg:      "get account success",
			Sequence: lastLt,
		}, nil
	}

}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	ret, err := w.getClient().PostSendTx(req.RawTx)
	if err != nil {
		return nil, err
	}
	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "success",
		TxHash: ret.Hash,
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	ret, err := w.getClient().PostEstimateFee(req.RawTx, req.Address)
	if err != nil {
		return nil, err
	}

	total, err := ret.SumFees()

	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		Gas:  strconv.Itoa(int(total)),
	}, nil
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

func (a *WalletAdaptor) getClient() *tonClient {
	return a.clients.BestClient().(*tonClient)
}

func GetTonTxMessage(ret *Tx, tx *Transactions) (*wallet2.TxMessage, error) {
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	totalAmount := big.NewInt(0)

	from_addrs = append(from_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, tx.InMsg.Source)})
	to_addrs = append(to_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, tx.InMsg.Destination)})

	if len(tx.InMsg.Value) > 0 {
		totalAmount = new(big.Int).Add(totalAmount, stringToInt(tx.InMsg.Value))
	}

	for _, out := range tx.OutMsgs {
		if len(out.Source) > 0 {
			from_addrs = append(from_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, out.Source)})
		}
		if len(out.Destination) > 0 {
			to_addrs = append(to_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, out.Destination)})
		}
		log.Info(totalAmount.String(), "value", out.Value)
		if len(out.Value) > 0 {
			totalAmount = new(big.Int).Sub(totalAmount, stringToInt(out.Value))
		}
	}

	txMsg := &wallet2.TxMessage{
		Hash:     tx.Hash,
		Froms:    from_addrs,
		Tos:      to_addrs,
		Fee:      tx.TotalFees,
		Values:   []*wallet2.Value{{Value: totalAmount.String()}},
		Status:   wallet2.TxStatus_Success,
		Datetime: strconv.Itoa(tx.Now),
		Height:   strconv.Itoa(tx.BlockRef.Seqno),
	}

	return txMsg, nil
}

func getUserFriendly(addressBook map[string]struct {
	UserFriendly string `json:"user_friendly"`
}, key string) string {
	if entry, ok := addressBook[key]; ok {
		return entry.UserFriendly
	}
	return ""
}

func stringToInt(amount string) *big.Int {
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}
