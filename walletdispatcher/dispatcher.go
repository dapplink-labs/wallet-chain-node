package walletdispatcher

import (
	"context"
	"github.com/SavourDao/savour-hd/rpc/common"
	"github.com/SavourDao/savour-hd/wallet"
	"github.com/SavourDao/savour-hd/wallet/solana"
	"runtime/debug"
	"strings"

	"github.com/SavourDao/savour-hd/config"
	wallet2 "github.com/SavourDao/savour-hd/rpc/wallet"
	"github.com/SavourDao/savour-hd/wallet/bitcoin"
	"github.com/SavourDao/savour-hd/wallet/ethereum"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommonRequest interface {
	GetChain() string
}

type CommonReply = wallet2.SupportCoinsResponse

type ChainType = string

type WalletDispatcher struct {
	registry map[ChainType]wallet.WalletAdaptor
}

func (d *WalletDispatcher) mustEmbedUnimplementedWalletServiceServer() {
	//TODO implement me
	panic("implement me")
}

func New(conf *config.Config) (*WalletDispatcher, error) {
	dispatcher := WalletDispatcher{
		registry: make(map[ChainType]wallet.WalletAdaptor),
	}
	walletAdaptorFactoryMap := map[string]func(conf *config.Config) (wallet.WalletAdaptor, error){
		bitcoin.ChainName:  bitcoin.NewChainAdaptor,
		ethereum.ChainName: ethereum.NewChainAdaptor,
		solana.ChainName:   solana.NewChainAdaptor,
	}
	supportedChains := []string{bitcoin.ChainName, ethereum.ChainName, solana.ChainName}
	for _, c := range conf.Chains {
		if factory, ok := walletAdaptorFactoryMap[c]; ok {
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("failed to setup chain", "chain", c, "error", err)
			}
			dispatcher.registry[c] = adaptor
		} else {
			log.Error("unsupported chain", "chain", c, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

func NewLocal(network config.NetWorkType) *WalletDispatcher {
	dispatcher := WalletDispatcher{
		registry: make(map[ChainType]wallet.WalletAdaptor),
	}

	walletAdaptorFactoryMap := map[string]func(network config.NetWorkType) wallet.WalletAdaptor{
		bitcoin.ChainName:  bitcoin.NewLocalChainAdaptor,
		ethereum.ChainName: ethereum.NewLocalWalletAdaptor,
	}
	supportedChains := []string{bitcoin.ChainName, ethereum.ChainName}

	for _, c := range supportedChains {
		if factory, ok := walletAdaptorFactoryMap[c]; ok {
			dispatcher.registry[c] = factory(network)
		}
	}
	return &dispatcher
}

func (d *WalletDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	chain := req.(CommonRequest).GetChain()
	log.Info(method, "chain", chain, "req", req)

	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *WalletDispatcher) preHandler(req interface{}) (resp *CommonReply) {
	chain := req.(CommonRequest).GetChain()
	if _, ok := d.registry[chain]; !ok {
		return &CommonReply{
			Error:   &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			Support: false,
		}
	}
	return nil
}

func (d *WalletDispatcher) GetSupportCoins(ctx context.Context, request *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *WalletDispatcher) GetNonce(ctx context.Context, request *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.NonceResponse{
			Error: &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
		}, nil
	}
	return d.registry[request.Chain].GetNonce(request)
}

func (d *WalletDispatcher) GetGasPrice(ctx context.Context, request *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.GasPriceResponse{
			Error: &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			Gas:   "",
		}, nil
	}
	return d.registry[request.Chain].GetGasPrice(request)
}

func (d *WalletDispatcher) SendTx(ctx context.Context, request *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.SendTxResponse{
			Error:  &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			TxHash: "",
		}, nil
	}
	return d.registry[request.Chain].SendTx(request)
}

func (d *WalletDispatcher) GetBalance(ctx context.Context, request *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.BalanceResponse{
			Error:   &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			Balance: "",
		}, nil
	}
	return d.registry[request.Chain].GetBalance(request)
}

func (d *WalletDispatcher) GetTxByAddress(ctx context.Context, request *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.TxAddressResponse{
			Error: &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			Tx:    nil,
		}, nil
	}
	return d.registry[request.Chain].GetTxByAddress(request)
}

func (d *WalletDispatcher) GetTxByHash(ctx context.Context, request *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.TxHashResponse{
			Error: &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			Tx:    nil,
		}, nil
	}
	return d.registry[request.Chain].GetTxByHash(request)
}

func (d *WalletDispatcher) GetAccount(ctx context.Context, request *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.AccountResponse{
			Error:         &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			AccountNumber: "",
			Sequence:      "",
		}, nil
	}
	return d.registry[request.Chain].GetAccount(request)
}

func (d *WalletDispatcher) GetUtxo(ctx context.Context, request *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.UtxoResponse{
			Error: &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
		}, nil
	}
	return d.registry[request.Chain].GetUtxo(request)
}

func (d *WalletDispatcher) GetMinRent(ctx context.Context, request *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.MinRentResponse{
			Error: &common.Error{Code: common.ReturnCode_ERROR, Brief: config.UnsupportedOperation, Detail: config.UnsupportedChain, CanRetry: true},
			Value: "",
		}, nil
	}
	return d.registry[request.Chain].GetMinRent(request)
}
