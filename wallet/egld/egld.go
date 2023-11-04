package egld

//
//import (
//	"github.com/savour-labs/wallet-chain-node/config"
//	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
//	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
//	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
//)
//
//const (
//	ChainName = "EGLD"
//	Coin      = "egld"
//)
//
//type WalletAdaptor struct {
//	fallback.WalletAdaptor
//	clients *multiclient.MultiClient
//}
//
//func NewChainAdaptor(conf *config.Config) (*WalletAdaptor, error) {
//	clients, err := newElgdClient(conf)
//	if err != nil {
//		return nil, err
//	}
//	multiClients := make([]multiclient.Client, len(clients))
//	for i, client := range clients {
//		multiClients[i] = client
//	}
//	return &WalletAdaptor{
//		clients: multiclient.New(multiClients),
//	}, nil
//}
//
//func (wa *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
//	//TODO implement me
//	client := wa.clients.BestClient().(EgldClient)
//	account := client.GetAccountBalance(req.Address)
//	balance := account.Balance
//	return &wallet2.BalanceResponse{Balance: balance}, nil
//}
//
//func (wa *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
//
//	//TODO implement me
//	panic("implement me")
//}
//
//func (wa *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
