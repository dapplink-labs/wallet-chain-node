package types

type RpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}
type GetBlockResult struct {
	Result struct {
		LedgerCurrentIndex int    `json:"ledger_current_index"`
		Status             string `json:"status"`
	} `json:"result"`
}

type GetBalance struct {
	Strict      bool   `json:"strict"`
	Account     string `json:"account"`
	LedgerIndex string `json:"ledger_index"`
	Queue       bool   `json:"queue"`
}
type GetBalanceRes struct {
	Forwarded bool `json:"forwarded"`
	Result    struct {
		AccountData struct {
			Account           string `json:"Account"`
			Balance           string `json:"Balance"`
			Flags             int    `json:"Flags"`
			LedgerEntryType   string `json:"LedgerEntryType"`
			OwnerCount        int    `json:"OwnerCount"`
			PreviousTxnID     string `json:"PreviousTxnID"`
			PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
			Sequence          int    `json:"Sequence"`
			Index             string `json:"index"`
		} `json:"account_data"`
		LedgerCurrentIndex int `json:"ledger_current_index"`
		QueueData          struct {
			TxnCount int `json:"txn_count"`
		} `json:"queue_data"`
		Status    string `json:"status"`
		Validated bool   `json:"validated"`
	} `json:"result"`
	Type     string `json:"type"`
	Warnings []struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	} `json:"warnings"`
}
