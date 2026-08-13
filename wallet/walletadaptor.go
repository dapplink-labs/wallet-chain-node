package wallet

import (
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
)

type WalletAdaptor interface {
	GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error)
	GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error)

	ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error)
	ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error)

	GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error)
	GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error)
	SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error)
	GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error)
	GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error)
	GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error)
	GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error)
	GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error)
	GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error)
	GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error)

	GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error)
	GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error)
	GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error)
	GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error)
	GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error)

	CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error)
	CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error)
	CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error)
	CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error)

	VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error)
	VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error)

	ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error)
	ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error)

	// GetBlockHeaderByNumber 根据区块高度获取区块头
	GetBlockHeaderByNumber(req *wallet2.BlockHeaderRequest) (*wallet2.BlockHeaderResponse, error)
	// GetBlockByNumber 根据区块高度获取区块
	GetBlockByNumber(req *wallet2.BlockInfoRequest) (*wallet2.BlockInfoResponse, error)
	// GetLatestSafeBlockHeader 获取最新安全区块
	GetLatestSafeBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error)
	// GetLatestFinalizedBlockHeader 获取最新最终确认区块
	GetLatestFinalizedBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error)
	// GetBlockHeaderByHash 根据区块hash获取区块头
	GetBlockHeaderByHash(req *wallet2.BlockHeaderByHashRequest) (*wallet2.BlockHeaderResponse, error)
	// GetBlockByRange 根据区块高度范围获取区块头
	GetBlockByRange(req *wallet2.BlockByRangeRequest) (*wallet2.BlockByRangeResponse, error)
	// GetTxReceiptByHash 根据交易hash获取交易收据
	GetTxReceiptByHash(req *wallet2.TxReceiptByHashRequest) (*wallet2.TxReceiptByHashResponse, error)
	// GetStorageHash 获取store hash
	GetStorageHash(req *wallet2.StorageHashRequest) (*wallet2.StorageHashResponse, error)
	// GetFilterLogs 获取区块内的日志
	GetFilterLogs(req *wallet2.FilterLogsRequest) (*wallet2.FilterLogsResponse, error)
	// GetTxCountByAddress 获取指定地址的交易数量
	GetTxCountByAddress(req *wallet2.TxCountByAddressRequest) (*wallet2.TxCountByAddressResponse, error)
	// GetSuggestGasPrice 获取建议的gas价格
	GetSuggestGasPrice(req *wallet2.SuggestGasPriceRequest) (*wallet2.SuggestGasPriceResponse, error)
	// GetSuggestGasTipCap 获取建议的gas费
	GetSuggestGasTipCap(req *wallet2.SuggestGasPriceRequest) (*wallet2.SuggestGasPriceResponse, error)
}
