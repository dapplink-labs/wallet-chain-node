package icp

import (
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/ic"
	"github.com/aviate-labs/agent-go/principal"
)

// client for the "ledger" canister.
type LedgerClient struct {
	agentClient *agent.Agent
	canisterId  principal.Principal
}

// NewAgent creates a new agent for the "ledger" canister.
func NewLedgerClient(config agent.Config) (*LedgerClient, error) {
	agentClient, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	return &LedgerClient{
		agentClient: agentClient,
		canisterId:  ic.LEDGER_PRINCIPAL,
	}, nil
}

// QueryBlocks calls the "query_blocks" method on the "ledger" canister.
func (ledgerClient LedgerClient) QueryBlocks(getBlockReq GetBlocksArgs) (*QueryBlocksResponse, error) {
	var resp QueryBlocksResponse
	if err := ledgerClient.agentClient.Query(
		ledgerClient.canisterId,
		"query_blocks",
		[]any{getBlockReq},
		[]any{&resp},
	); err != nil {
		return nil, err
	}
	return &resp, nil
}
