package solana

import (
	"context"
	"encoding/base64"
	"github.com/mr-tron/base58"
	"log"
	"math/big"
	"strconv"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"

	"github.com/savour-labs/wallet-chain-node/config"
)

type SolanaClient struct {
	RpcClient    rpc.RpcClient
	Client       *client.Client
	solanaConfig config.SolanaNode
}

func NewSolanaClients(conf *config.Config) (*SolanaClient, error) {
	endpoint := conf.Fullnode.Sol.RPCs[0].RPCURL
	rpcClient := rpc.NewRpcClient(endpoint)
	clientNew := client.NewClient(endpoint)
	return &SolanaClient{
		RpcClient:    rpcClient,
		Client:       clientNew,
		solanaConfig: conf.Fullnode.Sol,
	}, nil
}

func (sol *SolanaClient) GetLatestBlockHeight() (int64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return int64(res.Result), nil
}

func (sol *SolanaClient) GetBalance(address string) (string, error) {
	balance, err := sol.RpcClient.GetBalanceWithConfig(
		context.TODO(),
		address,
		rpc.GetBalanceConfig{
			Commitment: rpc.CommitmentProcessed,
		},
	)
	if err != nil {
		return "", err
	}

	var lamportsOnAccount = new(big.Float).SetUint64(balance.Result.Value)

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

func (sol *SolanaClient) GetAccount() (string, string, error) {
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

func (sol *SolanaClient) GetTxByHash(hash string) (*TxMessage, error) {
	out, err := sol.RpcClient.GetTransaction(
		context.Background(),
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
func (sol *SolanaClient) RequestAirdrop(address string) (string, error) {
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

func (sol *SolanaClient) SendTx(rawTx string) (string, error) {
	res, err := sol.RpcClient.SendTransactionWithConfig(
		context.Background(),
		base64.StdEncoding.EncodeToString([]byte(rawTx)),
		rpc.SendTransactionConfig{
			Encoding: rpc.SendTransactionConfigEncodingBase64,
		},
	)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", res.Error
	}
	return res.Result, nil
}

func (sol *SolanaClient) GetNonce(nonceAccountAddr string) (string, error) {
	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), nonceAccountAddr)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
		return "", err
	}
	return nonce, nil
}

func (sol *SolanaClient) GetMinRent() (string, error) {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		log.Fatalf("failed to get GetMinimumBalanceForRentExemption , err: %v", err)
		return "", err
	}
	return strconv.FormatUint(bal.Result, 10), nil
}
