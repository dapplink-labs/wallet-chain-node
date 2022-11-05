package near

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Transaction struct {
	TransactionHash              string `json:"transaction_hash" db:"transaction_hash"`
	IncludedInBlockHash          string `json:"included_in_block_hash" db:"included_in_block_hash"`
	IncludedInChunkHash          string `json:"included_in_chunk_hash" db:"included_in_chunk_hash"`
	IndexInChunk                 string `json:"index_in_chunk" db:"index_in_chunk"`
	BlockTimestamp               string `json:"block_timestamp" db:"block_timestamp"`
	SignerAccountId              string `json:"signer_account_id" db:"signer_account_id"`
	SignerPublicKey              string `json:"signer_public_key" db:"signer_public_key"`
	Nonce                        string `json:"nonce" db:"nonce"`
	ReceiverAccountId            string `json:"receiver_account_id" db:"receiver_account_id"`
	Signature                    string `json:"signature" db:"signature"`
	Status                       string `json:"status" db:"status"`
	ConvertedIntoReceiptId       string `json:"converted_into_receipt_id" db:"converted_into_receipt_id"`
	ReceiptConversionGasBurnt    string `json:"receipt_conversion_gas_burnt" db:"receipt_conversion_gas_burnt"`
	ReceiptConversionTokensBurnt string `json:"receipt_conversion_tokens_burnt" db:"receipt_conversion_tokens_burnt"`
}

type BlockTransaction struct {
	Transaction
	BlockHeight string `json:"block_height"`
	Amount      string `json:"amount"`
}

type TransactionActionArgs struct {
	Gas        int64  `json:"gas"`
	Deposit    string `json:"deposit"`
	ArgsBase64 string `json:"args_base64"`
	MethodName string `json:"method_name"`
}

type TransactionAction struct {
	TransactionHash    string `json:"transaction_hash" db:"transaction_hash"`
	ActionKind         string `json:"action_kind" db:"action_kind"`
	IndexInTransaction int64  `json:"index_in_transaction" db:"index_in_transaction"`
	Args               string `json:"args" db:"args"`
}

type Block struct {
	BlockHeight     string `json:"block_height" db:"block_height"`
	BlockHash       string `json:"block_hash" db:"block_hash"`
	PrevBlockHash   string `json:"prev_block_hash" db:"prev_block_hash"`
	BlockTimestamp  string `json:"block_timestamp" db:"block_timestamp"`
	TotalSupply     string `json:"total_supply" db:"total_supply"`
	GasPrice        string `json:"gas_price" db:"gas_price"`
	AuthorAccountId string `json:"author_account_id" db:"author_account_id"`
}

// var db *sqlx.DB
var db *sql.DB

// var db *gorm.DB
var err error

type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	fmt.Printf("> %s %q", query, args)
	return context.WithValue(ctx, "begin", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	fmt.Printf(". took: %s\n", time.Since(begin))
	return ctx, nil
}

func init() {
	//db, err = sqlx.Connect("postgres", "host=mainnet.db.explorer.indexer.near.dev user=public_readonly password=nearprotocol dbname=mainnet_explorer sslmode=disable binary_parameters=yes")

	db, err = sql.Open("postgres", "host=mainnet.db.explorer.indexer.near.dev user=public_readonly password=nearprotocol dbname=mainnet_explorer sslmode=disable binary_parameters=yes")
	//defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	//db.SetConnMaxLifetime(10)
	//sql.Register("postgres", sqlhooks.Wrap(&pq.Driver{}, &Hooks{}))
	//db, err = gorm.Open(postgres.Open(
	//	"postgres://public_readonly:nearprotocol@mainnet.db.explorer.indexer.near.dev/mainnet_explorer?binary_parameters=yes",
	//), &gorm.Config{
	//	PrepareStmt: true,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//db.Exec("SELECT NOW() ")

}
