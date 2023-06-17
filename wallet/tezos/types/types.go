package types

type AccountBalance struct {
	Chain   string `json:"chain"`
	Coin    string `json:"name"`
	Address string
	Balance string `json:"balance"`
}

type AccountCounter struct {
	Chain   string `json:"chain"`
	Coin    string `json:"name"`
	Address string `json:"addres"`
	Counter string `json:"counter"`
}

type AccountManagerKey struct {
	Chain      string `json:"chain"`
	Coin       string `json:"name"`
	Address    string `json:"addres"`
	ManagerKey string `json:"manager_key"`
}

type Transaction struct {
	Chain  string `json:"chain"`
	Coin   string `json:"name"`
	TxHash string `json:"tx_hash"`
}

type TransactionList struct {
}

type TransactionTxHash struct {
	Id            uint64  `json:"id"`
	Hash          string  `json:"hash"`
	Type          string  `json:"type"`
	Block         string  `json:"block"`
	Time          string  `json:"time"`
	Height        uint64  `json:"height"`
	Cycle         uint64  `json:"cycle"`
	Counter       uint64  `json:"counter"`
	OpN           uint64  `json:"op_n"`
	OpP           uint64  `json:"op_p"`
	Status        string  `json:"status"`
	IsSuccess     bool    `json:"is_success"`
	GasLimit      uint64  `json:"gas_limit"`
	GasUsed       uint64  `json:"gas_used"`
	StorageLimit  uint64  `json:"storage_limit"`
	Volume        float64 `json:"volume"`
	Fee           float64 `json:"fee"`
	Sender        string  `json:"sender"`
	Receiver      string  `json:"receiver"`
	Confirmations uint64  `json:"confirmations"`
}
