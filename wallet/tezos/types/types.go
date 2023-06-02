package types

type AccountBalance struct {
	Chain   string `json:"chain"`
	Coin    string `json:"name"`
	Addres  string
	Balance string `json:"balance"`
}

type AccountCounter struct {
	Chain   string `json:"chain"`
	Coin    string `json:"name"`
	Addres  string `json:"addres"`
	Counter string `json:"counter"`
}

type AccountManagerKey struct {
	Chain      string `json:"chain"`
	Coin       string `json:"name"`
	Addres     string `json:"addres"`
	ManagerKey string `json:"manager_key"`
}

type Transaction struct {
	Chain  string `json:"chain"`
	Coin   string `json:"name"`
	TxHash string `json:"tx_hash"`
}
