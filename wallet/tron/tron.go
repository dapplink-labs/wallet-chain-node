package tron

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	pb "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const TrxDecimals = 6

const (
	ChainName               = "Tron"
	TronSymbol              = "TRX"
	MaxTimeUntillExpiration = 24*60*60*1000 - 120000 //23hour58min, MaxTimeUntillExpiration is 24 hours in Tron
)

const (
	trc20TransferTopicLen        = 3
	trc20TransferTopic           = "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	trc20TransferAddrLen         = 32
	trc20TransferMethodSignature = "a9059cbb"
	defaultGasLimit              = 1000000
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients  *multiclient.MultiClient
	tronScan *TronScan
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

func NewWalletAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newTronClients(conf)
	if err != nil {
		log.Error("new tron client err", "err", err)
		return nil, err
	}
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}

	tronScan, err := NewTronScanClient(conf.OkLink.OkLinkBaseUrl, conf.OkLink.OkLinkApiKey, conf.OkLink.OkLinkTimeout)
	if err != nil {
		return nil, err
	}
	return &WalletAdaptor{
		clients:  multiclient.New(clis),
		tronScan: tronScan,
	}, nil
}

func NewLocalWalletAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newWalletAdaptor(newLocalTronClient(network))
}

func newWalletAdaptor(client *tronClient) wallet.WalletAdaptor {
	return &WalletAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *WalletAdaptor) getClient() *tronClient {
	return a.clients.BestClient().(*tronClient)
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	log.Info("GetBalance", "req", req)

	grpcClient := a.getClient().grpcClient

	var result *big.Int
	if req.ContractAddress != "0x00" {
		symbol, err := grpcClient.TRC20GetSymbol(req.ContractAddress)
		if err != nil {
			return &wallet2.BalanceResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "get balance fail",
			}, err
		}
		log.Info("Get symbol success", "symbol", symbol)
		result, err = grpcClient.TRC20ContractBalance(req.Address, req.ContractAddress)
		if err != nil {
			return &wallet2.BalanceResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "get balance fail",
			}, err
		}
	} else {
		acc, err := grpcClient.GetAccount(req.Address)
		if err != nil {
			fmt.Println("sssss", err)
			return &wallet2.BalanceResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "get balance fail",
			}, err
		}

		if req.Coin == TronSymbol {
			//TRX
			result = big.NewInt(acc.Balance)
		} else {
			//TRC10
			if r, exist := acc.AssetV2[req.Coin]; !exist {
				result = big.NewInt(0)
			} else {
				result = big.NewInt(r)
			}
		}
	}

	return &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "get balance success",
		Balance: result.String(),
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	var resp *account.TransactionResponse[account.AccountTxResponse]
	var err error
	if req.ContractAddress != "0x00" {
		resp, err = a.tronScan.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "token")
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = a.tronScan.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "normal")
		if err != nil {
			return nil, err
		}
	}
	var tx_list []*wallet2.TxMessage
	for _, tx := range resp.TransactionList {
		tx_list = append(tx_list, &wallet2.TxMessage{
			Hash:     tx.TxId,
			Froms:    []*wallet2.Address{{Address: tx.From}},
			Tos:      []*wallet2.Address{{Address: tx.To}},
			Fee:      tx.TxFee,
			Values:   []*wallet2.Value{{Value: tx.Amount}},
			Status:   wallet2.TxStatus_Success,
			Datetime: tx.TransactionTime,
			Height:   tx.Height,
		})
	}
	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transactions fail",
		Tx:   tx_list,
	}, err
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	log.Info("GetTxByHash", "req", req)
	grpcClient := a.getClient().grpcClient

	tx, err := grpcClient.GetTransactionByID(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail",
		}, err
	}

	r := tx.RawData.Contract
	if len(r) != 1 {
		err = fmt.Errorf("GetTxByHash, unsupport tx %v", req.Hash)
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail, unsupport tx",
		}, err
	}

	txi, err := grpcClient.GetTransactionInfoByID(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	var depositList []depositInfo
	switch r[0].Type {
	case core.Transaction_Contract_TransferContract:
		depositList, err = decodeTransferContract(r[0], req.Hash)
		if err != nil {
			return &wallet2.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}

	case core.Transaction_Contract_TransferAssetContract:
		depositList, err = decodeTransferAssetContract(r[0], req.Hash)
		if err != nil {
			return &wallet2.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}

	case core.Transaction_Contract_TriggerSmartContract:
		depositList, err = decodeTriggerSmartContract(r[0], txi, req.Hash)
		if err != nil {
			return &wallet2.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}
	default:
		err = fmt.Errorf("QueryTransaction, unsupport contract type %v, tx hash %v ", r[0].Type, req.Hash)
		log.Info(err.Error())
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	//Note: decodeTriggerSmartContract supports multi TRC20 transfer in single hash,  but assume we will initiate single TRC20 transfer
	// in single hash, QueryAccountTransaction is supposed to query self-initiated transaction
	if len(depositList) > 1 {
		err = fmt.Errorf("QueryTransaction, more than 1 deposit list %v, tx hash %v ", len(depositList), req.Hash)
		log.Info(err.Error())
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	from_addrs = append(from_addrs, &wallet2.Address{Address: depositList[0].fromAddr})
	to_addrs = append(to_addrs, &wallet2.Address{Address: depositList[0].toAddr})
	value_list = append(value_list, &wallet2.Value{Value: depositList[0].amount})
	if len(depositList) == 0 {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  "get tx by hash successes",
			Tx: &wallet2.TxMessage{
				Hash:   req.Hash,
				Froms:  from_addrs,
				Tos:    to_addrs,
				Fee:    big.NewInt(txi.GetFee()).String(),
				Status: wallet2.TxStatus_Success,
				Values: value_list,
				Type:   0,
				Height: strconv.FormatInt(txi.BlockNumber, 10),
			},
		}, nil
	} else {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  "get tx by hash successes",
			Tx: &wallet2.TxMessage{
				Hash:            req.Hash,
				Froms:           from_addrs,
				Tos:             to_addrs,
				Fee:             big.NewInt(txi.GetFee()).String(),
				Status:          wallet2.TxStatus_Success,
				Values:          value_list,
				Type:            0,
				Height:          strconv.FormatInt(txi.BlockNumber, 10),
				ContractAddress: depositList[0].contractAddr,
			},
		}, nil
	}
}

func (wa *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	log.Info("QueryNonce", "req", req)
	return &wallet2.NonceResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "get nonce success",
		Nonce: "0",
	}, nil
}

func (wa WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	log.Info("SendTx", "req", req)
	var tx core.Transaction
	err := pb.Unmarshal([]byte(req.RawTx), &tx)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	rawData, err := pb.Marshal(tx.GetRawData())
	hash := getHash(rawData)

	_, err = wa.getClient().grpcClient.Broadcast(&tx)
	if err != nil {
		log.Error("broadcast tx failed", "hash", hex.EncodeToString(hash), "err", err)
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	log.Info("broadcast tx success", "hash", hex.EncodeToString(hash))
	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "broadcast tx success",
		TxHash: hex.EncodeToString(hash),
	}, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	log.Info("ConvertAddress", "req", req)
	btcecPubKey, err := btcec.ParsePubKey(req.PublicKey)
	if err != nil {
		log.Error("btcec.ParsePubKey failed", "err", err)
		return &wallet2.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	addr := address.PubkeyToAddress(*btcecPubKey.ToECDSA()).String()

	log.Debug("ConvertAddress result", "pub", hex.EncodeToString(req.PublicKey), "address", addr)
	return &wallet2.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Address: addr,
	}, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	log.Info("ValidAddress", "req", req)
	ok := strings.HasPrefix(req.Address, "T")
	grpcClient := a.getClient().grpcClient
	if !ok {
		if !a.getClient().local {
			txi, err := grpcClient.GetAssetIssueByID(req.Address)
			if err != nil {
				log.Error("invalid TRC10 issuer", "err", err)
				return &wallet2.ValidAddressResponse{
					Code:  common.ReturnCode_ERROR,
					Msg:   err.Error(),
					Valid: false,
				}, err
			}

			if txi.Id != req.Address {
				err := fmt.Errorf("unmatched TRC10 issuer:%v", req.Address)
				log.Error(err.Error())
				return &wallet2.ValidAddressResponse{
					Code:  common.ReturnCode_ERROR,
					Msg:   err.Error(),
					Valid: false,
				}, err
			}
		}
		return &wallet2.ValidAddressResponse{
			Code:             common.ReturnCode_SUCCESS,
			Valid:            true,
			CanWithdrawal:    false,
			CanonicalAddress: req.Address,
		}, nil

	}
	_, err := address.Base58ToAddress(req.Address)
	if err != nil {
		log.Error("Base58ToAddress failed", "err", err)
		return &wallet2.ValidAddressResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   err.Error(),
			Valid: false,
		}, err
	}

	isTrc20 := false
	if !a.getClient().local {
		abi, err := grpcClient.GetContractABI(req.Address)
		if err != nil {
			return &wallet2.ValidAddressResponse{
				Code:  common.ReturnCode_ERROR,
				Valid: false,
				Msg:   err.Error(),
			}, err
		}

		if abi != nil {
			isTrc20 = true
		}
	}
	return &wallet2.ValidAddressResponse{
		Code:             common.ReturnCode_SUCCESS,
		Valid:            true,
		CanWithdrawal:    !isTrc20,
		CanonicalAddress: req.Address,
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	log.Info("QueryAccountTransactionFromData", "req", req)
	var tx core.TransactionRaw
	err := pb.Unmarshal(req.RawData, &tx)
	if err != nil {
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	return queryTransactionLocal(&tx, req.Symbol)
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	log.Info("QueryTransactionFromSignedData", "req", req)
	var tx core.Transaction
	err := pb.Unmarshal(req.SignedTxData, &tx)
	if err != nil {
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	return queryTransactionLocal(tx.GetRawData(), req.Symbol)
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	log.Info("CreateAccountSignedTransaction", "chain", req.Chain, "txData", hex.EncodeToString(req.TxData), "sig", hex.EncodeToString(req.Signature), "sig's len", len(req.Signature), "pubkey", hex.EncodeToString(req.PublicKey))
	rawData := req.TxData
	hash := getHash(rawData)

	//verify signature check R|S, omit V
	if ok := crypto.VerifySignature(req.PublicKey, hash, req.Signature[:64]); !ok {
		err := fmt.Errorf("fail to verify signature, chain:%v txdata:%v, signature:%v, pubkey:%v", req.Chain, hex.EncodeToString(req.TxData), hex.EncodeToString(req.Signature), hex.EncodeToString(req.PublicKey))
		log.Info(err.Error())
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	var txRaw core.TransactionRaw
	var tx core.Transaction

	err := pb.Unmarshal(req.TxData, &txRaw)
	if err != nil {
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	tx.RawData = &txRaw
	tx.Signature = append(tx.Signature, req.Signature)

	bz, err := pb.Marshal(&tx)
	if err != nil {
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return &wallet2.CreateSignedTxResponse{
		Code:         common.ReturnCode_SUCCESS,
		SignedTxData: bz,
		Hash:         hash,
	}, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	log.Info("CreateTransaction", "req", req)
	grpcClient := a.getClient().grpcClient
	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid amount",
		}, nil
	}

	gas, err := stringToInt64(req.GasLimit)
	if err != nil {
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	var txe *api.TransactionExtention
	if req.Symbol == TronSymbol {
		txe, err = grpcClient.Transfer(req.From, req.To, amount.Int64())
		if err != nil {
			return &wallet2.CreateAccountTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, nil
		}
	} else {
		//TRC10/TRC20 both should have issuer,TRC10's issuer = "1000315", TRC20's issuer = contractAddress
		isTrc10 := false
		if req.ContractAddress == "" {
			err := fmt.Errorf(" trc10 or trc20 token %v without issuer", req.Symbol)
			return &wallet2.CreateAccountTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, nil
		}

		_, err := address.Base58ToAddress(req.ContractAddress)
		if err != nil {
			isTrc10 = true
		}

		if isTrc10 {
			// for TRC10, symbol is sign in hbtc chain, contractadress is the sign in tron chain
			txe, err = grpcClient.TransferAsset(req.From, req.To, req.ContractAddress, amount.Int64())
			if err != nil {
				return &wallet2.CreateAccountTxResponse{
					Code: common.ReturnCode_ERROR,
					Msg:  err.Error(),
				}, nil
			}
		} else {
			txe, err = grpcClient.TRC20Send(req.From, req.To, req.ContractAddress, amount, gas)
			if err != nil {
				return &wallet2.CreateAccountTxResponse{
					Code: common.ReturnCode_ERROR,
					Msg:  err.Error(),
				}, nil
			}
		}
	}

	//update expiration and recalculate  hash
	txe.Transaction.RawData.Expiration = txe.Transaction.RawData.Timestamp + MaxTimeUntillExpiration
	txr, err := pb.Marshal(txe.GetTransaction().GetRawData())
	if err != nil {
		return &wallet2.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	hash := getHash(txr)
	txe.Txid = hash

	return &wallet2.CreateAccountTxResponse{
		Code:     common.ReturnCode_SUCCESS,
		TxData:   txr,
		SignHash: hash,
	}, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	log.Error("VerifySignedTransaction", "chain", req.Chain, "signTxData", hex.EncodeToString(req.SignedTxData), "sender", req.Sender)
	var tx core.Transaction
	err := pb.Unmarshal(req.SignedTxData, &tx)
	if err != nil {
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	rawData, err := pb.Marshal(tx.GetRawData())
	if err != nil {
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	hash := getHash(rawData)
	if len(tx.Signature) != 1 {
		err := fmt.Errorf("VerifySignedTransaction, len(tx.Signature) != 1")
		log.Error(err.Error())
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	pubKey, err := crypto.SigToPub(hash, tx.Signature[0])
	if err != nil {
		msg := fmt.Sprintf("SigToPub error, hash:%v, signature:%v, pubKey:%v", hex.EncodeToString(hash), hex.EncodeToString(tx.Signature[0]), hex.EncodeToString(crypto.CompressPubkey(pubKey)))
		log.Error(msg)
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
	addr := address.PubkeyToAddress(*pubKey)
	return &wallet2.VerifySignedTxResponse{
		Code:     common.ReturnCode_SUCCESS,
		Verified: addr.String() == expectedSender,
	}, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	return &wallet2.UtxoInsResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return &wallet2.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	return &wallet2.CreateUtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return &wallet2.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	return &wallet2.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	return &wallet2.AccountResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	return &wallet2.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "tron don't support this api",
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

func stringToInt64(amount string) (int64, error) {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return 0, fmt.Errorf("fail to convert string%v to int64", amount)
	}
	return intAmount.Int64(), nil
}

func getHash(bz []byte) []byte {
	h := sha256.New()
	h.Write(bz)
	hash := h.Sum(nil)
	return hash
}

type depositInfo struct {
	tokenID      string
	fromAddr     string
	toAddr       string
	amount       string
	index        int
	contractAddr string
}

func decodeTransferContract(txContract *core.Transaction_Contract, txHash string) ([]depositInfo, error) {
	var tc core.TransferContract
	if err := ptypes.UnmarshalAny(txContract.GetParameter(), &tc); err != nil {
		return nil, err
	}
	fromAddress := address.Address(tc.OwnerAddress).String()
	toAddress := address.Address(tc.ToAddress).String()
	var tronDepositInfo depositInfo
	tronDepositInfo.tokenID = TronSymbol
	tronDepositInfo.fromAddr = fromAddress
	tronDepositInfo.toAddr = toAddress
	tronDepositInfo.amount = big.NewInt(tc.Amount).String()
	tronDepositInfo.contractAddr = ""
	return []depositInfo{tronDepositInfo}, nil
}

func decodeTransferAssetContract(txContract *core.Transaction_Contract, txHash string) ([]depositInfo, error) {
	var err error
	var tc core.TransferAssetContract
	if err := ptypes.UnmarshalAny(txContract.GetParameter(), &tc); err != nil {
		log.Error("UnmarshalAny TransferAssetContract", "hash", txHash, "err", err)
		return nil, err
	}
	fromAddress := address.Address(tc.OwnerAddress).String()
	toAddress := address.Address(tc.ToAddress).String()
	assetName := string(tc.AssetName)

	//	log.Info("decodeTransferAssetContract", "hash", txHash, "symbol", assetName, "fromAddress", fromAddress, "toAddress", toAddress, "amount", tc.Amount)
	var trc10DepositInfo depositInfo
	trc10DepositInfo.fromAddr = fromAddress
	trc10DepositInfo.toAddr = toAddress
	trc10DepositInfo.amount = big.NewInt(tc.Amount).String()
	trc10DepositInfo.contractAddr = assetName
	return []depositInfo{trc10DepositInfo}, err
}

func decodeTriggerSmartContract(txContract *core.Transaction_Contract, txi *core.TransactionInfo, txHash string) ([]depositInfo, error) {
	var tsc core.TriggerSmartContract
	if err := pb.Unmarshal(txContract.GetParameter().GetValue(), &tsc); err != nil {
		log.Error("decodeTriggerSmartContractLocal", "err", err, "hash", txHash)
		return nil, err
	}

	//decode only trc20transferMethod
	trc20TransferMethodByte, _ := hex.DecodeString(trc20TransferMethodSignature)
	if ok := bytes.HasPrefix(tsc.Data, trc20TransferMethodByte); !ok {
		return nil, nil
	}

	contractAddr := address.Address(tsc.ContractAddress).String()

	var depositList []depositInfo
	// check transfer info in log
	for i, txLog := range txi.Log {
		logAddrByte := []byte{}

		// transfer log topics must be 3
		if len(txLog.Topics) != trc20TransferTopicLen {
			log.Info("decodeTriggerSmartContract", "hash's len of topics is invalid", txHash)
			continue
		}
		if hex.EncodeToString(txLog.Topics[0]) == trc20TransferTopic {
			if len(txLog.Topics[1]) != trc20TransferAddrLen || len(txLog.Topics[2]) != trc20TransferAddrLen {
				log.Debug("decodeTriggerSmartContract", "invalid transfer addr len", txHash)
				continue
			}
			//address is 20 bytes
			fromBytes := txLog.Topics[1][12:]
			toBytes := txLog.Topics[2][12:]
			logAddrByte = append([]byte{address.TronBytePrefix}, fromBytes...)
			fromAddr := address.Address(logAddrByte).String()
			logAddrByte = append([]byte{address.TronBytePrefix}, toBytes...)
			toAddr := address.Address(logAddrByte).String()
			amount := new(big.Int).SetBytes(txLog.Data)

			//	log.Info("decodeTriggerSmartContract", "hash", txHash, "from", fromAddr, "to", toAddr, "amount", amount)

			var trc20DepositInfo depositInfo
			trc20DepositInfo.amount = amount.String()
			trc20DepositInfo.fromAddr = fromAddr
			trc20DepositInfo.toAddr = toAddr
			trc20DepositInfo.index = i
			trc20DepositInfo.contractAddr = contractAddr
			depositList = append(depositList, trc20DepositInfo)

		} else {
			//	log.Debug("decodeTriggerSmartContract", "hash is not transfer method", txHash)
			continue
		}
	}

	return depositList, nil
}

// IMPORTANT, current support only 1 TRC20 transfer
func decodeTriggerSmartContractLocal(txContract *core.Transaction_Contract, txHash string) ([]depositInfo, error) {
	var tsc core.TriggerSmartContract
	if err := pb.Unmarshal(txContract.GetParameter().GetValue(), &tsc); err != nil {
		log.Error("decodeTriggerSmartContractLocal", "err", err, "hash", txHash)
		return nil, err
	}

	//decode only trc20transferMethod
	trc20TransferMethodByte, _ := hex.DecodeString(trc20TransferMethodSignature)
	if ok := bytes.HasPrefix(tsc.Data, trc20TransferMethodByte); !ok {
		return nil, nil
	}

	fromAddr := address.Address(tsc.OwnerAddress).String()
	contractAddr := address.Address(tsc.ContractAddress).String()

	start := len(trc20TransferMethodByte)
	end := start + trc20TransferAddrLen
	start = end - address.AddressLength + 1

	addressTron := make([]byte, 0)
	addressTron = append(addressTron, address.TronBytePrefix)
	addressTron = append(addressTron, tsc.Data[start:end]...)

	toAddr := address.Address(addressTron).String()
	amount := new(big.Int).SetBytes(tsc.Data[end:])

	var trc20DepositInfo depositInfo
	trc20DepositInfo.amount = amount.String()
	trc20DepositInfo.fromAddr = fromAddr
	trc20DepositInfo.contractAddr = contractAddr
	trc20DepositInfo.toAddr = toAddr
	trc20DepositInfo.index = 0
	return []depositInfo{trc20DepositInfo}, nil
}

// queryTransaction should not be called to decode locally build
func queryTransactionLocal(txRaw *core.TransactionRaw, symbol string) (*wallet2.AccountTxResponse, error) {
	bz, err := pb.Marshal(txRaw)
	hash := hex.EncodeToString(getHash(bz))

	r := txRaw.Contract
	if len(r) != 1 {
		err := fmt.Errorf("QueryTransactionFromSignedData, tx's len(contract): %v !=1", len(r))
		log.Error(err.Error())
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	var depositList []depositInfo
	switch r[0].Type {
	case core.Transaction_Contract_TransferContract:
		depositList, err = decodeTransferContract(r[0], hash)
		if err != nil {
			return &wallet2.AccountTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}

	case core.Transaction_Contract_TransferAssetContract:
		depositList, err = decodeTransferAssetContract(r[0], hash)
		if err != nil {
			return &wallet2.AccountTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}

	case core.Transaction_Contract_TriggerSmartContract:
		depositList, err = decodeTriggerSmartContractLocal(r[0], hash)
		if err != nil {
			return &wallet2.AccountTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err.Error(),
			}, err
		}
	default:
		err = fmt.Errorf("QueryTransaction, unsupport contract type %v, tx hash %v ", r[0].Type, hash)
		log.Info(err.Error())
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if len(depositList) > 1 {
		err = fmt.Errorf("QueryTransaction, more than 1 deposit list %v, tx hash %v ", len(depositList), hash)
		log.Info(err.Error())
		return &wallet2.AccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	gasLimit := txRaw.FeeLimit
	if gasLimit == 0 {
		gasLimit = defaultGasLimit
	}

	return &wallet2.AccountTxResponse{
		Code:            common.ReturnCode_SUCCESS,
		TxHash:          hash,
		From:            depositList[0].fromAddr,
		To:              depositList[0].toAddr,
		Amount:          depositList[0].amount,
		Memo:            "",
		Nonce:           0,
		GasPrice:        "1",
		GasLimit:        big.NewInt(gasLimit).String(),
		SignHash:        getHash(bz),
		ContractAddress: depositList[0].contractAddr,
	}, nil

}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}
