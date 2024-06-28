package icp

import (
	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
	"net/url"
	"time"
)

// icp0 is the default host for the Internet Computer.
var icp0, _ = url.Parse("https://icp0.io/")

// Config is the configuration for an Agent.
type IcpConfig struct {
	Identity      Identity
	IngressExpiry time.Duration
	ClientConfig  *ClientConfig
	FetchRootKey  bool
}

// Identity is an identity that can sign messages.
type Identity interface {
	// Sender方法返回与Identity相关联的主体（principal）的引用。在DFINITY的互联网计算机架构中，
	// 主体是一个全局唯一标识符，
	// 通常用于标识一个用户、智能合约或其他网络参与者。通过返回主体，我们可以知道是谁发送或签署了消息。
	// Sender returns the principal of the identity.
	Sender() []byte
	// Sign signs the given message.
	// Sign方法接受一个字节切片msg作为输入，代表需要签署的消息。然后，该方法使用与Identity相关联的私钥对消息进行加密签名，
	// 并返回签名的字节切片。这个签名可以被任何人使用相应的公钥来验证，以确保消息的完整性和来源。
	Sign(msg []byte) []byte
	// PublicKey returns the public key of the identity.
	PublicKey() []byte
}

type WalletAdaptor struct {
	fallback.WalletAdaptor
	clients *Client
}

func NewIcpControllersNode() *config.IcpControllersNode {
	return &config.IcpControllersNode{
		RegistryCanister:      "rwlgt-iiaaa-aaaaa-aaaaa-cai",
		GovernanceCanister:    "rrkah-fqaaa-aaaaa-aaaaq-cai",
		LedgerCanister:        "ryjl3-tyaaa-aaaaa-aaaba-cai",
		NnsRootCanister:       "r7inp-6aaaa-aaaaa-aaabq-cai",
		CyclesMintingCanister: "rkp4c-7iaaa-aaaaa-aaaca-cai",
		LifelineCanister:      "rno2w-sqaaa-aaaaa-aaacq-cai",
		GenesisTokenCanister:  "renrk-eyaaa-aaaaa-aaada-cai",
		SnsWasmCanister:       "qaa6y-5yaaa-aaaaa-aaafa-cai",
		IdentityCanister:      "rdmx6-jaaaa-aaaaa-aaadq-cai",
		NnsUiCanister:         "qoctq-giaaa-aaaaa-aaaea-cai",
	}
}
