package solana

import (
	"github.com/SavourDao/savour-core/config"
	wallet2 "github.com/SavourDao/savour-core/rpc/wallet"
	"github.com/SavourDao/savour-core/wallet"
	"github.com/SavourDao/savour-core/wallet/fallback"
	"github.com/SavourDao/savour-core/wallet/multiclient"
)

const (
	ChainName = "SOL"
	Symbol    = "SOL"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *multiclient.MultiClient
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance := a.getClient().GetBalance(req.Address)
	return &wallet2.BalanceResponse{
		Balance: balance,
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	txs := a.getClient().GetTxByAddress(req.Address)
	list := make([]*wallet2.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		list = append(list, &wallet2.TxMessage{
			Hash:   txs[i].TxHash,
			To:     txs[i].Dst,
			From:   txs[i].Src,
			Fee:    txs[i].TxHash,
			Status: true,
			Value:  string(rune(txs[i].Lamport)),
			Type:   1,
			Height: string(txs[i].Slot),
		})
	}
	return &wallet2.TxAddressResponse{
		Tx: list,
	}, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	_ = a.getClient().GetTxByHash(req.Hash)
	return &wallet2.TxHashResponse{
		Tx: &wallet2.TxMessage{
			//Hash:   tx.Result,
			//To:     tx.Dst,
			//From:   txs[i].Src,
			//Fee:    txs[i].TxHash,
			//Status: true,
			//Value:  string(rune(txs[i].Lamport)),
			//Type:   1,
			//Height: string(txs[i].Slot),
		},
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	adderss, _, _ := a.getClient().GetAccount()
	return &wallet2.AccountResponse{
		AccountNumber: adderss,
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	value := a.getClient().GetMinRent()
	return &wallet2.MinRentResponse{
		Value: value,
	}, nil
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

//func NewLocalWalletAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
//	return newWalletAdaptor(newLocalEthClient(network))
//}

func newWalletAdaptor(client *solanaClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *WalletAdaptor) getClient() *solanaClient {
	return a.clients.BestClient().(*solanaClient)
}

func (w *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return nil, nil
}

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	value := w.getClient().GetNonce(req.Address)
	return &wallet2.NonceResponse{
		Nonce: value,
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	return nil, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	value := w.getClient().SendTx()
	return &wallet2.SendTxResponse{
		TxHash: value,
	}, nil
}
