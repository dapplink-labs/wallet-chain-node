package algo

type Account struct {
	Address string
}

type TxByAddres struct {
	Address  string
	Page     uint64
	PageSize uint64
	Sort     string
}
