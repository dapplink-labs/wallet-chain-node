package types

import "math/big"

type SpendingOutpointsItem struct {
	N       uint64  `json:"n"`
	TxIndex big.Int `json:"tx_index"`
}

type PrevOut struct {
	Addr              string                  `json:"addr"`
	N                 uint64                  `json:"n"`
	Script            string                  `json:"script"`
	SpendingOutpoints []SpendingOutpointsItem `json:"spending_outpoints"`
	Spent             bool                    `json:"spent"`
	TxIndex           big.Int                 `json:"tx_index"`
	Type              uint64                  `json:"type"`
	Value             big.Int                 `json:"value"`
}

type InputItem struct {
	Sequence big.Int `json:"sequence"`
	Witness  string  `json:"witness"`
	Script   string  `json:"script"`
	Index    uint64  `json:"index"`
	PrevOut  PrevOut `json:"prev_out"`
}

type OutItem struct {
	Type              uint64                  `json:"type"`
	Spent             bool                    `json:"spent"`
	Value             big.Int                 `json:"value"`
	SpendingOutpoints []SpendingOutpointsItem `json:"spending_outpoints"`
	N                 uint64                  `json:"n"`
	TxIndex           big.Int                 `json:"tx_index"`
	Script            string                  `json:"script"`
	Addr              string                  `json:"addr"`
}

type TxsItem struct {
	Hash        string      `json:"hash"`
	Ver         uint64      `json:"ver"`
	VinSz       uint64      `json:"vin_sz"`
	VoutSz      uint64      `json:"vout_sz"`
	Size        uint64      `json:"size"`
	Weight      uint64      `json:"weight"`
	Fee         big.Int     `json:"fee"`
	RelayedBy   string      `json:"relayed_by"`
	LockTime    big.Int     `json:"lock_time"`
	TxIndex     uint64      `json:"tx_index"`
	DoubleSpend bool        `json:"double_spend"`
	Time        big.Int     `json:"time"`
	BlockIndex  big.Int     `json:"block_index"`
	BlockHeight big.Int     `json:"block_height"`
	Inputs      []InputItem `json:"inputs"`
	Out         []OutItem   `json:"out"`
	Result      big.Int     `json:"result"`
	Balance     big.Int     `json:"balance"`
}

type Transaction struct {
	Hash160       string    `json:"hash160"`
	Address       string    `json:"address"`
	NTx           uint64    `json:"n_tx"`
	NUnredeemed   big.Int   `json:"n_unredeemed"`
	TotalReceived big.Int   `json:"total_received"`
	TotalSent     big.Int   `json:"total_sent"`
	FinalBalance  big.Int   `json:"final_balance"`
	Txs           []TxsItem `json:"txs"`
}
