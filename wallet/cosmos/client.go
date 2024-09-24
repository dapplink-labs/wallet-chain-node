package cosmos

import (
	"context"
	"errors"
	"fmt"

	"github.com/armon/go-metrics"
	"google.golang.org/grpc"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"

	ed255192 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Client struct {
	grpcConn        *grpc.ClientConn
	bankClient      banktypes.QueryClient
	txServiceClient tx.ServiceClient
	authClient      authv1beta1.QueryClient
	bankKeeper      keeper.BaseKeeper
}

func NewClient(rpcTarget string) (*Client, error) {
	grpcConn, err := grpc.Dial(rpcTarget, grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(nil)))
	if err != nil {
		return nil, err
	}

	return &Client{
		grpcConn:        grpcConn,
		bankClient:      banktypes.NewQueryClient(grpcConn),
		txServiceClient: tx.NewServiceClient(grpcConn),
		authClient:      authv1beta1.NewQueryClient(grpcConn),
	}, nil
}

func (c *Client) Close() error {
	// close grpc conn
	return c.grpcConn.Close()
}

// GetBalance get balance with coin name and address
func (c *Client) GetBalance(ctx context.Context, coin, addr string) (*types.Coin, error) {
	address, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	resp, err := c.bankClient.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   coin,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetBalance(), nil
}

type Tx struct {
	MsgIndex int        `json:"msg_index"`
	Events   []*TxEvent `json:"events"`
}

type TxEvent struct {
	Type       string              `json:"type"`
	Attributes []*TxEventAttribute `json:"attributes"`
}

type TxEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetTxByHash get tx by hash
func (c *Client) GetTxByHash(ctx context.Context, hash string) (*tx.GetTxResponse, error) {
	return c.txServiceClient.GetTx(ctx, &tx.GetTxRequest{Hash: hash})
}

// GetTxByEvent get tx by event
func (c *Client) GetTxByEvent(ctx context.Context, event []string, page, limit uint64) (*tx.GetTxsEventResponse, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 1000
	}

	if len(event) == 0 {
		return nil, errors.New("event cannot be empty")
	}
	eventTmp := make([]string, 0)
	for _, v := range event {
		_ = v
		eventTmp = append(eventTmp, fmt.Sprintf("message.sender='%s'", "cosmos1qd99t24whd3hfg22r53x8uw9ps3rrctwxqvn4m"))
	}

	tmp := &tx.GetTxsEventRequest{
		Events: eventTmp,
		Pagination: &query.PageRequest{
			Offset: page,
			Limit:  limit,
		},
	}

	return c.txServiceClient.GetTxsEvent(ctx, tmp)
}

func (c *Client) GetAccount(ctx context.Context, addr string) (*authv1beta1.QueryAccountResponse, error) {
	return c.authClient.Account(ctx, &authv1beta1.QueryAccountRequest{Address: addr})
}

func (c *Client) SendTx(ctx context.Context, fromAddr, toAddr, coin string, amount int64) (*banktypes.MsgSendResponse, error) {
	newCtx := types.UnwrapSDKContext(ctx)

	fromAddress, err := types.AccAddressFromBech32(fromAddr)
	if err != nil {
		return nil, err
	}
	toAddress, err := types.AccAddressFromBech32(toAddr)
	if err != nil {
		return nil, err
	}

	if c.bankKeeper.BlockedAddr(toAddress) {
		return nil, errors.New(fmt.Sprintf("%s is not allowed to receive funds", toAddress))
	}

	msg := banktypes.NewMsgSend(fromAddress, toAddress, types.NewCoins(types.NewInt64Coin(coin, amount)))
	if err := c.bankKeeper.SendCoins(newCtx, fromAddress, toAddress, msg.Amount); err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "send"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	return &banktypes.MsgSendResponse{}, nil
}

func (c *Client) GetAddressFromPubKey(key []byte) string {
	// todo check
	pub := ed255192.PubKey{Key: key}
	return pub.Address().String()
}

func (c *Client) BroadcastTx(ctx context.Context, txByte []byte) (*tx.BroadcastTxResponse, error) {
	return c.txServiceClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		TxBytes: txByte,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
}
