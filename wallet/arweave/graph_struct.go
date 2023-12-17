package arweave

type TransactionList struct {
	Transactions Transactions `json:"transactions"`
}
type TransactionDetail struct {
	Transaction Transaction `json:"transaction"`
}
type Block struct {
	ID        string `json:"id"`
	Height    int64  `json:"height"`
	Timestamp int64  `json:"timestamp"`
}
type Fee struct {
	Ar string `json:"ar"`
}
type Owner struct {
	Address string `json:"address"`
}
type Quantity struct {
	Ar string `json:"ar"`
}
type Transaction struct {
	Block     Block    `json:"block"`
	Fee       Fee      `json:"fee"`
	ID        string   `json:"id"`
	Owner     Owner    `json:"owner"`
	Quantity  Quantity `json:"quantity"`
	Recipient string   `json:"recipient"`
}
type Edges struct {
	Cursor      string      `json:"cursor"`
	Transaction Transaction `json:"node"`
}
type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
}
type Transactions struct {
	Edges    []Edges  `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
}
