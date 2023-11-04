package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/savour-labs/wallet-chain-node/config"
)

const (
	omniPrefix = "6f6d6e69"
)

type btcClient struct {
	*rpcclient.Client
	chainConfig *chaincfg.Params
	compressed  bool
}

func newBtcClients(conf *config.Config) ([]*btcClient, error) {
	chainConfig := &chaincfg.TestNet3Params
	if conf.NetWork == "mainnet" {
		chainConfig = &chaincfg.MainNetParams
	} else if conf.NetWork == "regtest" {
		chainConfig = &chaincfg.RegressionNetParams
	}
	log.Info("btc client setup", "network", conf.NetWork)

	var clients []*btcClient
	for _, rpc := range conf.Fullnode.Btc.RPCs {
		client, err := rpcclient.New(&rpcclient.ConnConfig{
			HTTPPostMode: true,
			DisableTLS:   true,
			Host:         rpc.RPCURL,
			User:         rpc.RPCUser,
			Pass:         rpc.RPCPass,
		}, nil)
		if err != nil {
			log.Error("Fail to create BTC client", "err", err)
			continue
		}
		clients = append(clients, &btcClient{
			Client:      client,
			chainConfig: chainConfig,
			compressed:  true,
		})
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}

	return clients, nil
}

func newLocalBtcClient(network config.NetWorkType) *btcClient {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "",
		User:         "",
		Pass:         "",
	}, nil)

	if err != nil {
		panic("Fail to create BTC client")

	}
	var para *chaincfg.Params
	switch network {
	case config.MainNet:
		para = &chaincfg.MainNetParams
	case config.TestNet:
		para = &chaincfg.TestNet3Params
	case config.RegTest:
		para = &chaincfg.RegressionNetParams
	default:
		panic("unsupported network type")
	}
	return &btcClient{
		Client:      client,
		chainConfig: para,
		compressed:  true,
	}
}

// Close close the client connection
func (btc *btcClient) Close() {
	btc.Shutdown()
}

// GetNetwork get the current bitcoin network
func (btc *btcClient) GetNetwork() *chaincfg.Params {
	return btc.chainConfig
}

type EstimateSmartFeeResult struct {
	Feerate float64  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  int      `json:"blocks"`
}

// EstimateSmartFee provides an estimated fee  in bitcoins per kilobyte.
func (btc *btcClient) EstimateSmartFee(numBlocks int64) (EstimateSmartFeeResult, error) {
	var reply = EstimateSmartFeeResult{}

	params, err := marshal(numBlocks)
	if err != nil {
		return reply, err
	}

	data, err := btc.RawRequest("estimatesmartfee", params)
	if err != nil {
		return reply, errors.Wrap(err, "could not estimate fee")
	}

	err = json.Unmarshal(data, &reply)
	if err != nil {
		return reply, errors.Wrap(err, "could not unmarshal estimate smart fee result")
	}

	if len(reply.Errors) != 0 {
		return reply, errors.New(strings.Join(reply.Errors, ","))
	}

	return reply, nil
}

func marshal(numBlocks int64) ([]json.RawMessage, error) {
	var params []json.RawMessage
	numBlocksJSON, err := json.Marshal(numBlocks)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal numBlocks")
	}
	params = []json.RawMessage{numBlocksJSON}
	return params, nil
}

func (btc *btcClient) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	networkInfo, err := btc.GetNetworkInfo()
	defaultFeeRate := "0.02000000"
	if err != nil {
		log.Warn("failed to get btc networkinfo, use latest api")
		return btc.SendRawTransaction190001(tx, defaultFeeRate)
	}

	if networkInfo.Version >= 190001 {
		return btc.SendRawTransaction190001(tx, defaultFeeRate)
	}
	return btc.Client.SendRawTransaction(tx, false)
}

func (btc *btcClient) SendRawTransaction190001(tx *wire.MsgTx, maxFeeRateInBtcPerK string) (*chainhash.Hash, error) {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return nil, err
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}

	var params []json.RawMessage
	maxFeeRateJSON, err := json.Marshal(maxFeeRateInBtcPerK)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal numBlocks")
	}
	txHexJSON, err := json.Marshal(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal numBlocks")
	}
	params = []json.RawMessage{txHexJSON, maxFeeRateJSON}

	data, err := btc.RawRequest("sendrawtransaction", params)
	if err != nil {
		return nil, err
	}

	var txHashStr string
	err = json.Unmarshal(data, &txHashStr)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHashStr)
}

type GetBlockVerboseResult struct {
	Hash          string                 `json:"hash"`
	Confirmations int64                  `json:"confirmations"`
	StrippedSize  int32                  `json:"strippedsize"`
	Size          int32                  `json:"size"`
	Weight        int32                  `json:"weight"`
	Height        int64                  `json:"height"`
	Version       int32                  `json:"version"`
	VersionHex    string                 `json:"versionHex"`
	MerkleRoot    string                 `json:"merkleroot"`
	Tx            []*btcjson.TxRawResult `json:"tx,omitempty"`
	Time          int64                  `json:"time"`
	Nonce         uint32                 `json:"nonce"`
	Bits          string                 `json:"bits"`
	Difficulty    float64                `json:"difficulty"`
	PreviousHash  string                 `json:"previousblockhash"`
	NextHash      string                 `json:"nextblockhash,omitempty"`
}

func (btc *btcClient) GetBlockWithRawTransactionVerbose(blockHash *chainhash.Hash) (*GetBlockVerboseResult, error) {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	hashJSON, err := json.Marshal(hash)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal hash")
	}

	verboseJSON, err := json.Marshal(2)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal verbose")
	}

	params := []json.RawMessage{hashJSON, verboseJSON}

	data, err := btc.RawRequest("getblock", params)
	if err != nil {
		return nil, err
	}

	var result GetBlockVerboseResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type GetNetworkInfoResult struct {
	Version         int32                          `json:"version"`
	SubVersion      string                         `json:"subversion"`
	ProtocolVersion int32                          `json:"protocolversion"`
	LocalServices   string                         `json:"localservices"`
	LocalRelay      bool                           `json:"localrelay"`
	TimeOffset      int64                          `json:"timeoffset"`
	Connections     int32                          `json:"connections"`
	NetworkActive   bool                           `json:"networkactive"`
	Networks        []btcjson.NetworksResult       `json:"networks"`
	RelayFee        float64                        `json:"relayfee"`
	IncrementalFee  float64                        `json:"incrementalfee"`
	LocalAddresses  []btcjson.LocalAddressesResult `json:"localaddresses"`
	Warnings        string                         `json:"warnings"`
}

func (btc *btcClient) GetNetworkInfo() (*GetNetworkInfoResult, error) {
	data, err := btc.RawRequest("getnetworkinfo", nil)
	if err != nil {
		return nil, err
	}

	var result GetNetworkInfoResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (btc *btcClient) GetLatestBlockHeight() (int64, error) {
	return btc.Client.GetBlockCount()
}
