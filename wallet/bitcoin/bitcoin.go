package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"github.com/savour-labs/wallet-chain-node/wallet/multiclient"
)

const (
	confirms         = 1
	btcDecimals      = 8
	btcFeeBlocks     = 3
	ChainName        = "Bitcoin"
	Symbol           = "BTC"
	BlockChainApiUrl = "https://blockchain.info"
	OkLinkUrl        = "https://www.oklink.com"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients   *multiclient.MultiClient
	bcClient  *BcClient
	oklClient *OkLinkClient
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

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	clients, err := newBtcClients(conf)
	if err != nil {
		return nil, err
	}
	return newChainAdaptorWithClients(clients), nil
}

func NewLocalChainAdaptor(network config.NetWorkType) wallet.WalletAdaptor {
	return newChainAdaptorWithClients([]*btcClient{newLocalBtcClient(network)})
}

func newChainAdaptorWithClients(clients []*btcClient) *WalletAdaptor {
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}
	bcClient, err := NewBlockChainClient(BlockChainApiUrl)
	if err != nil {
		log.Error("new blockchain client fail", "err", err)
	}
	oklClient, err := NewOkLinkClient(OkLinkUrl)
	if err != nil {
		log.Error("new oklink client fail", "err", err)
	}
	return &WalletAdaptor{
		clients:   multiclient.New(clis),
		bcClient:  bcClient,
		oklClient: oklClient,
	}
}

func (a *WalletAdaptor) getClient() *btcClient {
	return a.clients.BestClient().(*btcClient)
}

func (w *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return &wallet2.SupportCoinsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
		Support: true,
	}, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	addressPubKey, err := btcutil.NewAddressPubKey(req.PublicKey, a.getClient().GetNetwork())
	if err != nil {
		return &wallet2.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Address: addressPubKey.EncodeAddress(),
	}, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	address, err := btcutil.DecodeAddress(req.Address, a.getClient().GetNetwork())
	if err != nil {
		return &wallet2.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if !address.IsForNet(a.getClient().GetNetwork()) {
		err := errors.New("address is not valid for this network")
		return &wallet2.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.ValidAddressResponse{
		Code:             common.ReturnCode_SUCCESS,
		Valid:            true,
		CanWithdrawal:    true,
		CanonicalAddress: address.String(),
	}, nil
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	balance, err := a.bcClient.GetAccountBalance(req.Address)
	if err != nil {
		return &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get btc balance fail",
			Balance: "0",
		}, err
	}
	return &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "get btc balance success",
		Balance: balance,
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	transaction, err := a.bcClient.GetTransactionsByAddress(req.Address, strconv.Itoa(int(req.Page)), strconv.Itoa(int(req.Pagesize)))
	if err != nil {
		return &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction list fail",
			Tx:   nil,
		}, err
	}
	var tx_list []*wallet2.TxMessage
	for _, ttxs := range transaction.Txs {
		var from_addrs []*wallet2.Address
		var to_addrs []*wallet2.Address
		var value_list []*wallet2.Value
		var direction int32
		for _, inputs := range ttxs.Inputs {
			from_addrs = append(from_addrs, &wallet2.Address{Address: inputs.PrevOut.Addr})
		}
		tx_fee := ttxs.Fee
		for _, out := range ttxs.Out {
			to_addrs = append(to_addrs, &wallet2.Address{Address: out.Addr})
			value_list = append(value_list, &wallet2.Value{Value: out.Value.String()})
		}
		datetime := ttxs.Time.String()
		if strings.EqualFold(req.Address, from_addrs[0].Address) {
			direction = 0
		} else {
			direction = 1
		}
		tx := &wallet2.TxMessage{
			Hash:            ttxs.Hash,
			Froms:           from_addrs,
			Tos:             to_addrs,
			Values:          value_list,
			Fee:             tx_fee.String(),
			Status:          wallet2.TxStatus_Success,
			Type:            direction,
			Height:          ttxs.BlockHeight.String(),
			ContractAddress: "0x00",
			Datetime:        datetime,
		}
		tx_list = append(tx_list, tx)
	}

	return &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction list success",
		Tx:   tx_list,
	}, err
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	transaction, err := a.bcClient.GetTransactionsByHash(req.Hash)
	if err != nil {
		return &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction list fail",
			Tx:   nil,
		}, err
	}
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	for _, inputs := range transaction.Inputs {
		from_addrs = append(from_addrs, &wallet2.Address{Address: inputs.PrevOut.Addr})
	}
	tx_fee := transaction.Fee
	for _, out := range transaction.Out {
		to_addrs = append(to_addrs, &wallet2.Address{Address: out.Addr})
		value_list = append(value_list, &wallet2.Value{Value: out.Value.String()})
	}
	datetime := transaction.Time.String()
	txMsg := &wallet2.TxMessage{
		Hash:            transaction.Hash,
		Froms:           from_addrs,
		Tos:             to_addrs,
		Values:          value_list,
		Fee:             tx_fee.String(),
		Status:          wallet2.TxStatus_Success,
		Type:            0,
		Height:          transaction.BlockHeight.String(),
		ContractAddress: "0x00",
		Datetime:        datetime,
	}
	return &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   txMsg,
	}, nil
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	utxoList, err := a.bcClient.GetAccountUtxo(req.Address)
	if err != nil {
		return &wallet2.UnspentOutputsResponse{
			Code:           common.ReturnCode_ERROR,
			Msg:            err.Error(),
			UnspentOutputs: nil,
		}, err
	}
	var unspentOutputList []*wallet2.UnspentOutput
	for _, value := range utxoList {
		unspentOutput := &wallet2.UnspentOutput{
			TxHashBigEndian: value.TxHashBigEndian,
			TxHash:          value.TxHash,
			TxOutputN:       value.TxOutputN,
			Script:          value.Script,
			Value:           value.Value,
			Confirmations:   value.Confirmations,
			TxIndex:         value.TxIndex,
		}
		unspentOutputList = append(unspentOutputList, unspentOutput)
	}
	return &wallet2.UnspentOutputsResponse{
		Code:           common.ReturnCode_SUCCESS,
		Msg:            "get unspent outputs success",
		UnspentOutputs: unspentOutputList,
	}, err
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	utxo := req.Vin
	txhash, err := chainhash.NewHashFromStr(utxo.Hash)
	if err != nil {
		log.Info("QueryUtxo NewHashFromStr", "err", err)
		return &wallet2.UtxoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	reply, err := a.getClient().GetTxOut(txhash, utxo.Index, true)
	if err != nil {
		log.Info("QueryUtxo GetTxOut", "err", err)
		return &wallet2.UtxoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	if reply == nil {
		log.Info("QueryUtxo GetTxOut", "err", "hash not found")
		err = errors.New("hash not found")
		return &wallet2.UtxoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	if btcToSatoshi(reply.Value).Int64() != utxo.Amount {
		log.Info("QueryUtxo GetTxOut", "err", "amount mismatch")
		err = errors.New("amount mismatch")
		return &wallet2.UtxoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	tx, err := a.getClient().GetRawTransactionVerbose(txhash)
	if err != nil {
		log.Info("QueryUtxo GetRawTransactionVerbose", "err", err)
		return &wallet2.UtxoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	if tx.Vout[utxo.Index].ScriptPubKey.Addresses[0] != utxo.Address {
		log.Info("QueryUtxo GetTxOut", "err", "address mismatch")
		err := errors.New("address mismatch")
		return &wallet2.UtxoResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.UtxoResponse{
		Code:    common.ReturnCode_SUCCESS,
		Unspent: true,
	}, nil
}

func (w *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	gasFee, err := w.oklClient.GetGasFee("btc")
	if err != nil {
		log.Info("QueryGasPrice", "err", err)
		return &wallet2.GasPriceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get gas price success",
		Gas:  gasFee,
	}, nil
}

func (w *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	r := bytes.NewReader([]byte(req.RawTx))
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	txHash, err := w.getClient().SendRawTransaction(&msgTx)
	if err != nil {
		return &wallet2.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if strings.Compare(msgTx.TxHash().String(), txHash.String()) != 0 {
		log.Error("BroadcastTransaction, txhash mismatch", "local hash", msgTx.TxHash().String(), "hash from net", txHash.String(), "signedTx", req.RawTx)
	}

	return &wallet2.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: txHash.String(),
	}, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	log.Info("QueryUtxoInsFromData", "req", req)
	vins, err := decodecommonVinsFromData(req.Data)
	if err != nil {
		return &wallet2.UtxoInsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.UtxoInsResponse{
		Code: common.ReturnCode_SUCCESS,
		Vins: vins,
	}, err
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	res, err := a.decodeTx(req.RawData, req.Vins, false)
	if err != nil {
		log.Info("QueryTransactionFromData decodeTx", "err", err)

		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.UtxoTxResponse{
		Code:       common.ReturnCode_SUCCESS,
		SignHashes: res.SignHashes,
		Status:     wallet2.TxStatus_Other,
		Vins:       res.Vins,
		Vouts:      res.Vouts,
		CostFee:    res.CostFee.String(),
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	res, err := a.decodeTx(req.SignedTxData, req.Vins, true)
	if err != nil {
		log.Info("QueryTransactionFromSignedData decodeTx", "err", err)

		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.UtxoTxResponse{
		Code:       common.ReturnCode_SUCCESS,
		TxHash:     res.Hash,
		Status:     wallet2.TxStatus_Other,
		Vins:       res.Vins,
		Vouts:      res.Vouts,
		CostFee:    res.CostFee.String(),
		SignHashes: res.SignHashes,
	}, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	r := bytes.NewReader(req.TxData)
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		log.Error("CreateSignedTransaction msgTx.Deserialize", "err", err)

		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if len(req.Signatures) != len(msgTx.TxIn) {
		log.Error("CreateSignedTransaction invalid params", "err", "Signature number mismatch Txin number")
		err = errors.New("Signature number != Txin number")
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if len(req.PublicKeys) != len(msgTx.TxIn) {
		log.Error("CreateSignedTransaction invalid params", "err", "Pubkey number mismatch Txin number")
		err = errors.New("Pubkey number != Txin number")
		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// assemble signatures
	for i, in := range msgTx.TxIn {
		btcecPub, err2 := btcec.ParsePubKey(req.PublicKeys[i])
		if err2 != nil {
			log.Error("CreateSignedTransaction ParsePubKey", "err", err2)
			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		var pkData []byte
		if btcec.IsCompressedPubKey(req.PublicKeys[i]) {
			pkData = btcecPub.SerializeCompressed()
		} else {
			pkData = btcecPub.SerializeUncompressed()
		}

		// verify transaction
		preTx, err2 := a.getClient().GetRawTransactionVerbose(&in.PreviousOutPoint.Hash)
		if err2 != nil {
			log.Error("CreateSignedTransaction GetRawTransactionVerbose", "err", err2)

			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		log.Info("CreateSignedTransaction ", "from address", preTx.Vout[in.PreviousOutPoint.Index].ScriptPubKey.Addresses[0])

		fromAddress, err2 := btcutil.DecodeAddress(preTx.Vout[in.PreviousOutPoint.Index].ScriptPubKey.Addresses[0], a.getClient().GetNetwork())
		if err2 != nil {
			log.Error("CreateSignedTransaction DecodeAddress", "err", err2)

			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		fromPkScript, err2 := txscript.PayToAddrScript(fromAddress)
		if err2 != nil {
			log.Error("CreateSignedTransaction PayToAddrScript", "err", err2)

			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		// creat sigscript and verify
		if len(req.Signatures[i]) < 64 {
			err2 = errors.New("Invalid signature length")
			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}
		var r *btcec.ModNScalar
		R := r.SetInt(r.SetBytes((*[32]byte)(req.Signatures[i][0:32])))
		var s *btcec.ModNScalar
		S := s.SetInt(r.SetBytes((*[32]byte)(req.Signatures[i][32:64])))
		btcecSig := ecdsa.NewSignature(R, S)
		sig := append(btcecSig.Serialize(), byte(txscript.SigHashAll))
		sigScript, err2 := txscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
		if err2 != nil {
			log.Error("CreateSignedTransaction NewScriptBuilder", "err", err2)

			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}

		msgTx.TxIn[i].SignatureScript = sigScript
		amount := btcToSatoshi(preTx.Vout[in.PreviousOutPoint.Index].Value).Int64()
		log.Info("CreateSignedTransaction ", "amount", preTx.Vout[in.PreviousOutPoint.Index].Value, "int amount", amount)

		vm, err2 := txscript.NewEngine(fromPkScript, &msgTx, i, txscript.StandardVerifyFlags, nil, nil, amount)
		if err2 != nil {
			log.Error("CreateSignedTransaction NewEngine", "err", err2)

			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err2.Error(),
			}, err2
		}
		if err3 := vm.Execute(); err3 != nil {
			log.Error("CreateSignedTransaction NewEngine Execute", "err", err3)

			return &wallet2.CreateSignedTxResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  err3.Error(),
			}, err3
		}

	}

	// serialize tx
	buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))

	err = msgTx.Serialize(buf)
	if err != nil {
		log.Error("CreateSignedTransaction tx Serialize", "err", err)

		return &wallet2.CreateSignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	hash := msgTx.TxHash()
	return &wallet2.CreateSignedTxResponse{
		Code:         common.ReturnCode_SUCCESS,
		SignedTxData: buf.Bytes(),
		Hash:         (&hash).CloneBytes(),
	}, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	vinNum := len(req.Vins)
	var totalAmountIn, totalAmountOut int64

	if vinNum == 0 {
		err := fmt.Errorf("no Vin in req:%v", req)
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	fee, ok := big.NewInt(0).SetString(req.Fee, 0)
	if !ok {
		err := errors.New("CreateTransaction, fail to get fee")
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	for _, in := range req.Vins {
		totalAmountIn += in.Amount
	}
	for _, out := range req.Vouts {
		totalAmountOut += out.Amount
	}
	if totalAmountIn != totalAmountOut+fee.Int64() {
		err := errors.New("CreateTransaction, total amount in != total amount out + fee")
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	rawTx, err := a.createRawTx(req.Vins, req.Vouts)
	if err != nil {
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	err = rawTx.Serialize(buf)
	if err != nil {
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	signHashes, err := a.calcSignHashes(req.Vins, req.Vouts)
	if err != nil {
		return &wallet2.CreateUtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	log.Info("CreateTransaction", "usigned tx", hex.EncodeToString(buf.Bytes()))

	return &wallet2.CreateUtxoTxResponse{
		Code:       common.ReturnCode_SUCCESS,
		TxData:     buf.Bytes(),
		SignHashes: signHashes,
	}, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	_, err := a.decodeTx(req.SignedTxData, req.Vins, true)
	if err != nil {
		log.Error("VerifySignedTransaction", "decodeTx err", err)
		return &wallet2.VerifySignedTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &wallet2.VerifySignedTxResponse{
		Code:     common.ReturnCode_SUCCESS,
		Verified: true,
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	return &wallet2.AccountResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	return &wallet2.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (w *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	return &wallet2.NonceResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return &wallet2.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	return &wallet2.AccountTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	return &wallet2.AccountTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return &wallet2.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	return &wallet2.CreateAccountTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "bitcoin don't support this api",
	}, nil
}

func (a *WalletAdaptor) queryTransaction(txhash *chainhash.Hash) (*wallet2.UtxoTxResponse, error) {
	tx, err := a.getClient().GetRawTransactionVerbose(txhash)
	if err != nil {
		if rpcErr, ok := err.(*btcjson.RPCError); ok && rpcErr.Code == btcjson.ErrRPCBlockNotFound {
			return &wallet2.UtxoTxResponse{
				Code:   common.ReturnCode_SUCCESS,
				Status: wallet2.TxStatus_NotFound,
			}, nil
		}
		log.Error("queryTransaction GetRawTransactionVerbose", "err", err)
		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	if tx == nil || txhash.String() != tx.Txid {
		log.Error("queryTransaction txid mismatch")

		return &wallet2.UtxoTxResponse{
			Code:   common.ReturnCode_SUCCESS,
			Status: wallet2.TxStatus_NotFound,
		}, nil
	}

	if tx.Confirmations < confirms {
		log.Error("queryTransaction confirmes too low", "tx confirms", tx.Confirmations, "need confirms", confirms)

		return &wallet2.UtxoTxResponse{
			Code:   common.ReturnCode_SUCCESS,
			Status: wallet2.TxStatus_Pending,
		}, nil
	}

	blockHash, _ := chainhash.NewHashFromStr(tx.BlockHash)
	block, err := a.getClient().GetBlockVerbose(blockHash)
	if err != nil {
		log.Error("queryTransaction GetBlockVerbose", "err", err)

		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	reply, err := a.assembleUtxoTransactionReply(tx, block.Height, block.Time, func(txid string, index uint32) (int64, string, error) {
		preHash, err2 := chainhash.NewHashFromStr(txid)
		if err2 != nil {
			return 0, "", err2
		}
		preTx, err2 := a.getClient().GetRawTransactionVerbose(preHash)
		if err2 != nil {
			return 0, "", err2
		}
		amount := btcToSatoshi(preTx.Vout[index].Value).Int64()

		return amount, preTx.Vout[index].ScriptPubKey.Addresses[0], nil
	})
	if err != nil {
		return &wallet2.UtxoTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return reply, nil
}

func (a *WalletAdaptor) assembleUtxoTransactionReply(tx *btcjson.TxRawResult, blockHeight, blockTime int64, getPrevTxInfo func(txid string, index uint32) (int64, string, error)) (*wallet2.UtxoTxResponse, error) {
	var totalAmountIn, totalAmountOut int64
	ins := make([]*wallet2.Vin, 0, len(tx.Vin))
	outs := make([]*wallet2.Vout, 0, len(tx.Vout))
	for _, in := range tx.Vin {
		amount, address, err := getPrevTxInfo(in.Txid, in.Vout)
		if err != nil {
			return nil, err
		}
		totalAmountIn += amount

		t := wallet2.Vin{
			Hash:    in.Txid,
			Index:   in.Vout,
			Amount:  amount,
			Address: address,
		}
		ins = append(ins, &t)
	}

	for index, out := range tx.Vout {
		amount := btcToSatoshi(out.Value).Int64()
		addr := ""
		if len(out.ScriptPubKey.Addresses) > 0 {
			addr = out.ScriptPubKey.Addresses[0]
		}

		totalAmountOut += amount
		t := wallet2.Vout{
			Address: addr,
			Amount:  amount,
			Index:   uint32(index),
		}
		outs = append(outs, &t)
	}
	gasUsed := totalAmountIn - totalAmountOut
	reply := &wallet2.UtxoTxResponse{
		Code:        common.ReturnCode_SUCCESS,
		TxHash:      tx.Txid,
		Status:      wallet2.TxStatus_Success,
		Vins:        ins,
		Vouts:       outs,
		CostFee:     strconv.FormatInt(gasUsed, 10),
		BlockHeight: uint64(blockHeight),
		BlockTime:   uint64(blockTime),
	}
	return reply, nil
}

func (a *WalletAdaptor) assembleUtxoTransactionReplyForTxHash(tx *btcjson.TxRawResult, blockHeight, blockTime int64, getPrevTxInfo func(txid string, index uint32) (int64, string, error)) (*wallet2.TxHashResponse, error) {
	var totalAmountIn, totalAmountOut int64
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	var value_list []*wallet2.Value
	var direction int32
	for _, in := range tx.Vin {
		amount, address, err := getPrevTxInfo(in.Txid, in.Vout)
		if err != nil {
			return nil, err
		}
		totalAmountIn += amount
		from_addrs = append(from_addrs, &wallet2.Address{Address: address})
		value_list = append(value_list, &wallet2.Value{Value: strconv.FormatInt(totalAmountIn, 10)})
	}
	for _, out := range tx.Vout {
		amount := btcToSatoshi(out.Value).Int64()
		totalAmountOut += amount
		addr := ""
		if len(out.ScriptPubKey.Addresses) > 0 {
			addr = out.ScriptPubKey.Addresses[0]
		}
		to_addrs = append(to_addrs, &wallet2.Address{Address: addr})
		value_list = append(value_list, &wallet2.Value{Value: strconv.FormatInt(totalAmountOut, 10)})
	}
	gasUsed := totalAmountIn - totalAmountOut
	reply := &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get tx success",
		Tx: &wallet2.TxMessage{
			Hash:   tx.Hash,
			Status: wallet2.TxStatus_Success,
			Froms:  from_addrs,
			Tos:    to_addrs,
			Fee:    strconv.FormatInt(gasUsed, 10),
			Values: value_list,
			Height: strconv.FormatInt(blockHeight, 10),
			Type:   direction,
		},
	}
	return reply, nil
}

func btcToSatoshi(btcCount float64) *big.Int {
	amount := strconv.FormatFloat(btcCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(btcDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}

func (a *WalletAdaptor) createRawTx(ins []*wallet2.Vin, outs []*wallet2.Vout) (*wire.MsgTx, error) {
	if len(ins) == 0 || len(outs) == 0 {
		return nil, errors.New("invalid len in or out")
	}

	rawTx := wire.NewMsgTx(wire.TxVersion)
	for _, in := range ins {
		// convert string hash to a bitcoin hash
		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
		if err != nil {
			return nil, err
		}
		// make a txin
		txIn := wire.NewTxIn(wire.NewOutPoint(utxoHash, in.Index), nil, nil)
		// add txIn to transaction
		rawTx.AddTxIn(txIn)
	}

	// build tx output
	for _, out := range outs {
		if strings.HasPrefix(out.Address, omniPrefix) {
			toPkScript, err := buildOmniScript(out.Address)
			if err != nil {
				return nil, err
			}
			rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
			continue
		}

		toAddress, err := btcutil.DecodeAddress(out.Address, a.getClient().GetNetwork())
		if err != nil {
			return nil, err
		}
		// build the pkScript
		toPkScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return nil, err
		}
		// add txOut to transaction
		rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
	}
	return rawTx, nil
}

func buildOmniScript(addr string) ([]byte, error) {
	omniData, err := hex.DecodeString(addr)
	if err != nil {
		return nil, errors.New("parse omni data error")
	}
	return txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(omniData).Script()
}

type DecodeTxRes struct {
	Hash       string
	SignHashes [][]byte
	Vins       []*wallet2.Vin
	Vouts      []*wallet2.Vout
	CostFee    *big.Int
}

func (a *WalletAdaptor) decodeTx(txData []byte, vins []*wallet2.Vin, sign bool) (*DecodeTxRes, error) {
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(bytes.NewReader(txData))
	if err != nil {
		return nil, err
	}

	offline := true
	if len(vins) == 0 {
		offline = false
	}
	if offline && len(vins) != len(msgTx.TxIn) {
		return nil, status.Error(codes.InvalidArgument, "the length of deserialized tx's in differs from vin in req")
	}

	ins, totalAmountIn, err := a.decodeVins(msgTx, offline, vins, sign)
	if err != nil {
		return nil, err
	}

	outs, totalAmountOut, err := a.decodeVouts(msgTx)
	if err != nil {
		return nil, err
	}

	// build the pkScript and Generate signhash for each Vin,
	signHashes, err := a.calcSignHashes(ins, outs)
	if err != nil {
		return nil, err
	}

	res := DecodeTxRes{
		SignHashes: signHashes,
		Vins:       ins,
		Vouts:      outs,
		CostFee:    totalAmountIn.Sub(totalAmountIn, totalAmountOut),
	}
	if sign {
		res.Hash = msgTx.TxHash().String()
	}
	return &res, nil
}

func (a *WalletAdaptor) decodeVouts(msgTx wire.MsgTx) ([]*wallet2.Vout, *big.Int, error) {
	outs := make([]*wallet2.Vout, 0, len(msgTx.TxOut))
	totalAmountOut := big.NewInt(0)
	for _, out := range msgTx.TxOut {
		var t wallet2.Vout
		_, pubkeyAddrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, a.getClient().GetNetwork())
		if err != nil {
			return nil, nil, err
		}
		t.Address = pubkeyAddrs[0].EncodeAddress()
		t.Amount = out.Value
		totalAmountOut.Add(totalAmountOut, big.NewInt(t.Amount))
		outs = append(outs, &t)
	}
	return outs, totalAmountOut, nil
}

func (a *WalletAdaptor) decodeVins(msgTx wire.MsgTx, offline bool, vins []*wallet2.Vin, sign bool) ([]*wallet2.Vin, *big.Int, error) {
	// verify signatures and decode
	ins := make([]*wallet2.Vin, 0, len(msgTx.TxIn))
	totalAmountIn := big.NewInt(0)
	for index, in := range msgTx.TxIn {
		vin, err := a.getVin(offline, vins, index, in)
		if err != nil {
			return nil, nil, err
		}

		if sign {
			err = a.verifySign(vin, msgTx, index)
			if err != nil {
				return nil, nil, err
			}
		}

		totalAmountIn.Add(totalAmountIn, big.NewInt(vin.Amount))
		ins = append(ins, vin)
	}
	return ins, totalAmountIn, nil
}

func (a *WalletAdaptor) getVin(offline bool, vins []*wallet2.Vin, index int, in *wire.TxIn) (*wallet2.Vin, error) {
	var vin *wallet2.Vin
	if offline {
		vin = vins[index]
	} else {
		preTx, err := a.getClient().GetRawTransactionVerbose(&in.PreviousOutPoint.Hash)
		if err != nil {
			return nil, err
		}
		out := preTx.Vout[in.PreviousOutPoint.Index]
		vin = &wallet2.Vin{
			Hash:    "",
			Index:   0,
			Amount:  btcToSatoshi(out.Value).Int64(),
			Address: out.ScriptPubKey.Addresses[0],
		}
	}
	vin.Hash = in.PreviousOutPoint.Hash.String()
	vin.Index = in.PreviousOutPoint.Index
	return vin, nil
}

func (a *WalletAdaptor) verifySign(vin *wallet2.Vin, msgTx wire.MsgTx, index int) error {
	fromAddress, err := btcutil.DecodeAddress(vin.Address, a.getClient().GetNetwork())
	if err != nil {
		return err
	}

	fromPkScript, err := txscript.PayToAddrScript(fromAddress)
	if err != nil {
		return err
	}

	vm, err := txscript.NewEngine(fromPkScript, &msgTx, index, txscript.StandardVerifyFlags, nil, nil, vin.Amount)
	if err != nil {
		return err
	}
	return vm.Execute()
}

func (a *WalletAdaptor) calcSignHashes(Vins []*wallet2.Vin, Vouts []*wallet2.Vout) ([][]byte, error) {

	rawTx, err := a.createRawTx(Vins, Vouts)
	if err != nil {
		return nil, err
	}

	signHashes := make([][]byte, len(Vins))
	for i, in := range Vins {
		from := in.Address
		fromAddr, err := btcutil.DecodeAddress(from, a.getClient().GetNetwork())
		if err != nil {
			log.Info("DecodeAddress err", "from", from, "err", err)
			return nil, err
		}
		fromPkScript, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			log.Info("PayToAddrScript err", "err", err)
			return nil, err
		}

		signHash, err := txscript.CalcSignatureHash(fromPkScript, txscript.SigHashAll, rawTx, i)
		if err != nil {
			log.Info("CalcSignatureHash err", "err", err)
			return nil, err
		}
		signHashes[i] = signHash
	}
	return signHashes, nil
}

func decodecommonVinsFromData(data []byte) ([]*wallet2.Vin, error) {
	r := bytes.NewReader(data)
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return nil, err
	}
	ins := make([]*wallet2.Vin, len(msgTx.TxIn))
	for i, in := range msgTx.TxIn {
		t := wallet2.Vin{}
		t.Hash = in.PreviousOutPoint.Hash.String()
		t.Index = in.PreviousOutPoint.Index
		ins[i] = &t
	}
	return ins, nil
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}
