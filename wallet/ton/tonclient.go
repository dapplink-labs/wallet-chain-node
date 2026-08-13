package ton

import (
	"bytes"
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/savour-labs/wallet-chain-node/config"
)

type tonClient struct {
	client     *liteclient.ConnectionPool
	api        ton.APIClientWrapped
	endpoint   string
	httpClient *http.Client
}

type SendTxResult struct {
	Hash string `json:"message_hash"`
}

type EstimateFeeResult struct {
	InFwdFee   string `json:"in_fwd_fee"`
	StorageFee string `json:"storage_fee"`
	GasFee     string `json:"gas_fee"`
	FwdFee     string `json:"fwd_fee"`
}

type WalletInfoResult struct {
	Balance             string `json:"balance"`
	WalletType          string `json:"wallet_type"`
	Seqno               uint64 `json:"seqno"`
	WalletId            uint64 `json:"wallet_id"`
	LastTransactionLt   string `json:"last_transaction_lt"`
	LastTransactionHash string `json:"last_transaction_hash"`
	Status              string `json:"status"`
}

func (e EstimateFeeResult) SumFees() (int64, error) {
	var sum int64

	// 定义一个帮助函数，将字符串转换为 int64
	parseAndAdd := func(s string) error {
		value, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		sum += value
		return nil
	}

	// 将每个字段转换并相加
	if err := parseAndAdd(e.InFwdFee); err != nil {
		return 0, err
	}
	if err := parseAndAdd(e.StorageFee); err != nil {
		return 0, err
	}
	if err := parseAndAdd(e.GasFee); err != nil {
		return 0, err
	}
	if err := parseAndAdd(e.FwdFee); err != nil {
		return 0, err
	}

	return sum, nil
}

type Transactions struct {
	Account       string `json:"account"`
	Hash          string `json:"hash"`
	Lt            string `json:"lt"`
	Now           int    `json:"now"`
	OrigStatus    string `json:"orig_status"`
	EndStatus     string `json:"end_status"`
	TotalFees     string `json:"total_fees"`
	PrevTransHash string `json:"prev_trans_hash"`
	PrevTransLt   string `json:"prev_trans_lt"`
	Description   struct {
		Type   string `json:"type"`
		Action struct {
			Valid       bool `json:"valid"`
			Success     bool `json:"success"`
			NoFunds     bool `json:"no_funds"`
			ResultCode  int  `json:"result_code"`
			TotActions  int  `json:"tot_actions"`
			MsgsCreated int  `json:"msgs_created"`
			SpecActions int  `json:"spec_actions"`
			TotMsgSize  struct {
				Bits  string `json:"bits"`
				Cells string `json:"cells"`
			} `json:"tot_msg_size"`
			StatusChange    string `json:"status_change"`
			TotalFwdFees    string `json:"total_fwd_fees"`
			SkippedActions  int    `json:"skipped_actions"`
			ActionListHash  string `json:"action_list_hash"`
			TotalActionFees string `json:"total_action_fees"`
		} `json:"action"`
		Aborted  bool `json:"aborted"`
		CreditPh struct {
			Credit string `json:"credit"`
		} `json:"credit_ph"`
		Destroyed bool `json:"destroyed"`
		ComputePh struct {
			Mode             int    `json:"mode"`
			Type             string `json:"type"`
			Success          bool   `json:"success"`
			GasFees          string `json:"gas_fees"`
			GasUsed          string `json:"gas_used"`
			VMSteps          int    `json:"vm_steps"`
			ExitCode         int    `json:"exit_code"`
			GasLimit         string `json:"gas_limit"`
			GasCredit        string `json:"gas_credit"`
			MsgStateUsed     bool   `json:"msg_state_used"`
			AccountActivated bool   `json:"account_activated"`
			VMInitStateHash  string `json:"vm_init_state_hash"`
			VMFinalStateHash string `json:"vm_final_state_hash"`
		} `json:"compute_ph"`
		StoragePh struct {
			StatusChange         string `json:"status_change"`
			StorageFeesCollected string `json:"storage_fees_collected"`
		} `json:"storage_ph"`
		CreditFirst bool `json:"credit_first"`
	} `json:"description"`
	BlockRef struct {
		Workchain int    `json:"workchain"`
		Shard     string `json:"shard"`
		Seqno     int    `json:"seqno"`
	} `json:"block_ref"`
	InMsg struct {
		Hash           string      `json:"hash"`
		Source         string      `json:"source"`
		Destination    string      `json:"destination"`
		Value          string      `json:"value"`
		FwdFee         interface{} `json:"fwd_fee"`
		IhrFee         interface{} `json:"ihr_fee"`
		CreatedLt      interface{} `json:"created_lt"`
		CreatedAt      interface{} `json:"created_at"`
		Opcode         string      `json:"opcode"`
		IhrDisabled    interface{} `json:"ihr_disabled"`
		Bounce         interface{} `json:"bounce"`
		Bounced        interface{} `json:"bounced"`
		ImportFee      string      `json:"import_fee"`
		MessageContent struct {
			Hash    string      `json:"hash"`
			Body    string      `json:"body"`
			Decoded interface{} `json:"decoded"`
		} `json:"message_content"`
		InitState interface{} `json:"init_state"`
	} `json:"in_msg"`
	OutMsgs []struct {
		Hash           string      `json:"hash"`
		Source         string      `json:"source"`
		Destination    string      `json:"destination"`
		Value          string      `json:"value"`
		FwdFee         string      `json:"fwd_fee"`
		IhrFee         string      `json:"ihr_fee"`
		CreatedLt      string      `json:"created_lt"`
		CreatedAt      string      `json:"created_at"`
		Opcode         interface{} `json:"opcode"`
		IhrDisabled    bool        `json:"ihr_disabled"`
		Bounce         bool        `json:"bounce"`
		Bounced        bool        `json:"bounced"`
		ImportFee      interface{} `json:"import_fee"`
		MessageContent struct {
			Hash    string      `json:"hash"`
			Body    string      `json:"body"`
			Decoded interface{} `json:"decoded"`
		} `json:"message_content"`
		InitState interface{} `json:"init_state"`
	} `json:"out_msgs"`
	AccountStateBefore struct {
		Hash          string      `json:"hash"`
		Balance       string      `json:"balance"`
		AccountStatus string      `json:"account_status"`
		FrozenHash    interface{} `json:"frozen_hash"`
		CodeHash      string      `json:"code_hash"`
		DataHash      string      `json:"data_hash"`
	} `json:"account_state_before"`
	AccountStateAfter struct {
		Hash          string      `json:"hash"`
		Balance       string      `json:"balance"`
		AccountStatus string      `json:"account_status"`
		FrozenHash    interface{} `json:"frozen_hash"`
		CodeHash      string      `json:"code_hash"`
		DataHash      string      `json:"data_hash"`
	} `json:"account_state_after"`
	McBlockSeqno int `json:"mc_block_seqno"`
}

type Tx struct {
	Transactions []Transactions `json:"transactions"`
	AddressBook  map[string]struct {
		UserFriendly string `json:"user_friendly"`
	} `json:"address_book"`
}

func (t tonClient) GetLatestBlockHeight() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func newTonClients(conf *config.Config) ([]*tonClient, error) {
	var clients []*tonClient
	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), conf.Fullnode.Ton.RPCs[0].RPCURL)

	if err != nil {
		return nil, err
	}

	client := liteclient.NewConnectionPool()

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return nil, err
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	var endpoint = conf.Fullnode.Ton.TpApiUrl

	clients = append(clients, &tonClient{
		client:     client,
		api:        api,
		endpoint:   endpoint,
		httpClient: new(http.Client),
	})
	return clients, nil
}

func (c *tonClient) GetTxByTxHash(txHash string) (*Tx, error) {
	route := fmt.Sprintf("/transactions?hash=%s", txHash)

	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}

	res := new(Tx)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *tonClient) GetTxByAddr(addr string) (*Tx, error) {
	route := fmt.Sprintf("/transactions?account=%s", addr)

	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}

	res := new(Tx)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *tonClient) PostSendTx(boc string) (*SendTxResult, error) {
	route := fmt.Sprintf("/message")

	payload := struct {
		Boc string `json:"boc"`
	}{Boc: boc}

	ret, err := c.doPost(route, payload)

	if err != nil {
		return nil, err
	}
	res := new(SendTxResult)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *tonClient) WalletInfo(address string) (*WalletInfoResult, error) {
	route := fmt.Sprintf("/walletInformation?address=%s&use_v2=true", address)
	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}
	res := new(WalletInfoResult)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *tonClient) PostEstimateFee(boc string, add string) (*EstimateFeeResult, error) {
	route := fmt.Sprintf("/estimateFee")

	payload := struct {
		Address string `json:"address"`
		Body    string `json:"body"`
	}{Address: add, Body: boc}

	ret, err := c.doPost(route, payload)

	if err != nil {
		return nil, err
	}
	res := new(EstimateFeeResult)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *tonClient) doPost(route string, input interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+route, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// set header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, application/x-bcs")

	// send req
	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("resp status code is not ok, code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithDetails(err, "ioutil.ReadAll fail", "resp", resp, "input", input)
	}

	return body, nil
}

func (c *tonClient) doGet(route string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.endpoint + route)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("resp status code is not ok, code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithDetails(err, "ioutil.ReadAll fail", "resp", resp, "route", route)
	}

	return body, nil
}
