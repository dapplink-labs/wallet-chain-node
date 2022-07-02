package service

import (
	"context"
	"github.com/SavourDao/savour-core/rpc/savourrpc/go-savourrpc/common"
	"github.com/SavourDao/savour-core/rpc/savourrpc/go-savourrpc/wallet"
)

type WalletRpcServer struct {
	wallet.WalletServiceServer
}

func (this *WalletRpcServer) GetSupportCoins(ctx context.Context, req *wallet.SupportCoinsRequest) (*wallet.SupportCoinsResponse, error) {
	return &wallet.SupportCoinsResponse{
		Error:        &common.Error{Code: 20000},
		SupportCoins: nil,
	}, nil
}

func (this *WalletRpcServer) GetNonce(ctx context.Context, req *wallet.NonceRequest) (*wallet.NonceResponse, error) {
	return &wallet.NonceResponse{
		Error: &common.Error{Code: 20000},
		Nonce: "",
	}, nil
}

func (this *WalletRpcServer) GetGasPrice(ctx context.Context, req *wallet.GasPriceRequest) (*wallet.GasPriceResponse, error) {
	return &wallet.GasPriceResponse{
		Error: &common.Error{Code: 20000},
		Gas:   "",
	}, nil
}

func (this *WalletRpcServer) SendTx(ctx context.Context, req *wallet.SendTxRequest) (*wallet.SendTxResponse, error) {
	return &wallet.SendTxResponse{
		Error:  &common.Error{Code: 20000},
		TxHash: "",
	}, nil
}
