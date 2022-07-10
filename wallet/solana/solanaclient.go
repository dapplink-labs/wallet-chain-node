package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SavourDao/savour-core/config"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memoprog"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
)

type solanaClient struct {
	RpcClient        rpc.RpcClient
	Client           client.Client
	solanaConfig     config.SolanaNode
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

//func newLocalSolanaClient(network config.NetWorkType) *solanaClient {
//	var para *params.ChainConfig
//	switch network {
//	case config.MainNet:
//		para = params.MainnetChainConfig
//	case config.TestNet:
//		para = params.RopstenChainConfig
//	case config.RegTest:
//		para = params.AllCliqueProtocolChanges
//	default:
//		panic("unsupported network type")
//	}
//	return &solanaClient{
//		Client:           &solanaClient,
//		chainConfig:      para,
//		cacheBlockNumber: nil,
//		local:            true,
//	}
//}

func (sol *solanaClient) GetLatestBlockHeight() (int64, error) {
	return 0, nil
}
func newSolanaClients(conf *config.Config) ([]*solanaClient, error) {
	endpoint := rpc.DevnetRPCEndpoint
	if conf.Fullnode.Sol.NetWork == "testnet" {
		endpoint = rpc.TestnetRPCEndpoint
	} else if conf.Fullnode.Sol.NetWork == "mainnet" {
		endpoint = rpc.MainnetRPCEndpoint
	}
	var clients []*solanaClient
	rpcClient := rpc.NewRpcClient(endpoint)
	clients = append(clients, &solanaClient{
		RpcClient:    rpcClient,
		solanaConfig: conf.Fullnode.Sol,
	})
	return clients, nil
}

func (sol *solanaClient) GetBalance(address string) string {
	balance, err := sol.RpcClient.GetBalanceWithConfig(
		context.TODO(),
		address,
		rpc.GetBalanceConfig{
			Commitment: rpc.CommitmentProcessed,
		},
	)
	if err != nil {
		log.Fatalf("failed to get balance with cfg, err: %v", err)
	}

	var lamportsOnAccount = new(big.Float).SetUint64(balance.Result.Value)
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(1000000000))
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
	account := types.NewAccount()
	fmt.Println(account.PublicKey.ToBase58())
	fmt.Println(base58.Encode(account.PrivateKey))

	address := account.PublicKey.ToBase58()
	private := base58.Encode(account.PrivateKey)
	return address, private, nil
}

func (sol *solanaClient) GetTxByHash(hash string) rpc.GetTransactionResponse {
	out, err := sol.RpcClient.GetTransaction(
		context.TODO(),
		hash,
	)
	if err != nil {
		panic(err)
	}
	return out
}

func (sol *solanaClient) RequestAirdrop(address string) {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	sig, err := c.RequestAirdrop(
		context.TODO(),
		address,
		1e9, // lamports (1 SOL = 10^9 lamports)
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
	}
	fmt.Println(sig)
}

// FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
var feePayer, _ = types.AccountFromBase58("4TMFNY9ntAn3CHzguSAvDNLPRoQTaK3sWbQQXdDXaE6KWRBLufGL6PJdsD2koiEe3gGmMdRK3aAw7sikGNksHJrN")

// 9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde
var alice, _ = types.AccountFromBase58("4voSPg3tYuWbKzimpQK9EbXHmuyy5fUrtXvpLDMLkmY6TRncaTHAKGD8jUg3maB5Jbrd9CkQg4qjJMyN6sQvnEF2")

func (sol *solanaClient) SendTx() {

	// get nonce account
	nonceAccountPubkey := common.PublicKeyFromString("DJyNpXgggw1WGgjTVzFsNjb3fuQZVMqhoakvSBfX9LYx")
	nonceAccount, err := sol.Client.GetNonceAccount(context.Background(), nonceAccountPubkey.ToBase58())
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
	}

	// create a tx
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer, alice},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: nonceAccount.Nonce.ToBase58(),
			Instructions: []types.Instruction{
				sysprog.AdvanceNonceAccount(sysprog.AdvanceNonceAccountParam{
					Nonce: nonceAccountPubkey,
					Auth:  alice.PublicKey,
				}),
				memoprog.BuildMemo(memoprog.BuildMemoParam{
					Memo: []byte("use nonce"),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	sig, err := sol.Client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println("txhash", sig)
}

func (sol *solanaClient) GetNonce(address string) string {
	nonceAccountAddr := "DJyNpXgggw1WGgjTVzFsNjb3fuQZVMqhoakvSBfX9LYx"
	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), nonceAccountAddr)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
	}

	fmt.Println("nonce", nonce)
	return nonce
}

func (sol *solanaClient) GetMinRent() string {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
	}
	fmt.Println("nonce", bal.Result)
	return strconv.FormatUint(bal.Result, 10)
}
