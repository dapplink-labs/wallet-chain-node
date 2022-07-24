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

func newLocalSolanaClient(network config.NetWorkType) *solanaClient {
	var endpoint string
	switch network {
	case config.MainNet:
		endpoint = rpc.MainnetRPCEndpoint
	case config.TestNet:
		endpoint = rpc.TestnetRPCEndpoint
	default:
		panic("unsupported network type")
	}
	rpcClient := rpc.NewRpcClient(endpoint)

	return &solanaClient{
		RpcClient: rpcClient,
	}
}

func (sol *solanaClient) GetLatestBlockHeight() (int64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return int64(res.Result), nil
}

func newSolanaClients(conf *config.Config) ([]*solanaClient, error) {
	endpoint := rpc.DevnetRPCEndpoint
	if conf.NetWork == "testnet" {
		endpoint = rpc.TestnetRPCEndpoint
	} else if conf.NetWork == "mainnet" {
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

func (sol *solanaClient) GetBalance(address string) (string, error) {
	balance, err := sol.RpcClient.GetBalanceWithConfig(
		context.TODO(),
		address,
		rpc.GetBalanceConfig{
			Commitment: rpc.CommitmentProcessed,
		},
	)
	if err != nil {
		log.Fatalf("failed to get balance with cfg, err: %v", err)
		return "", err
	}

	var lamportsOnAccount = new(big.Float).SetUint64(balance.Result.Value)
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(1000000000))
	return solBalance.String(), nil
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

func (sol *solanaClient) GetTxByAddress(address string, page uint32, size uint32) ([]GetTxByAddressTx, error) {
	offset := (page - 1) * size
	url := sol.solanaConfig.PublicUrl + "/account/solTransfers?limit=" + strconv.FormatInt(int64(size), 10) + "&account=" + address + "&offset=" + strconv.FormatInt(int64(offset), 10)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var res GetTxByAddressRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res.Data, nil
}

func (sol *solanaClient) GetAccount() (string, string, error) {
	account := types.NewAccount()
	address := account.PublicKey.ToBase58()
	private := base58.Encode(account.PrivateKey)
	return address, private, nil
}

type Header struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}
type Instructions struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIDIndex int    `json:"programIdIndex"`
}
type Message struct {
	AccountKeys     []string       `json:"accountKeys"`
	Header          Header         `json:"header"`
	Instructions    []Instructions `json:"instructions"`
	RecentBlockhash string         `json:"recentBlockhash"`
}
type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

type TxMessage struct {
	Hash   string
	From   string
	To     string
	Fee    string
	Status bool
	Value  string
	Type   int32
	Height string
}

func (sol *solanaClient) GetTxByHash(hash string) (*TxMessage, error) {
	out, err := sol.RpcClient.GetTransaction(
		context.TODO(),
		hash,
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
		return nil, err
	}
	message := out.Result.Transaction.(map[string]interface{})["message"]
	accountKeys := message.((map[string]interface{}))["accountKeys"].([]interface{})
	signatures := out.Result.Transaction.(map[string]interface{})["signatures"].([]interface{})
	_hash := signatures[0]
	if out.Result.Meta.Err != nil || len(out.Result.Meta.LogMessages) == 0 || _hash == "" {
		log.Fatalf("not found tx, err: %v", err)
		return nil, err
	}
	var txMessage []*TxMessage
	for i := 0; i < len(accountKeys); i++ {
		to := accountKeys[i].(string)
		amount := out.Result.Meta.PostBalances[i] - out.Result.Meta.PreBalances[i]

		if to != "" && amount > 0 {
			txMessage = append(txMessage, &TxMessage{
				Hash:   hash,
				From:   "",
				To:     to,
				Fee:    strconv.FormatUint(out.Result.Meta.Fee, 10),
				Status: true,
				Value:  strconv.FormatInt(amount, 10),
				Type:   1,
				Height: strconv.FormatUint(out.Result.Slot, 10),
			})
		}
	}

	for i := 0; i < len(out.Result.Meta.PostTokenBalances); i++ {
		postToken := out.Result.Meta.PostTokenBalances[i]

		preTokenBalance := getPreTokenBalance(out.Result.Meta.PreTokenBalances, postToken.AccountIndex)
		if preTokenBalance == nil {
			continue
		}
		postAmount, _ := strconv.ParseFloat(postToken.UITokenAmount.Amount, 64)
		preAmount, _ := strconv.ParseFloat(preTokenBalance.UITokenAmount.Amount, 64)
		amount := postAmount - preAmount
		if amount > 0 {
			txMessage = append(txMessage, &TxMessage{
				Hash:   hash,
				From:   "",
				To:     postToken.Owner,
				Fee:    strconv.FormatUint(out.Result.Meta.Fee, 10),
				Status: true,
				Value:  strconv.FormatFloat(amount, 'E', -1, 10),
				Type:   1,
				Height: strconv.FormatUint(out.Result.Slot, 10),
			})
		}
	}
	if len(txMessage) > 0 {
		return txMessage[0], nil
	}
	log.Fatalf("not found tx, err: %v", err)
	return nil, err
}

func getPreTokenBalance(preTokenBalance []rpc.TransactionMetaTokenBalance, accountIndex uint64) *rpc.TransactionMetaTokenBalance {
	for j := 0; j < len(preTokenBalance); j++ {
		preToken := preTokenBalance[j]
		if preToken.AccountIndex == accountIndex {
			return &preTokenBalance[j]
		}
	}
	return nil
}
func (sol *solanaClient) RequestAirdrop(address string) (string, error) {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	sig, err := c.RequestAirdrop(
		context.TODO(),
		address,
		1e9, // lamports (1 SOL = 10^9 lamports)
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
		return "", err
	}
	return sig, nil
}

func (sol *solanaClient) SendTx(pri string) (string, error) {

	nonceAccountPubkey := common.PublicKeyFromString(sol.solanaConfig.NonceAccountAddr)
	nonceAccount, err := sol.Client.GetNonceAccount(context.Background(), nonceAccountPubkey.ToBase58())
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
		return "", err
	}
	var feePayer, _ = types.AccountFromBase58(sol.solanaConfig.FeeAccountPriKey)
	var userAccount, _ = types.AccountFromBase58(pri)

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer, userAccount},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: nonceAccount.Nonce.ToBase58(),
			Instructions: []types.Instruction{
				sysprog.AdvanceNonceAccount(sysprog.AdvanceNonceAccountParam{
					Nonce: nonceAccountPubkey,
					Auth:  userAccount.PublicKey,
				}),
				memoprog.BuildMemo(memoprog.BuildMemoParam{
					Memo: []byte("use nonce"),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
		return "", err
	}

	sig, err := sol.Client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
		return "", err
	}

	return sig, nil
}

func (sol *solanaClient) GetNonce() (string, error) {
	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), sol.solanaConfig.NonceAccountAddr)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
		return "", err
	}
	return nonce, nil
}

func (sol *solanaClient) GetMinRent() (string, error) {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		log.Fatalf("failed to get GetMinimumBalanceForRentExemption , err: %v", err)
		return "", err

	}
	return strconv.FormatUint(bal.Result, 10), nil
}
