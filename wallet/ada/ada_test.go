package ada

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/savour-labs/wallet-hd-chain/config"
	wallet2 "github.com/savour-labs/wallet-hd-chain/rpc/wallet"
	"github.com/test-go/testify/assert"
	"log"
	"math/big"
	"testing"
)

func TestClient_GetBalance(t *testing.T) {
	//config.yml
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	a, err := NewChainAdaptor(conf)
	assert.NoError(t, err)

	ret, err := a.GetBalance(&wallet2.BalanceRequest{
		Coin:    "stake",
		Address: "",
	})
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetNonce(t *testing.T) {
	//config.yml
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	a, err := NewChainAdaptor(conf)
	assert.NoError(t, err)

	ret, err := a.GetNonce(&wallet2.NonceRequest{
		Address: "",
	})
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetGasPrice(t *testing.T) {
	//config.yml
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	a, err := NewChainAdaptor(conf)
	assert.NoError(t, err)
	// 创建一个 GasPriceRequest 实例
	req := &wallet2.GasPriceRequest{}
	// 调用 GetGasPrice 方法
	resp, err := a.GetGasPrice(req)
	if err != nil {
		t.Fatalf("GetGasPrice failed: %v", err)
	}
	fmt.Println(resp)
}

func TestClient_SendTx(t *testing.T) {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	a, err := NewChainAdaptor(conf)
	assert.NoError(t, err)

	// 构造已签名交易
	signedTx := types.NewTransaction(
		1,                              // nonce
		common.HexToAddress(""),        // to
		big.NewInt(300000000000000000), // value (0.3 PETH)
		22680,                          // gasLimit
		big.NewInt(1),                  // gasPrice
		nil,                            // data
	)
	privateKeyBytes, err := hex.DecodeString("")
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	signedTx, err = types.SignTx(signedTx, types.HomesteadSigner{}, privateKey)
	if err != nil {
	}

	var buf bytes.Buffer
	if err := signedTx.EncodeRLP(&buf); err != nil {
		t.Fatalf("EncodeRLP failed: %v", err)
	}

	// 构造 SendTxRequest
	req := &wallet2.SendTxRequest{
		RawTx: hexutil.Encode(buf.Bytes()),
	}

	// 调用 SendTx 方法
	ret, err := a.SendTx(req)
	if err != nil {
		t.Fatalf("SendTx failed: %v", err)
	}
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetAccountTxFromSignedData(t *testing.T) {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	a, err := NewChainAdaptor(conf)
	assert.NoError(t, err)

	// 将私钥从十六进制字符串转换为字节数组
	privateKeyBytes, err := hex.DecodeString("")
	if err != nil {
		log.Fatal(err)
	}

	// 从字节数组创建 ECDSA 私钥
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	// 创建一个未签名的交易
	unsignedTx := types.NewTransaction(
		1,                              // nonce
		common.HexToAddress(""),        // to
		big.NewInt(300000000000000000), // value (0.3 PETH)
		22680,                          // gasLimit
		big.NewInt(1),                  // gasPrice
		nil,                            // data
	)

	// 签名交易
	signedTx, err := types.SignTx(unsignedTx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 获取已签名的交易数据
	signedTxData, err := signedTx.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Signed transaction data: %x\n", signedTxData)
	ret, err := a.GetAccountTxFromSignedData(&wallet2.TxFromSignedDataRequest{
		SignedTxData: signedTxData,
	})
	assert.NoError(t, err)
	fmt.Println(ret)
}
