package account

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
	Id int `json:"id"`
}
