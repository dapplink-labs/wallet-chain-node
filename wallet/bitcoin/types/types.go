package types

import "math/big"

type AccountBalance struct {
	FinalBalance  big.Int `json:"final_balance"`
	NTx           big.Int `json:"n_tx"`
	TotalReceived big.Int `json:"total_received"`
}

type UnspentOutput struct {
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxHash          string `json:"tx_hash"`
	TxOutputN       uint64 `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           uint64 `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   uint64 `json:"confirmations"`
	TxIndex         uint64 `json:"tx_index"`
}

type UnspentOutputList struct {
	Notice         string          `json:"notice"`
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

type GasFee struct {
	ChainFullName       string `json:"chainFullName"`
	ChainShortName      string `json:"chainShortName"`
	Symbol              string `json:"symbol"`
	BestTransactionFee  string `json:"bestTransactionFee"`
	RecommendedGasPrice string `json:"recommendedGasPrice"`
	RapidGasPrice       string `json:"rapidGasPrice"`
	StandardGasPrice    string `json:"standardGasPrice"`
	SlowGasPrice        string `json:"slowGasPrice"`
}

type GasFeeData struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []GasFee `json:"data"`
}
