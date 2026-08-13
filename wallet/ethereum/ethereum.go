package ethereum

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
	"github.com/the-web3/etherscan-api"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const (
	ChainName = "Ethereum"
	Coin      = "ETH"
)

var (
	Startblock = 0
	Endblock   = 999999999
)

type WalletAdaptor struct {
	clients      *multiclient.MultiClient
	etherscanCli *etherscan.Client
}

func (a *WalletAdaptor) GetBlockHeaderByHash(req *wallet2.BlockHeaderByHashRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetLatestSafeBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetLatestFinalizedBlockHeader(req *wallet2.BasicRequest) (*wallet2.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

type BlockResponse struct {
	Block *types.Block
	Err   error
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	return rpc.BlockNumber(number.Int64()).String()
}

// BatchGetBlocks 批量获取指定区块范围内的区块信息
func BatchGetBlocks(client *rpc.Client, startBlock, endBlock *big.Int) ([]BlockResponse, error) {
	// 计算区块数量
	blockCount := new(big.Int).Sub(endBlock, startBlock).Uint64() + 1
	// 创建请求数组和响应数组
	batchRequests := make([]rpc.BatchElem, blockCount)
	blockResponses := make([]BlockResponse, blockCount)
	// 遍历并构造批量请求
	for i := uint64(0); i < blockCount; i++ {
		blockNum := new(big.Int).Add(startBlock, big.NewInt(int64(i)))
		batchRequests[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []any{toBlockNumArg(blockNum), true},
			Result: new(types.Block),
		}
	}
	// 使用 BatchCallContext 发送批量请求
	err := client.BatchCallContext(context.Background(), batchRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to batch call: %v", err)
	}
	// 处理每个响应
	for i, req := range batchRequests {
		if req.Error != nil {
			blockResponses[i] = BlockResponse{Err: req.Error}
		} else {
			blockResponses[i] = BlockResponse{Block: req.Result.(*types.Block), Err: nil}
		}
	}
	return blockResponses, nil
}

func (a *WalletAdaptor) GetBlockByRange(req *wallet2.BlockByRangeRequest) (*wallet2.BlockByRangeResponse, error) {
	var startBlock, endBlock = stringToBigInt(req.Start), stringToBigInt(req.End)
	blockResponses, err := BatchGetBlocks(a.getClient().Client.Client(), startBlock, endBlock)
	if err != nil {
		log.Error("Error retrieving blocks :" + err.Error())
	}
	fmt.Println(blockResponses)
	return &wallet2.BlockByRangeResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		//Blocks: blocks,
	}, nil
}

func (a *WalletAdaptor) GetTxReceiptByHash(req *wallet2.TxReceiptByHashRequest) (*wallet2.TxReceiptByHashResponse, error) {
	receipt, err := a.getClient().TransactionReceipt(context.Background(), ethcommon.HexToHash(req.Hash))
	if err != nil {
		return nil, err
	}
	return &wallet2.TxReceiptByHashResponse{
		Code:              common.ReturnCode_SUCCESS,
		Msg:               "success",
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.String(),
		Status:            receipt.Status,
		TransactionIndex:  uint64(receipt.TransactionIndex),
		EffectiveGasPrice: receipt.EffectiveGasPrice.String(),
		GasUsed:           receipt.GasUsed,
	}, nil
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

func (a *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	return &wallet2.BlockResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil

}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newEthClients(conf)
	if err != nil {
		return nil, err
	}
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}
	return &WalletAdaptor{
		clients:      multiclient.New(clis),
		etherscanCli: NewEtherscanClient(conf.Fullnode.Eth.TpApiUrl, conf.Fullnode.Eth.TpApiKey),
	}, nil
}

func NewLocalWalletAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newWalletAdaptor(newLocalEthClient(network))
}

func newWalletAdaptor(client *ethClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *WalletAdaptor) getClient() *ethClient {
	return a.clients.BestClient().(*ethClient)
}

func stringToInt(amount string) *big.Int {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}

func (a *WalletAdaptor) blockNumber() *big.Int {
	return a.getClient().blockNumber()
}

func (a *WalletAdaptor) makeSigner() (types.Signer, error) {
	height := a.blockNumber()
	if height == nil {
		err := fmt.Errorf("fail to get height in making signer")
		return nil, err
	}
	log.Info("make signer", "height", height.Uint64())
	return types.MakeSigner(a.getClient().chainConfig, height, 1000), nil
}

func (a *WalletAdaptor) makeSignerOffline(height int64) types.Signer {
	if height == 0 {
		height = math.MaxInt64
	}
	return types.MakeSigner(a.getClient().chainConfig, big.NewInt(height), 1000)
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	var result *big.Int
	var err error
	if len(req.ContractAddress) > 0 {
		erc20Balance, err := a.etherscanCli.TokenBalance(req.ContractAddress, req.Address)
		if err != nil {
			return &wallet2.BalanceResponse{
				Code:    common.ReturnCode_ERROR,
				Msg:     "get erc20 balance fail",
				Balance: "0",
			}, err
		}
		result = (*big.Int)(erc20Balance)
	} else {
		result, err = a.getClient().BalanceAt(context.TODO(), ethcommon.HexToAddress(req.Address), nil)
	}
	if err != nil {
		log.Error("get balance error", "err", err)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}, err
	} else {
		// balanceCache.Add(key, result)
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "get balance success",
			Balance: result.String(),
		}, nil
	}
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	var tx_list []*wallet2.TxMessage
	if req.ContractAddress == "0x00" {
		txs, err := a.etherscanCli.NormalTxByAddress(req.Address, &Startblock, &Endblock, int(req.Page), int(req.Pagesize), true)
		if err != nil {
			return &wallet2.TxAddressResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}
		for _, ktx := range txs {
			if ktx.Value.Int().String() == "0" {
				continue
			}
			var from_addrs []*wallet2.Address
			var to_addrs []*wallet2.Address
			var value_list []*wallet2.Value
			var direction int32
			from_addrs = append(from_addrs, &wallet2.Address{Address: ktx.From})
			to_addrs = append(to_addrs, &wallet2.Address{Address: ktx.To})
			value_list = append(value_list, &wallet2.Value{Value: ktx.Value.Int().String()})
			bigIntGasUsed := int64(ktx.GasUsed)
			bigIntGasPrice := big.NewInt(ktx.GasPrice.Int().Int64())
			tx_fee := bigIntGasPrice.Int64() * bigIntGasUsed
			datetime := ktx.TimeStamp.Time().Format("2006-01-02 15:04:05")
			if strings.EqualFold(req.Address, ktx.From) {
				direction = 0 // 转出
			} else {
				direction = 1 // 转入
			}
			fmt.Println("ktx.Input===", ktx.Input)
			tx := &wallet2.TxMessage{
				Hash:            ktx.Hash,
				Froms:           from_addrs,
				Tos:             to_addrs,
				Values:          value_list,
				Fee:             strconv.FormatInt(tx_fee, 10),
				Status:          wallet2.TxStatus_Success,
				Type:            direction,
				Height:          strconv.Itoa(ktx.BlockNumber),
				ContractAddress: ktx.ContractAddress,
				Datetime:        datetime,
			}
			tx_list = append(tx_list, tx)
		}
	} else if req.ContractAddress == "0xSwap" {
		txs, err := a.etherscanCli.SwapTransactions(req.Address, &Startblock, &Endblock, int(req.Page), int(req.Pagesize), true)
		if err != nil {
			return &wallet2.TxAddressResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}
		for _, ktx := range txs {
			var from_addrs []*wallet2.Address
			var to_addrs []*wallet2.Address
			var value_list []*wallet2.Value
			var direction int32
			from_addrs = append(from_addrs, &wallet2.Address{Address: ktx.From})
			to_addrs = append(to_addrs, &wallet2.Address{Address: ktx.To})
			value_list = append(value_list, &wallet2.Value{Value: "0"})
			bigIntGasUsed := int64(ktx.GasUsed)
			bigIntGasPrice := big.NewInt(ktx.GasPrice.Int().Int64())
			tx_fee := bigIntGasPrice.Int64() * bigIntGasUsed
			datetime := ktx.TimeStamp.Time().Format("2006-01-02 15:04:05")
			direction = 3
			tx := &wallet2.TxMessage{
				Hash:            ktx.Hash,
				Froms:           from_addrs,
				Tos:             to_addrs,
				Values:          value_list,
				Fee:             strconv.FormatInt(tx_fee, 10),
				Status:          wallet2.TxStatus_Success,
				Type:            direction,
				Height:          strconv.Itoa(ktx.BlockNumber),
				ContractAddress: ktx.ContractAddress,
				Datetime:        datetime,
			}
			tx_list = append(tx_list, tx)
		}
	} else {
		txs, err := a.etherscanCli.ERC20Transfers(&req.ContractAddress, &req.Address, &Startblock, &Endblock, int(req.Page), int(req.Pagesize), true)
		if err != nil {
			return &wallet2.TxAddressResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}
		for _, ktx := range txs {
			fmt.Println("ktx:", ktx.Value.Int())
			var from_addrs []*wallet2.Address
			var to_addrs []*wallet2.Address
			var value_list []*wallet2.Value
			var direction int32
			from_addrs = append(from_addrs, &wallet2.Address{Address: ktx.From})
			to_addrs = append(to_addrs, &wallet2.Address{Address: ktx.To})
			value_list = append(value_list, &wallet2.Value{Value: ktx.Value.Int().String()})
			bigIntGasUsed := int64(ktx.GasUsed)
			bigIntGasPrice := big.NewInt(ktx.GasPrice.Int().Int64())
			tx_fee := bigIntGasPrice.Int64() * bigIntGasUsed
			datetime := ktx.TimeStamp.Time().Format("2006-01-02 15:04:05")
			if strings.EqualFold(req.Address, ktx.From) {
				direction = 0 // 转出
			} else {
				direction = 1 // 转入
			}
			tx := &wallet2.TxMessage{
				Hash:            ktx.Hash,
				Froms:           from_addrs,
				Tos:             to_addrs,
				Values:          value_list,
				Fee:             strconv.FormatInt(tx_fee, 10),
				Status:          wallet2.TxStatus_Success,
				Type:            direction,
				Height:          strconv.Itoa(ktx.BlockNumber),
				ContractAddress: ktx.ContractAddress,
				Datetime:        datetime,
			}
			tx_list = append(tx_list, tx)
		}
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get address tx success",
		Tx:   tx_list,
	}, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	tx, _, err := a.getClient().TransactionByHash(context.TODO(), ethcommon.HexToHash(req.Hash))
	if err != nil {
		if err == ethereum.NotFound {
			return &wallet2.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "Ethereum Tx NotFound",
			}, nil
		}
		log.Error("get transaction error", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Ethereum Tx NotFound",
		}, nil
	}
	receipt, err := a.getClient().TransactionReceipt(context.TODO(), ethcommon.HexToHash(req.Hash))
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Get transaction receipt error",
		}, nil
	}
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	from_addrs = append(from_addrs, &wallet2.Address{Address: ""})
	to_addrs = append(to_addrs, &wallet2.Address{Address: tx.To().Hex()})
	value_list = append(value_list, &wallet2.Value{Value: tx.Value().String()})
	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &wallet2.TxMessage{
			Hash:            tx.Hash().Hex(),
			Index:           uint32(receipt.TransactionIndex),
			Froms:           from_addrs,
			Tos:             to_addrs,
			Values:          value_list,
			Fee:             tx.GasFeeCap().String(),
			Status:          wallet2.TxStatus_Success,
			Type:            0,
			Height:          receipt.BlockNumber.String(),
			ContractAddress: tx.To().String(),
		},
	}, nil
}

func (a *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return &wallet2.SupportCoinsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "this coin support",
		Support: true,
	}, nil
}

func (a *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	var bockHeight *big.Int
	nonce, err := a.getClient().NonceAt(context.TODO(), ethcommon.HexToAddress(req.Address), bockHeight)
	if err != nil {
		log.Error("get nonce failed", "err", err)
		return &wallet2.NonceResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   "get nonce failed",
			Nonce: strconv.FormatUint(nonce, 10),
		}, nil
	}
	return &wallet2.NonceResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "get nonce success",
		Nonce: strconv.FormatUint(nonce, 10),
	}, nil
}

func (a *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	price, err := a.getClient().SuggestGasPrice(context.TODO())
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &wallet2.GasPriceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get gas price fail",
			Gas:  "",
		}, nil
	}
	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get gas price success",
		Gas:  price.String(),
	}, nil
}

func (a *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	txbytes, err := hexutil.Decode(req.RawTx)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "Send tx fail(Decode)",
			TxHash: "",
		}, err
	}
	txSigned := new(types.Transaction)
	if err := rlp.DecodeBytes(txbytes, txSigned); err != nil {
		log.Error("signedTx DecodeBytes failed", "err", err)
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "Send tx fail(DecodeBytes)",
			TxHash: "",
		}, err
	}
	log.Info("broadcast tx", "tx", hexutil.Encode([]byte(req.RawTx)))
	txHash := fmt.Sprintf("0x%x", txSigned.Hash())
	if err := a.getClient().SendTransaction(context.TODO(), txSigned); err != nil {
		log.Error("braoadcast tx failed", "tx_hash", txHash, "err", err)
		return &wallet2.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    err.Error(),
			TxHash: "",
		}, err
	}
	log.Info("braoadcast tx success", "tx_hash", txHash)
	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "Send and braoadcast tx success",
		TxHash: txHash,
	}, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	publicKey, err := btcec.ParsePubKey(req.PublicKey)
	if err != nil {
		log.Error(" btcec.ParsePubKey failed", "err", err)
		return &wallet2.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	log.Info("convert req pub to address", "address", crypto.PubkeyToAddress(*publicKey.ToECDSA()).String())

	return &wallet2.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Address: crypto.PubkeyToAddress(*publicKey.ToECDSA()).String(),
	}, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	valid := ethcommon.IsHexAddress(req.Address)
	stdAddr := ethcommon.HexToAddress(req.Address)
	log.Info("valid address", "address", req.Address, "valid", valid, "standardAddreess", stdAddr.String())

	isContract := false
	if !a.getClient().local {
		isContract = a.getClient().isContractAddress(stdAddr)
	}
	return &wallet2.ValidAddressResponse{
		Code:             common.ReturnCode_SUCCESS,
		Valid:            valid,
		CanWithdrawal:    !isContract,
		CanonicalAddress: stdAddr.String(),
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	rawTx := new(types.Transaction)
	if err := rlp.DecodeBytes(req.RawData, rawTx); err != nil {
		log.Error("signedTx unmarlshal failed", "err", err)
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	reply, err := a.queryRawTransaction(req.Chain != req.Symbol, rawTx)
	if err != nil {
		log.Error("queryRawTransaction failed", "err", err)
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	signer := a.makeSignerOffline(req.Height)
	reply.SignHash = signer.Hash(rawTx).Bytes()
	return reply, nil
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	signedTx := new(types.Transaction)
	if err := rlp.DecodeBytes(req.SignedTxData, signedTx); err != nil {
		log.Error("signedTx unmarlshal failed", "err", err)
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  err.Error(),
		}, err
	}
	return a.queryTransaction(req.Chain != req.Symbol, signedTx, nil, 0, a.makeSignerOffline(req.Height))
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(req.TxData, tx); err != nil {
		log.Error("tx unmarlshal failed", "err", err)
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	signer, err := a.makeSigner()
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, req.Signature)
	if err != nil {
		log.Error("tx WithSignature failed", "err", err)
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	signedTxData, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Error("signedTx EncodeToBytes failed", "err", err)
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &wallet2.CreateSignedTxResponse{
		Code:         common.ReturnCode_SUCCESS,
		SignedTxData: signedTxData,
		Hash:         signedTx.Hash().Bytes(),
	}, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	if !ethcommon.IsHexAddress(req.From) {
		log.Info("invalid from address", "from", req.From)
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid from address",
		}, nil
	}

	if !ethcommon.IsHexAddress(req.To) {
		log.Info("invalid to address", "to", req.To)
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid to address",
		}, nil
	}

	if req.Amount == "" {
		log.Info("amount can not be zero", "amount", req.Amount)
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "zero amount",
		}, nil
	}

	if req.GasLimit == "" {
		log.Info("gas uints can not be zero", "gas_unit", req.GasLimit)
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "zero gasunits",
		}, nil
	}

	if req.GasPrice == "" {
		log.Info("gas price can not be zero", "gas_price", req.GasPrice)
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "zero gasprice",
		}, nil
	}

	// get nonce from chain
	nonce := req.Nonce
	// calc amount
	assetAmount := stringToInt(req.Amount)
	if assetAmount == nil {
		log.Error("convert asset amount failed")
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "convert asset amount failed",
		}, nil
	}

	// convert gas price
	gasPrice := stringToInt(req.GasPrice)
	if gasPrice == nil {
		log.Error("convert gasPrice failed")
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "convert gasPrice failed",
		}, nil
	}

	gasLimit := stringToInt(req.GasLimit)
	if gasLimit == nil {
		log.Error("convert gasLimit failed")
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "convert gasLimit failed",
		}, nil
	}

	var err error
	// make transaction
	var tx *types.Transaction
	if len(req.ContractAddress) > 0 {
		tx, err = a.getClient().Erc20RawTransfer(req.ContractAddress, nonce, ethcommon.HexToAddress(req.To), assetAmount,
			gasLimit.Uint64(), gasPrice)
		if err != nil {
			log.Error("ERC20 tx raw transfer failed", "err", err)
			return nil, err
		}
	} else {
		tx = types.NewTransaction(nonce, ethcommon.HexToAddress(req.To), assetAmount, gasLimit.Uint64(), gasPrice, nil)
	}
	txData, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Error("tx EncodeToBytes failed", "err", err)
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	signer, err := a.makeSigner()
	if err != nil {
		return nil, err
	}

	return &wallet2.CreateAccountTxResponse{
		Code:     common.ReturnCode_SUCCESS,
		TxData:   txData,
		SignHash: signer.Hash(tx).Bytes(),
	}, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	signedTx := new(types.Transaction)
	if err := rlp.DecodeBytes(req.SignedTxData, signedTx); err != nil {
		log.Error("signedTx DecodeBytess failed", "err", err)
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	signer := a.makeSignerOffline(req.Height)
	sender, err := signer.Sender(signedTx)
	if err != nil {
		log.Error("failed to get sender from signed tx", "err", err)
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	var expectedSender string
	if len(req.Addresses) > 0 {
		expectedSender = req.Addresses[0]
	} else {
		expectedSender = req.Sender
	}
	return &wallet2.VerifySignedTxResponse{
		Code:     common.ReturnCode_SUCCESS,
		Verified: sender == ethcommon.HexToAddress(expectedSender),
	}, err
}

// queryTransaction retrieve transaction information from a signed data.
func (a *WalletAdaptor) queryTransaction(isERC20 bool, tx *types.Transaction, receipt *types.Receipt, blockNumber uint64, signer types.Signer) (*wallet2.AccountTxResponse, error) {
	reply, err := a.queryRawTransaction(isERC20, tx)
	if err != nil {
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	reply.SignHash = signer.Hash(tx).Bytes()
	gasUsed := new(big.Int)
	if receipt != nil {
		gasUsed = gasUsed.SetUint64(receipt.GasUsed).Mul(gasUsed, tx.GasPrice())
		if isERC20 {
			// Check ERC20 Transfer event log
			err := a.validateAndQueryERC20TransferReceipt(ethcommon.HexToAddress(reply.ContractAddress),
				reply.From, reply.To, reply.Amount, receipt)
			if err != nil {
				return &wallet2.AccountTxResponse{
					Code: common.ReturnCode_ERROR,
					Msg:  err.Error(),
				}, nil
			}
		}
	}

	log.Info("QueryTransaction", "from", reply.From,
		"block_number", blockNumber,
		"gas_used", decimal.NewFromBigInt(gasUsed, 0).String())
	reply.From = reply.From
	reply.CostFee = decimal.NewFromBigInt(gasUsed, 0).String()
	reply.BlockHeight = blockNumber
	reply.Status = wallet2.TxStatus_Success
	reply.TxHash = tx.Hash().String()
	return reply, nil
}

// queryRawTransaction retrieve transaction information from a raw(unsigned) data.
func (a *WalletAdaptor) queryRawTransaction(isERC20 bool, rawTx *types.Transaction) (*wallet2.AccountTxResponse, error) {
	var amount *big.Int
	var to ethcommon.Address
	contractAddress := ""
	var err error
	if isERC20 {
		// erc20 transfer transaction
		to, amount, err = a.validateAndQueryERC20RawTransfer(*rawTx.To(), rawTx)
		if err != nil {
			return nil, err
		}
		contractAddress = rawTx.To().String()
	} else {
		// ether transaction
		to = *rawTx.To()
		amount = rawTx.Value()
	}
	log.Info("QueryRawTransaction",
		"is_erc20", isERC20,
		"to", to.String(),
		"nonce", rawTx.Nonce(),
		"amount", decimal.NewFromBigInt(amount, 0).String(),
		"gas_limit", decimal.NewFromBigInt(big.NewInt(int64(rawTx.Gas())), 0).String(),
		"gas_price", decimal.NewFromBigInt(rawTx.GasPrice(), 0).String())

	return &wallet2.AccountTxResponse{
		Code:            common.ReturnCode_SUCCESS,
		To:              to.String(),
		Nonce:           rawTx.Nonce(),
		Amount:          decimal.NewFromBigInt(amount, 0).String(),
		GasLimit:        decimal.NewFromBigInt(big.NewInt(int64(rawTx.Gas())), 0).String(),
		GasPrice:        decimal.NewFromBigInt(rawTx.GasPrice(), 0).String(),
		ContractAddress: contractAddress,
	}, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	return &wallet2.UtxoInsResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return &wallet2.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	return &wallet2.CreateUtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return &wallet2.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	return &wallet2.AccountResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func stringToBigInt(s string) (num *big.Int) {
	if s != "" {
		n, bol := big.NewInt(0).SetString(s, 10)
		if bol == true {
			num = n
		}
	}
	return
}

// GetBlockHeaderByNumber 根据区块号获取区块头
func (a *WalletAdaptor) GetBlockHeaderByNumber(req *wallet2.BlockHeaderRequest) (*wallet2.BlockHeaderResponse, error) {
	header, err := a.getClient().HeaderByNumber(context.Background(), stringToBigInt(req.GetHeight()))
	if err != nil {
		return &wallet2.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	return &wallet2.BlockHeaderResponse{
		Code:            common.ReturnCode_SUCCESS,
		Msg:             "get block header by number success",
		ParentHash:      header.ParentHash.Hex(),
		UncleHash:       header.UncleHash.Hex(),
		Coinbase:        header.Coinbase.Hex(),
		Root:            header.Root.Hex(),
		TxHash:          header.TxHash.Hex(),
		ReceiptHash:     header.ReceiptHash.Hex(),
		Number:          header.Number.String(),
		Difficulty:      header.Difficulty.String(),
		GasLimit:        header.GasLimit,
		GasUsed:         header.GasUsed,
		Time:            header.Time,
		MixDigest:       header.MixDigest.Hex(),
		BaseFee:         header.BaseFee.String(),
		WithdrawalsHash: header.WithdrawalsHash.Hex(),
	}, nil
}

// GetBlockByNumber 根据区块号获取区块
func (a *WalletAdaptor) GetBlockByNumber(req *wallet2.BlockInfoRequest) (*wallet2.BlockInfoResponse, error) {
	block, err := a.getClient().BlockByNumber(context.Background(), stringToBigInt(req.GetHeight()))
	if err != nil {
		return &wallet2.BlockInfoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	var transactions []*wallet2.BlockInfoTransactionList
	for _, tx := range block.Transactions() {
		signer := types.LatestSignerForChainID(tx.ChainId())
		from, err := types.Sender(signer, tx)
		if err != nil {
			log.Error("Failed to get sender from transaction: %v", err)
		}
		toAddress := ""
		if tx.To() != nil {
			toAddress = tx.To().Hex()
		}
		transactions = append(transactions, &wallet2.BlockInfoTransactionList{
			Hash:   tx.Hash().Hex(),
			To:     toAddress,
			From:   from.Hex(),
			Time:   tx.Time().String(),
			Amount: tx.Value().String(),
		})
	}
	return &wallet2.BlockInfoResponse{
		Code:         common.ReturnCode_SUCCESS,
		Msg:          "get block by number success",
		Hash:         block.Hash().Hex(),
		Transactions: transactions,
		BaseFee:      block.BaseFee().String(),
	}, nil
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	return &wallet2.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	return &wallet2.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}

type semaphore chan struct{}

func (s semaphore) Acquire() {
	s <- struct{}{}
}

func (s semaphore) Release() {
	<-s
}
