package solana

import (
	"github.com/SavourDao/savour-core/config"
	"github.com/SavourDao/savour-core/rpc/common"
	wallet2 "github.com/SavourDao/savour-core/rpc/wallet"
	"github.com/SavourDao/savour-core/wallet"
	"github.com/SavourDao/savour-core/wallet/fallback"
	"github.com/SavourDao/savour-core/wallet/multiclient"
	"github.com/ethereum/go-ethereum/log"
)

const (
	ChainName = "sol"
	Symbol    = "sol"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *multiclient.MultiClient
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, err := a.getClient().GetBalance(req.Address)
	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Error:   &common.Error{Code: 404},
			Balance: "0",
		}, err
	} else {
		return &wallet2.BalanceResponse{
			Error:   &common.Error{Code: 2000},
			Balance: balance,
		}, nil
	}
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	txs, err := a.getClient().GetTxByAddress(req.Address, req.Page, req.Pagesize)
	list := make([]*wallet2.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		list = append(list, &wallet2.TxMessage{
			Hash:   txs[i].TxHash,
			To:     []*wallet2.Address{{Address: txs[i].Dst}},
			From:   []*wallet2.Address{{Address: txs[i].Src}},
			Fee:    txs[i].TxHash,
			Status: true,
			Value:  []*wallet2.Value{{Value: string(rune(txs[i].Lamport))}},
			Type:   1,
			Height: string(rune(txs[i].Slot)),
		})
	}
	if err != nil {
		log.Error("get GetTxByAddress error", "err", err)
		return &wallet2.TxAddressResponse{
			Error: &common.Error{Code: 404},
			Tx:    nil,
		}, err
	} else {
		return &wallet2.TxAddressResponse{
			Error: &common.Error{Code: 2000},
			Tx:    list,
		}, nil
	}

}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	tx, err := a.getClient().GetTxByHash(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Error: &common.Error{Code: 404},
			Tx:    nil,
		}, err
	}
	return &wallet2.TxHashResponse{
		Tx: &wallet2.TxMessage{
			Hash:   tx.Hash,
			To:     []*wallet2.Address{{Address: tx.To}},
			From:   []*wallet2.Address{{Address: tx.From}},
			Fee:    tx.Fee,
			Status: tx.Status,
			Value:  []*wallet2.Value{{Value: tx.Value}},
			Type:   tx.Type,
			Height: tx.Height,
		},
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	address, _, err := a.getClient().GetAccount()
	if err != nil {
		log.Error("get GetAccount error", "err", err)
		return &wallet2.AccountResponse{
			Error:         &common.Error{Code: 404},
			AccountNumber: "",
		}, err
	} else {
		return &wallet2.AccountResponse{
			Error:         &common.Error{Code: 2000},
			AccountNumber: address,
		}, nil
	}
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	value, err := a.getClient().GetMinRent()
	if err != nil {
		log.Error("get GetMinRent error", "err", err)
		return &wallet2.MinRentResponse{
			Error: &common.Error{Code: 404},
			Value: "",
		}, err
	} else {
		return &wallet2.MinRentResponse{
			Error: &common.Error{Code: 2000},
			Value: value,
		}, nil
	}
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newSolanaClients(conf)
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

func (a *WalletAdaptor) getClient() *solanaClient {
	return a.clients.BestClient().(*solanaClient)
}

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	value, err := w.getClient().GetNonce()
	if err != nil {
		log.Error("get GetNonce error", "err", err)
		return &wallet2.NonceResponse{
			Error: &common.Error{Code: 404},
			Nonce: "",
		}, err
	} else {
		return &wallet2.NonceResponse{
			Error: &common.Error{Code: 2000},
			Nonce: value,
		}, nil
	}

}

func newWalletAdaptor(client *solanaClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func NewLocalWalletAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newWalletAdaptor(newLocalSolanaClient(network))
}

func (w *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return nil, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	return nil, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	value, err := w.getClient().SendTx("")
	if err != nil {
		log.Error("get GetNonce error", "err", err)
		return &wallet2.SendTxResponse{
			Error:  &common.Error{Code: 404},
			TxHash: "",
		}, err
	} else {
		return &wallet2.SendTxResponse{
			Error:  &common.Error{Code: 2000},
			TxHash: value,
		}, nil
	}

}
