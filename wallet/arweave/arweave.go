package arweave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	// ChainName is used to register the chain to walletdispatcher.WalletDispatcher
	ChainName = "Arweave"

	// Coin is the coin name
	Coin = "AR"
)

type GetTxByHashResp struct {
	Format    int          `json:"format"`
	ID        string       `json:"id"`
	LastTx    string       `json:"last_tx"`
	Owner     string       `json:"owner"`
	Tags      []*TxTagResp `json:"tags"`
	Target    string       `json:"target"`
	Quantity  string       `json:"quantity"`
	DataRoot  string       `json:"data_root"`
	DataSize  string       `json:"data_size"`
	Reward    string       `json:"reward"`
	Signature string       `json:"signature"`
}

type TxTagResp struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetTxStatusResp struct {
	BlockHeight           int    `json:"block_height"`
	BlockIndepHash        string `json:"block_indep_hash"`
	NumberOfConfirmations int    `json:"number_of_confirmations"`
}

type GraphQLQuery struct {
	Query string `json:"query"`
}

// GetTxByAddressResp is the top-level structure representing the entire JSON from GraphQL response.
type GetTxByAddressResp struct {
	Data Data `json:"data"`
}

// Data represents the data field in the JSON.
type Data struct {
	Transactions Transactions `json:"transactions"`
}

// Transactions represents the transactions field in the JSON.
type Transactions struct {
	Edges []Edge `json:"edges"`
}

// Edge represents each edge in the transactions field.
type Edge struct {
	Node Node `json:"node"`
}

// Node represents each node in the edge.
type Node struct {
	ID        string    `json:"id"`
	Anchor    string    `json:"anchor"`
	Signature string    `json:"signature"`
	Recipient string    `json:"recipient"`
	Owner     Owner     `json:"owner"`
	Fee       Fee       `json:"fee"`
	Quantity  Quantity  `json:"quantity"`
	Data      DataField `json:"data"`
	Tags      []Tag     `json:"tags"`
	Block     *Block    `json:"block"`  // // Use pointer to distinguish null from empty struct
	Parent    *string   `json:"parent"` // Use pointer to distinguish null from empty string
}

// Owner represents the owner field in the node.
type Owner struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

// Fee represents the fee field in the node.
type Fee struct {
	Winston string `json:"winston"`
	Ar      string `json:"ar"`
}

// Quantity represents the quantity field in the node.
type Quantity struct {
	Winston string `json:"winston"`
	Ar      string `json:"ar"`
}

// DataField represents the data field in the node.
type DataField struct {
	Size string `json:"size"`
	Type string `json:"type"`
}

// Tag represents each tag in the tags array.
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Block represents the block field in the node.
type Block struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Height    int64  `json:"height"`
	Previous  string `json:"previous"`
}

type WalletAdaptor struct {
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

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	return &WalletAdaptor{}, nil
}

func (a *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetBlock(req *wallet2.BlockRequest) (*wallet2.BlockResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	walletAddress := req.Address
	url := "https://arweave.net/wallet/" + walletAddress + "/balance"

	arWeaveResp, err := http.Get(url)
	if err != nil {
		log.Error("Get arweave balance failed", "err", err)
		resp := &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}
		return resp, err
	}

	balance, err := ioutil.ReadAll(arWeaveResp.Body)
	if err != nil {
		log.Error("Read balance failed", "err", err)
		resp := &wallet2.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}
		return resp, err
	}

	resp := &wallet2.BalanceResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "get balance success",
		Balance: string(balance),
	}
	return resp, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	address := req.Address
	url := "https://arweave.net/graphql"

	// build a GraphQL query
	query := fmt.Sprintf(`
    {
    transactions(
        owners: ["%s"]
    ) {
        edges {
            node {
                id
                anchor
                signature
                recipient
                owner {
                    address
                    key
                }
                fee {
                    winston
                    ar
                }
                quantity {
                    winston
                    ar
                }
                data {
                    size
                    type
                }
                tags {
                    name
                    value
                }
                block {
                    id
                    timestamp
                    height
                    previous
                }
                parent {
                    id
                }
            }
        }
    }
}`, address)

	requestBody, err := json.Marshal(GraphQLQuery{Query: query})
	if err != nil {
		log.Error("Marshal graphql query failed", "err", err)
		resp := &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	// send POST request
	graphResp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Get tx by address failed", "err", err)
		resp := &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	// read response body
	graphBody, err := ioutil.ReadAll(graphResp.Body)
	if err != nil {
		log.Error("Read tx by address failed", "err", err)
		resp := &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	txAddressResp := &GetTxByAddressResp{}
	err = json.Unmarshal(graphBody, txAddressResp)
	if err != nil {
		log.Error("Unmarshal tx by address failed", "err", err)
		resp := &wallet2.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	txList := make([]*wallet2.TxMessage, 0, len(txAddressResp.Data.Transactions.Edges))
	for _, edge := range txAddressResp.Data.Transactions.Edges {
		tx := &wallet2.TxMessage{
			Hash:            edge.Node.ID,
			Index:           0,
			Froms:           []*wallet2.Address{{Address: edge.Node.Owner.Address}},
			Tos:             []*wallet2.Address{{Address: edge.Node.Recipient}},
			Values:          []*wallet2.Value{{Value: edge.Node.Quantity.Ar}},
			Fee:             edge.Node.Fee.Ar,
			Status:          0,
			Type:            0,
			Height:          strconv.FormatInt(edge.Node.Block.Height, 10),
			ContractAddress: "",
			Datetime:        time.Unix(edge.Node.Block.Timestamp, 0).Format("2006-01-02 15:04:05"),
		}
		txList = append(txList, tx)
	}

	resp := &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get address tx success",
		Tx:   txList,
	}

	return resp, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	hash := req.Hash
	txHashUrl := "https://arweave.net/tx/" + hash
	txHashResp, err := http.Get(txHashUrl)
	if err != nil {
		log.Error("Get tx by hash failed", "err", err)
		resp := &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	tx, err := ioutil.ReadAll(txHashResp.Body)
	if err != nil {
		log.Error("Read tx by hash failed", "err", err)
		resp := &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	txResp := &GetTxByHashResp{}
	err = json.Unmarshal(tx, txResp)
	if err != nil {
		log.Error("Unmarshal tx failed", "err", err)
		resp := &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	txStatusUrl := "https://arweave.net/tx/" + hash + "/status"
	txStatusResp, err := http.Get(txStatusUrl)
	if err != nil {
		log.Error("Get tx status failed", "err", err)
		resp := &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	status, err := ioutil.ReadAll(txStatusResp.Body)
	if err != nil {
		log.Error("Read tx status failed", "err", err)
		resp := &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	statusResp := &GetTxStatusResp{}
	err = json.Unmarshal(status, statusResp)
	if err != nil {
		log.Error("Unmarshal tx status failed", "err", err)
		resp := &wallet2.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}
		return resp, err
	}

	froms := make([]*wallet2.Address, 0, 1)
	from := &wallet2.Address{
		Address: txResp.Owner,
	}
	froms = append(froms, from)

	tos := make([]*wallet2.Address, 0, 1)
	to := &wallet2.Address{
		Address: txResp.Target,
	}
	tos = append(tos, to)

	values := make([]*wallet2.Value, 0, 1)
	value := &wallet2.Value{
		Value: txResp.Quantity,
	}
	values = append(values, value)

	resp := &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &wallet2.TxMessage{
			Hash:            txResp.ID,
			Index:           0,
			Froms:           froms,         // owner
			Tos:             tos,           // target
			Values:          values,        // quantity
			Fee:             txResp.Reward, // reward
			Status:          0,
			Type:            0,
			Height:          strconv.Itoa(statusResp.BlockHeight), // in status resp
			ContractAddress: "",
			Datetime:        txResp.Target, // same as target
		},
	}

	return resp, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	panic("implement me")
}
