package near

import (
	"bytes"
	"encoding/json"
	"github.com/SavourDao/savour-hd/config"
	"github.com/ethereum/go-ethereum/params"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"
)

type nearClient struct {
	RpcClient        rpc.RpcClient
	Client           client.Client
	nodeConfig       config.Node
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

func newLocalClient(network config.NetWorkType) *nearClient {
	var endpoint string
	switch network {
	case config.MainNet:
		endpoint = rpc.MainnetRPCEndpoint
	case config.TestNet:
		endpoint = rpc.TestnetRPCEndpoint
	default:
		panic("unsupported network type")
	}
	rpcClient := rpc.NewRpcClient(endpoint)
	return &nearClient{
		RpcClient: rpcClient,
	}
}

type GetBlockReq struct {
	Finality string `json:"finality"`
}
type GetBlockResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Author string `json:"author"`
		Chunks []struct {
			BalanceBurnt         string        `json:"balance_burnt"`
			ChunkHash            string        `json:"chunk_hash"`
			EncodedLength        int           `json:"encoded_length"`
			EncodedMerkleRoot    string        `json:"encoded_merkle_root"`
			GasLimit             int64         `json:"gas_limit"`
			GasUsed              int64         `json:"gas_used"`
			HeightCreated        int           `json:"height_created"`
			HeightIncluded       int           `json:"height_included"`
			OutcomeRoot          string        `json:"outcome_root"`
			OutgoingReceiptsRoot string        `json:"outgoing_receipts_root"`
			PrevBlockHash        string        `json:"prev_block_hash"`
			PrevStateRoot        string        `json:"prev_state_root"`
			RentPaid             string        `json:"rent_paid"`
			ShardID              int           `json:"shard_id"`
			Signature            string        `json:"signature"`
			TxRoot               string        `json:"tx_root"`
			ValidatorProposals   []interface{} `json:"validator_proposals"`
			ValidatorReward      string        `json:"validator_reward"`
		} `json:"chunks"`
		Header struct {
			Approvals             []interface{} `json:"approvals"`
			BlockMerkleRoot       string        `json:"block_merkle_root"`
			BlockOrdinal          int           `json:"block_ordinal"`
			ChallengesResult      []interface{} `json:"challenges_result"`
			ChallengesRoot        string        `json:"challenges_root"`
			ChunkHeadersRoot      string        `json:"chunk_headers_root"`
			ChunkMask             []bool        `json:"chunk_mask"`
			ChunkReceiptsRoot     string        `json:"chunk_receipts_root"`
			ChunkTxRoot           string        `json:"chunk_tx_root"`
			ChunksIncluded        int           `json:"chunks_included"`
			EpochID               string        `json:"epoch_id"`
			EpochSyncDataHash     interface{}   `json:"epoch_sync_data_hash"`
			GasPrice              string        `json:"gas_price"`
			Hash                  string        `json:"hash"`
			Height                int           `json:"height"`
			LastDsFinalBlock      string        `json:"last_ds_final_block"`
			LastFinalBlock        string        `json:"last_final_block"`
			LatestProtocolVersion int           `json:"latest_protocol_version"`
			NextBpHash            string        `json:"next_bp_hash"`
			NextEpochID           string        `json:"next_epoch_id"`
			OutcomeRoot           string        `json:"outcome_root"`
			PrevHash              string        `json:"prev_hash"`
			PrevHeight            int           `json:"prev_height"`
			PrevStateRoot         string        `json:"prev_state_root"`
			RandomValue           string        `json:"random_value"`
			RentPaid              string        `json:"rent_paid"`
			Signature             string        `json:"signature"`
			Timestamp             int64         `json:"timestamp"`
			TimestampNanosec      string        `json:"timestamp_nanosec"`
			TotalSupply           string        `json:"total_supply"`
			ValidatorProposals    []interface{} `json:"validator_proposals"`
			ValidatorReward       string        `json:"validator_reward"`
		} `json:"header"`
	} `json:"result"`
	ID int `json:"id"`
}

func (c *nearClient) Request() error {
	jsonStr := []byte(`{"jsonrpc": "2.0", "method": "block", "params": {"finality": "final"},"id": 1}`)
	url := "https://rpc.mainnet.near.org"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := GetBlockResult{}
	_ = json.Unmarshal(body, &res)
	return nil
}

//func (c *nearClient) GetLatestBlockHeight() (int, error) {
//	rpcClient := jsonrpc.NewClient("https://rpc.mainnet.near.org")
//	ret, err := rpcClient.Call(context.Background(), "block", &GetBlockReq{Finality: "final"})
//	if err != nil {
//		return 0, err
//	}
//	result, ok := (ret.Result).(GetBlockResult)
//	if !ok {
//		return 0, err
//	}
//	return result.Result.Header.Height, nil
//}

func (c *nearClient) GetLatestBlockHeight() (int, error) {
	jsonStr := []byte(`{"jsonrpc": "2.0", "method": "block", "params": {"finality": "final"},"id": 1}`)
	url := c.nodeConfig.RPCs[0].RPCURL
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := GetBlockResult{}
	_ = json.Unmarshal(body, &res)
	return res.Result.Header.Height, nil
}

func newNearClients(conf *config.Config) ([]*nearClient, error) {
	endpoint := rpc.DevnetRPCEndpoint
	if conf.NetWork == "testnet" {
		endpoint = rpc.TestnetRPCEndpoint
	} else if conf.NetWork == "mainnet" {
		endpoint = rpc.MainnetRPCEndpoint
	}
	var clients []*nearClient
	rpcClient := rpc.NewRpcClient(endpoint)
	clients = append(clients, &nearClient{
		RpcClient:  rpcClient,
		nodeConfig: conf.Fullnode.Near,
	})
	return clients, nil
}
