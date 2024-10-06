package solana

import (
	"fmt"
	"github.com/dapplink-labs/chain-explorer-api/common/account"

	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
)

const (
	ChainName = "Solana"
	Symbol    = "SOL"
	Coin      = "SOL"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client *SolanaClient
	sol    *SolScan
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

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
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

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, err := a.client.GetBalance(req.Address)
	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance error",
			Balance: "0",
		}, err
	} else {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
			Balance: balance,
		}, nil
	}
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	var resp *account.TransactionResponse[account.AccountTxResponse]
	var err error
	fmt.Println("req.ContractAddress", req.ContractAddress)
	if req.ContractAddress != "0x00" && req.ContractAddress != "" {
		log.Info("Spl token transfer record")
		resp, err = a.sol.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "spl")
	} else {
		log.Info("Sol transfer record")
		resp, err = a.sol.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "sol")
	}
	if err != nil {
		log.Error("get GetTxByAddress error", "err", err)
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx list fail",
			Tx:   nil,
		}, err
	} else {
		txs := resp.TransactionList
		list := make([]*wallet2.TxMessage, 0, len(txs))
		for i := 0; i < len(txs); i++ {
			list = append(list, &wallet2.TxMessage{
				Hash:   txs[i].TxId,
				Tos:    []*wallet2.Address{{Address: txs[i].To}},
				Froms:  []*wallet2.Address{{Address: txs[i].From}},
				Fee:    txs[i].TxId,
				Status: wallet2.TxStatus_Success,
				Values: []*wallet2.Value{{Value: txs[i].Amount}},
				Type:   1,
				Height: txs[i].Height,
			})
		}
		fmt.Println("resp", resp)
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  "get tx list success",
			Tx:   list,
		}, nil
	}

}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	tx, err := a.client.GetTxByHash(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
			Tx:   nil,
		}, err
	}
	return &wallet2.TxHashResponse{
		Tx: &wallet2.TxMessage{
			Hash:   tx.Hash,
			Tos:    []*wallet2.Address{{Address: tx.To}},
			Froms:  []*wallet2.Address{{Address: tx.From}},
			Fee:    tx.Fee,
			Status: wallet2.TxStatus_Success,
			Values: []*wallet2.Value{{Value: tx.Value}},
			Type:   tx.Type,
			Height: tx.Height,
		},
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	return &wallet2.AccountResponse{
		Code:          common.ReturnCode_ERROR,
		Msg:           "do not support",
		AccountNumber: "",
	}, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	value, err := a.client.GetMinRent()
	if err != nil {
		log.Error("get GetMinRent error", "err", err)
		return &wallet2.MinRentResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   err.Error(),
			Value: "",
		}, err
	} else {
		return &wallet2.MinRentResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "get mint rent success",
			Value: value,
		}, nil
	}
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	cli, err := NewSolanaClients(conf)

	sol, err := NewSolScanClient(conf.Fullnode.Sol.SolScanBaseUrl, conf.Fullnode.Sol.SolScanApiKey, conf.Fullnode.Sol.SolScanBaseTimeout)
	if err != nil {
		return nil, err
	}
	return &WalletAdaptor{
		client: cli,
		sol:    sol,
	}, nil
}

func (a *WalletAdaptor) getClient() *SolanaClient {
	return a.client
}

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	value, err := w.getClient().GetNonce(req.Address)
	if err != nil {
		log.Error("get GetNonce error", "err", err)
		return &wallet2.NonceResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   err.Error(),
			Nonce: "",
		}, err
	} else {
		return &wallet2.NonceResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "get nonce success",
			Nonce: value,
		}, nil
	}

}

func newWalletAdaptor(client *SolanaClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		client: client,
	}
}

func (w *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return &wallet2.SupportCoinsResponse{
		Code:    common.ReturnCode_ERROR,
		Msg:     "do not support",
		Support: false,
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "do not support",
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	return &wallet2.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "do not support",
	}, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	value, err := w.getClient().SendTx("")
	if err != nil {
		log.Error("get GetNonce error", "err", err)
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "send tx fail",
			TxHash: "",
		}, err
	} else {
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_SUCCESS,
			Msg:    "send tx success",
			TxHash: value,
		}, nil
	}

}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}
