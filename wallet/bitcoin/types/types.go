package types

import "math/big"

type AccountBalance struct {
	FinalBalance  big.Int `json:"final_balance"`
	NTx           big.Int `json:"n_tx"`
	TotalReceived big.Int `json:"total_received"`
}

type UnspentOutputs struct {
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxHash          string `json:"tx_hash"`
	TxOutputN       string `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           int64  `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   int64  `json:"confirmations"`
	TxIndex         int64  `json:"tx_index"`
}
