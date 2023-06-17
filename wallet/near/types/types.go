package types

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/savour-labs/wallet-hd-chain/wallet/near/keys"
)

type RpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

type GetAccessKeyRequest struct {
	Finality    string `json:"finality"`
	RequestType string `json:"request_type"`
	AccountID   string `json:"account_id"`
	PublicKey   string `json:"public_key"`
}

type GetAccessKeyResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockHash   string `json:"block_hash"`
		BlockHeight int    `json:"block_height"`
		Nonce       int64  `json:"nonce"`
		Permission  string `json:"permission"`
	} `json:"result"`
	Error struct {
		Name  string `json:"name"`
		Cause struct {
			Name string `json:"name"`
			Info struct {
				ErrorMessage string `json:"error_message"`
			} `json:"info"`
		} `json:"cause"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	ID int `json:"id"`
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
type SendTxResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		ReceiptsOutcome []struct {
			BlockHash string `json:"block_hash"`
			Id        string `json:"id"`
			Outcome   struct {
				ExecutorId string        `json:"executor_id"`
				GasBurnt   int64         `json:"gas_burnt"`
				Logs       []interface{} `json:"logs"`
				Metadata   struct {
					GasProfile []interface{} `json:"gas_profile"`
					Version    int           `json:"version"`
				} `json:"metadata"`
				ReceiptIds []string `json:"receipt_ids"`
				Status     struct {
					SuccessValue string `json:"SuccessValue"`
				} `json:"status"`
				TokensBurnt string `json:"tokens_burnt"`
			} `json:"outcome"`
			Proof []struct {
				Direction string `json:"direction"`
				Hash      string `json:"hash"`
			} `json:"proof"`
		} `json:"receipts_outcome"`
		Status struct {
			SuccessValue string `json:"SuccessValue"`
		} `json:"status"`
		Transaction struct {
			Actions []struct {
				Transfer struct {
					Deposit string `json:"deposit"`
				} `json:"Transfer"`
			} `json:"actions"`
			Hash       string `json:"hash"`
			Nonce      int64  `json:"nonce"`
			PublicKey  string `json:"public_key"`
			ReceiverId string `json:"receiver_id"`
			Signature  string `json:"signature"`
			SignerId   string `json:"signer_id"`
		} `json:"transaction"`
		TransactionOutcome struct {
			BlockHash string `json:"block_hash"`
			Id        string `json:"id"`
			Outcome   struct {
				ExecutorId string        `json:"executor_id"`
				GasBurnt   int64         `json:"gas_burnt"`
				Logs       []interface{} `json:"logs"`
				Metadata   struct {
					GasProfile interface{} `json:"gas_profile"`
					Version    int         `json:"version"`
				} `json:"metadata"`
				ReceiptIds []string `json:"receipt_ids"`
				Status     struct {
					SuccessReceiptId string `json:"SuccessReceiptId"`
				} `json:"status"`
				TokensBurnt string `json:"tokens_burnt"`
			} `json:"outcome"`
			Proof []struct {
				Direction string `json:"direction"`
				Hash      string `json:"hash"`
			} `json:"proof"`
		} `json:"transaction_outcome"`
	} `json:"result"`
	Id    int `json:"id"`
	Error any `json:"error"`
}

type GetBalanceRes struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Amount        string `json:"amount"`
		BlockHash     string `json:"block_hash"`
		BlockHeight   int    `json:"block_height"`
		CodeHash      string `json:"code_hash"`
		Locked        string `json:"locked"`
		StoragePaidAt int    `json:"storage_paid_at"`
		StorageUsage  int    `json:"storage_usage"`
	} `json:"result"`
	ID int `json:"id"`
}

type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

type GetBalanceParam struct {
	Finality    string `json:"finality"`
	RequestType string `json:"request_type"`
	AccountID   string `json:"account_id"`
}

type Config struct {
	Signer    keys.KeyPair
	NetworkID string
	RPCClient *rpc.Client
}
