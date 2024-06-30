package icp

import (
	"fmt"
	"github.com/aviate-labs/agent-go"
	"net/url"
	"testing"
)

var ic0, _ = url.Parse("https://ic0.app/")

func Test_icp_status(t *testing.T) {
	c := agent.NewClient(agent.ClientConfig{Host: ic0})
	status, _ := c.Status()
	fmt.Println(status.Version)
	fmt.Println(status)
}

//func Test_icp_QueryBlocks_heignt(t *testing.T) {
//	// 我们还不知道最后一个区块是什么，因此首先要查询区块高度。
//	ledgerClient, _ := NewLedgerClient(agent.DefaultConfig)
//
//	blockHeight, err := ledgerClient.QueryBlocks(GetBlocksArgs{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	//我们可以查询分类账的第一个区块。
//	oldestBlock := blockHeight.FirstBlockIndex
//	fmt.Println("oldestBlock ", oldestBlock)
//	//我们可以查询分类账的最后一个区块。
//	lastBlock := blockHeight.ChainLength
//	fmt.Println("lastBlock ", lastBlock)
//}
//
//func Test_icp_QueryBlocks(t *testing.T) {
//	// 我们还不知道最后一个区块是什么，因此首先要查询区块高度。
//	ledgerClient, _ := NewLedgerClient(agent.DefaultConfig)
//	blockHeight, err := ledgerClient.QueryBlocks(GetBlocksArgs{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	//我们可以查询分类账的最后一个区块。
//	lastBlock := blockHeight.ChainLength
//
//	// Query the last 10 blocks.
//	response, err := ledgerClient.QueryBlocks(GetBlocksArgs{
//		Start:  lastBlock - 10,
//		Length: 10,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for i, block := range response.Blocks {
//		operation := block.Transaction.Operation
//		if transfer := operation.Transfer; transfer != nil {
//			var from principal.AccountIdentifier
//			copy(from[:], transfer.From)
//
//			var to principal.AccountIdentifier
//			copy(to[:], transfer.To)
//
//			fmt.Printf("Block %d: %s -> %s: %.2f ICP.\n", int(lastBlock)+i, from, to, float64(transfer.Amount.E8s)/1e8)
//		} else if burn := operation.Burn; burn != nil {
//			var from principal.AccountIdentifier
//			copy(from[:], burn.From)
//
//			fmt.Printf("Block %d: %s: %.2f ICP burned.\n", int(lastBlock)+i, from, float64(burn.Amount.E8s)/1e8)
//		} else if mint := operation.Mint; mint != nil {
//			var to principal.AccountIdentifier
//			copy(to[:], mint.To)
//
//			fmt.Printf("Block %d: %s: %.2f ICP minted.\n", int(lastBlock)+i, to, float64(mint.Amount.E8s)/1e8)
//		}
//	}
//}
