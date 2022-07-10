package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SavourDao/savour-core/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"
)

type solanaClient struct {
	*rpc.Client
	solanaConfig     config.SolanaNode
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

// newSolanaClients init the eth client
func newSolanaClients(conf *config.Config) ([]*solanaClient, error) {
	//endpoint := rpc.MainNetBeta_RPC
	endpoint := rpc.DevNet_RPC
	var clients []*solanaClient

	client := rpc.New(endpoint)
	clients = append(clients, &solanaClient{
		Client:       client,
		solanaConfig: conf.Fullnode.Sol,
	})

	return clients, nil
}

func (sol *solanaClient) GetBalance(address string) string {
	pubKey := solana.MustPublicKeyFromBase58(address)
	out, err := sol.Client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}

	var lamportsOnAccount = new(big.Float).SetUint64(out.Value)
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
	return solBalance.String()
}

type GetTxByAddressRes struct {
	Data []GetTxByAddressTx
}

type GetTxByAddressTx struct {
	ID                  string `json:"_id"`
	Src                 string `json:"src"`
	Dst                 string `json:"dst"`
	Lamport             int    `json:"lamport"`
	BlockTime           int    `json:"blockTime"`
	Slot                int    `json:"slot"`
	TxHash              string `json:"txHash"`
	Fee                 int    `json:"fee"`
	Status              string `json:"status"`
	Decimals            int    `json:"decimals"`
	TxNumberSolTransfer int    `json:"txNumberSolTransfer"`
}

func (sol *solanaClient) GetTxByAddress(address string) []GetTxByAddressTx {
	url := sol.solanaConfig.PublicUrl + "/account/solTransfers?limit=20&account=" + address
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var res GetTxByAddressRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err)
	}
	return res.Data
}

func (sol *solanaClient) GetAccount() (string, string, error) {
	account := solana.NewWallet()
	address := account.PublicKey().String()
	private := account.PrivateKey.String()
	return address, private, nil
}

func (sol *solanaClient) GetTxByHash(hash string) *rpc.GetTransactionResult {
	txSig := solana.MustSignatureFromBase58(hash)
	out, err := sol.Client.GetTransaction(
		context.TODO(),
		txSig,
		&rpc.GetTransactionOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		panic(err)
	}
	_, err = solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
	if err != nil {
		panic(err)
	}
	return out
}

func (sol *solanaClient) RequestAirdrop(address string) {
	publicKey, _ := solana.PublicKeyFromBase58(address)

	if true {
		// Airdrop 5 sol to the account so it will have something to transfer:
		out, err := sol.Client.RequestAirdrop(
			context.TODO(),
			publicKey,
			solana.LAMPORTS_PER_SOL*1,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			panic(err)
		}
		fmt.Println("airdrop transaction signature:", out)
		time.Sleep(time.Second * 5)
	}
}

func (sol *solanaClient) SendTx() {

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	if err != nil {
		panic(err)
	}

	recent, err := sol.Client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	// The public key of the account that you will send sol TO:
	accountTo := solana.MustPublicKeyFromBase58("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	// The amount to send (in lamports);
	// 1 sol = 1000000000 lamports
	amount := uint64(3333)
	accountFrom, err := solana.PrivateKeyFromBase58("2qbHWS73PBBhmn3RusnQYRjsBax1ZLuieNXTyAtNp6ZkeRSJ8BaDcheM73FWwE4vg2cvxHG6pcvnZzVojUVGgecH")
	//accountFrom, err := solana.PrivateKeyFromSolanaKeygenFile("/path/to/.config/solana/id.json")

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				accountFrom.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		panic(err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}
	spew.Dump(tx)
	// Pretty print the transaction:
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		sol.Client,
		wsClient,
		tx,
	)
	if err != nil {
		panic(err)
	}
	spew.Dump(sig)
}
