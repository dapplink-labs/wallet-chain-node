package walletdispatcher

import (
	"context"
	"github.com/savour-labs/wallet-chain-node/wallet/arweave"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/ada"
	"github.com/savour-labs/wallet-chain-node/wallet/arbitrum"
	"github.com/savour-labs/wallet-chain-node/wallet/avalanche"
	"github.com/savour-labs/wallet-chain-node/wallet/base"
	"github.com/savour-labs/wallet-chain-node/wallet/binance"
	"github.com/savour-labs/wallet-chain-node/wallet/bitcoin"
	"github.com/savour-labs/wallet-chain-node/wallet/eosio"
	"github.com/savour-labs/wallet-chain-node/wallet/ethereum"
	"github.com/savour-labs/wallet-chain-node/wallet/evmos"
	"github.com/savour-labs/wallet-chain-node/wallet/heco"
	"github.com/savour-labs/wallet-chain-node/wallet/linea"
	"github.com/savour-labs/wallet-chain-node/wallet/mantle"
	"github.com/savour-labs/wallet-chain-node/wallet/near"
	"github.com/savour-labs/wallet-chain-node/wallet/optimism"
	"github.com/savour-labs/wallet-chain-node/wallet/polygon"
	"github.com/savour-labs/wallet-chain-node/wallet/solana"
	"github.com/savour-labs/wallet-chain-node/wallet/tron"
	"github.com/savour-labs/wallet-chain-node/wallet/xrp"
	"github.com/savour-labs/wallet-chain-node/wallet/zksync"
)

type CommonRequest interface {
	GetChain() string
}

type CommonReply = wallet2.SupportCoinsResponse

type ChainType = string

type WalletDispatcher struct {
	registry map[ChainType]wallet.WalletAdaptor
}

func New(conf *config.Config) (*WalletDispatcher, error) {
	dispatcher := WalletDispatcher{
		registry: make(map[ChainType]wallet.WalletAdaptor),
	}
	walletAdaptorFactoryMap := map[string]func(conf *config.Config) (wallet.WalletAdaptor, error){
		bitcoin.ChainName:   bitcoin.NewChainAdaptor,
		ethereum.ChainName:  ethereum.NewChainAdaptor,
		solana.ChainName:    solana.NewChainAdaptor,
		arbitrum.ChainName:  arbitrum.NewChainAdaptor,
		base.ChainName:      base.NewChainAdaptor,
		linea.ChainName:     linea.NewChainAdaptor,
		mantle.ChainName:    mantle.NewChainAdaptor,
		tron.ChainName:      tron.NewWalletAdaptor,
		zksync.ChainName:    zksync.NewChainAdaptor,
		optimism.ChainName:  optimism.NewChainAdaptor,
		polygon.ChainName:   polygon.NewChainAdaptor,
		binance.ChainName:   binance.NewChainAdaptor,
		heco.ChainName:      heco.NewChainAdaptor,
		avalanche.ChainName: avalanche.NewChainAdaptor,
		evmos.ChainName:     evmos.NewChainAdaptor,
		near.ChainName:      near.NewChainAdaptor,
		xrp.ChainName:       xrp.NewChainAdaptor,
		eosio.ChainName:     eosio.NewChainAdaptor,
		ada.ChainName:       ada.NewChainAdaptor,
		arweave.ChainName:   arweave.NewChainAdaptor,
	}
	supportedChains := []string{
		bitcoin.ChainName, ethereum.ChainName, solana.ChainName, arbitrum.ChainName, base.ChainName, linea.ChainName,
		mantle.ChainName, tron.ChainName, zksync.ChainName, optimism.ChainName, polygon.ChainName, binance.ChainName,
		heco.ChainName, avalanche.ChainName, evmos.ChainName, near.ChainName, xrp.ChainName, eosio.ChainName, ada.ChainName,
		arweave.ChainName,
	}
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
		bitcoin.ChainName:   bitcoin.NewLocalChainAdaptor,
		ethereum.ChainName:  ethereum.NewLocalWalletAdaptor,
		solana.ChainName:    solana.NewLocalWalletAdaptor,
		arbitrum.ChainName:  arbitrum.NewLocalWalletAdaptor,
		zksync.ChainName:    zksync.NewLocalWalletAdaptor,
		optimism.ChainName:  optimism.NewLocalWalletAdaptor,
		polygon.ChainName:   polygon.NewLocalWalletAdaptor,
		binance.ChainName:   binance.NewLocalWalletAdaptor,
		heco.ChainName:      heco.NewLocalWalletAdaptor,
		avalanche.ChainName: avalanche.NewLocalWalletAdaptor,
		evmos.ChainName:     evmos.NewLocalWalletAdaptor,
	}
	supportedChains := []string{
		bitcoin.ChainName, ethereum.ChainName, solana.ChainName, arbitrum.ChainName,
		zksync.ChainName, optimism.ChainName, polygon.ChainName, binance.ChainName,
		heco.ChainName, avalanche.ChainName, evmos.ChainName,
	}
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
			Code:    common.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Support: false,
		}
	}
	return nil
}

func (d *WalletDispatcher) GetSupportCoins(ctx context.Context, request *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.SupportCoinsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[request.Chain].GetSupportCoins(request)
}

func (d *WalletDispatcher) GetBlock(ctx context.Context, request *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block number fail",
		}, nil
	}
	return d.registry[request.Chain].GetBlock(request)
}

func (d *WalletDispatcher) GetNonce(ctx context.Context, request *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.NonceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[request.Chain].GetNonce(request)
}

func (d *WalletDispatcher) GetGasPrice(ctx context.Context, request *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.GasPriceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
			Gas:  "",
		}, nil
	}
	return d.registry[request.Chain].GetGasPrice(request)
}

func (d *WalletDispatcher) SendTx(ctx context.Context, request *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    config.UnsupportedOperation,
			TxHash: "",
		}, nil
	}
	return d.registry[request.Chain].SendTx(request)
}

func (d *WalletDispatcher) GetBalance(ctx context.Context, request *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Balance: "",
		}, nil
	}
	return d.registry[request.Chain].GetBalance(request)
}

func (d *WalletDispatcher) GetTxByAddress(ctx context.Context, request *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
			Tx:   nil,
		}, nil
	}
	return d.registry[request.Chain].GetTxByAddress(request)
}

func (d *WalletDispatcher) GetTxByHash(ctx context.Context, request *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
			Tx:   nil,
		}, nil
	}
	return d.registry[request.Chain].GetTxByHash(request)
}

func (d *WalletDispatcher) GetAccount(ctx context.Context, request *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.AccountResponse{
			Code:          common.ReturnCode_ERROR,
			Msg:           config.UnsupportedOperation,
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
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[request.Chain].GetUtxo(request)
}

func (d *WalletDispatcher) GetUnspentOutputs(ctx context.Context, request *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.UnspentOutputsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[request.Chain].GetUnspentOutputs(request)
}

func (d *WalletDispatcher) GetMinRent(ctx context.Context, request *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.MinRentResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   config.UnsupportedOperation,
			Value: "",
		}, nil
	}
	return d.registry[request.Chain].GetMinRent(request)
}

func (d *WalletDispatcher) ConvertAddress(ctx context.Context, request *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].ConvertAddress(request)
}

func (d *WalletDispatcher) ValidAddress(ctx context.Context, request *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].ValidAddress(request)
}

func (d *WalletDispatcher) GetUtxoInsFromData(ctx context.Context, request *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.UtxoInsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].GetUtxoInsFromData(request)
}

func (d *WalletDispatcher) GetAccountTxFromData(ctx context.Context, request *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].GetAccountTxFromData(request)
}

func (d *WalletDispatcher) GetUtxoTxFromData(ctx context.Context, request *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].GetUtxoTxFromData(request)
}

func (d *WalletDispatcher) GetAccountTxFromSignedData(ctx context.Context, request *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].GetAccountTxFromSignedData(request)
}

func (d *WalletDispatcher) GetUtxoTxFromSignedData(ctx context.Context, request *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].GetUtxoTxFromSignedData(request)
}

func (d *WalletDispatcher) CreateAccountSignedTx(ctx context.Context, request *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].CreateAccountSignedTx(request)
}

func (d *WalletDispatcher) CreateAccountTx(ctx context.Context, request *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].CreateAccountTx(request)
}

func (d *WalletDispatcher) CreateUtxoSignedTx(ctx context.Context, request *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].CreateUtxoSignedTx(request)
}

func (d *WalletDispatcher) CreateUtxoTx(ctx context.Context, request *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].CreateUtxoTx(request)
}

func (d *WalletDispatcher) VerifyAccountSignedTx(ctx context.Context, request *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].VerifyAccountSignedTx(request)
}

func (d *WalletDispatcher) VerifyUtxoSignedTx(ctx context.Context, request *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].VerifyUtxoSignedTx(request)
}

func (d *WalletDispatcher) ABIBinToJSON(ctx context.Context, request *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.ABIBinToJSONResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].ABIBinToJSON(request)
}

func (d *WalletDispatcher) ABIJSONToBin(ctx context.Context, request *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet2.ABIJSONToBinResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain,
		}, nil
	}
	return d.registry[request.Chain].ABIJSONToBin(request)
}
