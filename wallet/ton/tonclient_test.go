package ton

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/savour-labs/wallet-chain-node/config"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/xssnick/tonutils-go/address"
	"log"
	"math/big"
	"testing"
)

func TestTonClient_GetAccountBalance(t *testing.T) {
	var client = getClient()
	block, err := client.api.CurrentMasterchainInfo(context.Background())
	acc, err := client.api.GetAccount(context.Background(), block, address.MustParseAddr("EQAUAHcUab66DpOV2GaT_QDuSagpMdIn0x6aMmO3_fPVMyD8"))
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(acc.State.Balance.String())
}

func TestWalletAdaptor_GetAccount(t *testing.T) {
	var client = getClient()
	block, err := client.api.CurrentMasterchainInfo(context.Background())
	acc, err := client.api.GetAccount(context.Background(), block, address.MustParseAddr("EQAUAHcUab66DpOV2GaT_QDuSagpMdIn0x6aMmO3_fPVMyD8"))
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(acc)
}

func TestTonClient_GetTxByAddress(t *testing.T) {
	api := getClient().api
	block, _ := api.CurrentMasterchainInfo(context.Background())
	account, _ := api.GetAccount(context.Background(), block, address.MustParseAddr("UQCQCLTvR0XYTyM0uxh_H8kLAR7u7v98pEKZKpbq8w2per6d"))
	list, _ := api.ListTransactions(context.Background(), address.MustParseAddr("UQCQCLTvR0XYTyM0uxh_H8kLAR7u7v98pEKZKpbq8w2per6d"), 20, account.LastTxLT, account.LastTxHash)
	utils.PrettyPrint(list)
}

func TestTonClient_GetTxByTxHash(t *testing.T) {
	//ret, _ := getClient().GetTxByTxHash("abe0fd4bd6c862855249d1216506c41f4aa420329bc6c15e79f539c1c151beef")
	ret, _ := getClient().GetTxByTxHash("abe0fd4bd6c862855249d1216506c41f4aa420329bc6c15e79f539c1c151beef")
	tx := ret.Transactions[0]
	var from_addrs []*wallet2.Address
	var to_addrs []*wallet2.Address
	totalAmount := big.NewInt(0)
	from_addrs = append(from_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, tx.InMsg.Source)})
	to_addrs = append(to_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, tx.InMsg.Source)})
	if len(tx.InMsg.Value) > 0 {
		totalAmount = new(big.Int).Add(totalAmount, stringToInt(tx.InMsg.Value))
	}

	for _, out := range tx.OutMsgs {
		if len(out.Source) > 0 {
			to_addrs = append(from_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, out.Source)})
		}
		if len(out.Destination) > 0 {
			to_addrs = append(to_addrs, &wallet2.Address{Address: getUserFriendly(ret.AddressBook, out.Destination)})
		}
		totalAmount = new(big.Int).Sub(new(big.Int).Abs(totalAmount), stringToInt(out.Value))
	}
	println(totalAmount)
}

func TestTonClient_GetTxByAddr(t *testing.T) {
	ret, _ := getClient().PostSendTx("te6ccgEBAgEAtQAB4YgAKADuKNN9dB0nK7DNJ/oB3JNQUmOkT6Y9NGTHb/vnqmYAlPJ8R05sBiUrrWj3y47PPho6pClfbMMP76C78XHsAIoIt6BSkCpHRLnhh4jYmt051T80TxUShTkQDseR0e9oWU1NGLs0IRQgAAAAKAAcAQB+QgBh+O1HZtR8ZqFdZxEnUfjbS49VP6piSWpoyCk9p6ufvhzEtAAAAAAAAAAAAAAAAAAAAAAAADExMTgwMzgw")
	println(ret.Hash)

}

func TestTonMessage(t *testing.T) {
	api := getClient().api
	b, err := hex.DecodeString("b0e4eb37bc5929491899d2a50f52f0a4613d3a48e56245267fdecff392ead89b7e4fdf79bf78566b85b73787e5739ab4306350d7ad1adc50be9c57fe2102bfcc")
	w, err := wallet.FromPrivateKey(api, b, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}
	log.Println("master proof checks are completed successfully, now communication is 100% safe!")

	balance, err := w.GetBalance(context.Background(), block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	if balance.Nano().Uint64() >= 3000000 {
		addr := address.MustParseAddr("EQDD8dqOzaj4zUK6ziJOo_G2lx6qf1TEktTRkFJ7T1c_fPQb")

		log.Println("sending transaction and waiting for confirmation...")

		// if destination wallet is not initialized (or you don't care)
		// you should set bounce to false to not get money back.
		// If bounce is true, money will be returned in case of not initialized destination wallet or smart-contract error
		bounce := false

		transfer, err := w.BuildTransfer(addr, tlb.MustFromTON("0.01"), bounce, "11180380")
		if err != nil {
			log.Fatalln("Transfer err:", err.Error())
			return
		}

		ext, err := w.BuildExternalMessageForMany(context.Background(), []*wallet.Message{transfer})

		req, err := tlb.ToCell(ext)

		println(base64.StdEncoding.EncodeToString(req.ToBOCWithFlags(false)))

		//if err != nil {
		//	log.Fatalln("Transfer err:", err.Error())
		//	return
		//}
		//
		//tx, block, err := w.SendWaitTransaction(context.Background(), transfer)
		//if err != nil {
		//	log.Fatalln("SendWaitTransaction err:", err.Error())
		//	return
		//}
		//
		//balance, err = w.GetBalance(context.Background(), block)
		//if err != nil {
		//	log.Fatalln("GetBalance err:", err.Error())
		//	return
		//}
		//
		//log.Printf("transaction confirmed at block %d, hash: %s balance left: %s", block.SeqNo,
		//	base64.StdEncoding.EncodeToString(tx.Hash), balance.String())
		//
		//return
	}

	//log.Println("not enough balance:", balance.String())
}

func getClient() *tonClient {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, _ := config.New(*f)
	var clinets, _ = newTonClients(conf)
	return clinets[0]
}
