package dot

type RuntimeVersionResponse struct {
	SpecName           string  `json:"specName"`           //链名称
	SpecVersion        float64 `json:"specVersion"`        //版本
	TransactionVersion float64 `json:"transactionVersion"` //事物版本
}

type Transaction struct {
	SpecName           string  `json:"specName"`           //链名称
	SpecVersion        float64 `json:"specVersion"`        //版本
	TransactionVersion float64 `json:"transactionVersion"` //事物版本
}
