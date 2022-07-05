package solana

import (
	"context"
	"github.com/SavourDao/savour-core/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"math/big"
	"sync"
)

type solanaClient struct {
	*rpc.Client
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

// newSolanaClients init the eth client
func newSolanaClients(conf *config.Config) ([]*solanaClient, error) {
	endpoint := rpc.MainNetBeta_RPC
	var clients []*solanaClient

	client := rpc.New(endpoint)
	clients = append(clients, &solanaClient{
		Client: client,
	})
	return clients, nil
}

func (sol *solanaClient) GetLatestBlockHeight() string {
	pubKey := solana.MustPublicKeyFromBase58("7xLk17EQQ5KLDLDe44wCmupJKJjTGd8hs3eSVVhCx932")
	out, err := sol.Client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}
	spew.Dump(out)
	spew.Dump(out.Value) // total lamports on the account; 1 sol = 1000000000 lamports

	var lamportsOnAccount = new(big.Float).SetUint64(uint64(out.Value))
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
	return solBalance.String()
}
