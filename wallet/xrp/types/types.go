package types

type RpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}
type GetBlockResult struct {
	Result struct {
		LedgerCurrentIndex int    `json:"ledger_current_index"`
		Status             string `json:"status"`
	} `json:"result"`
}

type GetBalance struct {
	Strict      bool   `json:"strict"`
	Account     string `json:"account"`
	LedgerIndex string `json:"ledger_index"`
	Queue       bool   `json:"queue"`
}

type GetTx struct {
	Transaction string `json:"transaction"`
}

type SendTx struct {
	TxBlob string `json:"tx_blob"`
}

type GetBalanceRes struct {
	Forwarded bool `json:"forwarded"`
	Result    struct {
		AccountData struct {
			Account           string `json:"Account"`
			Balance           string `json:"Balance"`
			Flags             int    `json:"Flags"`
			LedgerEntryType   string `json:"LedgerEntryType"`
			OwnerCount        int    `json:"OwnerCount"`
			PreviousTxnID     string `json:"PreviousTxnID"`
			PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
			Sequence          int    `json:"Sequence"`
			Index             string `json:"index"`
		} `json:"account_data"`
		LedgerCurrentIndex int `json:"ledger_current_index"`
		QueueData          struct {
			TxnCount int `json:"txn_count"`
		} `json:"queue_data"`
		Status    string `json:"status"`
		Validated bool   `json:"validated"`
	} `json:"result"`
	Type     string `json:"type"`
	Warnings []struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	} `json:"warnings"`
}

type GetAccountTx struct {
	Account string `json:"account"`
	Binary  bool   `json:"binary"`
	Forward bool   `json:"forward"`
	//LedgerIndexMax int    `json:"ledger_index_max"`
	//LedgerIndexMin int    `json:"ledger_index_min"`
}

type Transaction struct {
	Amount      string `json:"amount"`
	Fee         string `json:"fee"`
	To          string `json:"to"`
	From        string `json:"from"`
	Hash        string `json:"hash"`
	BlockHeight string `json:"block_height"`
}
type RequestBody struct {
	Method string `json:"method"`
	Params any    `json:"params"`
}

type SendTxResult struct {
	Result struct {
		TxJson struct {
			Hash string `json:"hash"`
		} `json:"tx_json"`
	}
}
type GetTxByHashResult struct {
	Result struct {
		Account string `json:"Account"`
		Amount  struct {
			Currency string `json:"currency"`
			Issuer   string `json:"issuer"`
			Value    string `json:"value"`
		} `json:"Amount"`
		Destination     string `json:"Destination"`
		Fee             string `json:"Fee"`
		Flags           int64  `json:"Flags"`
		Sequence        int    `json:"Sequence"`
		SigningPubKey   string `json:"SigningPubKey"`
		TransactionType string `json:"TransactionType"`
		TxnSignature    string `json:"TxnSignature"`
		Date            int    `json:"date"`
		Hash            string `json:"hash"`
		InLedger        int    `json:"inLedger"`
		LedgerIndex     int    `json:"ledger_index"`
		Meta            struct {
			AffectedNodes []struct {
				ModifiedNode struct {
					FinalFields struct {
						Account    string `json:"Account"`
						Balance    string `json:"Balance"`
						Flags      int    `json:"Flags"`
						OwnerCount int    `json:"OwnerCount"`
						Sequence   int    `json:"Sequence"`
					} `json:"FinalFields"`
					LedgerEntryType string `json:"LedgerEntryType"`
					LedgerIndex     string `json:"LedgerIndex"`
					PreviousFields  struct {
						Balance  string `json:"Balance"`
						Sequence int    `json:"Sequence"`
					} `json:"PreviousFields"`
					PreviousTxnID     string `json:"PreviousTxnID"`
					PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
				} `json:"ModifiedNode"`
			} `json:"AffectedNodes"`
			TransactionIndex  int    `json:"TransactionIndex"`
			TransactionResult string `json:"TransactionResult"`
			DeliveredAmount   struct {
				Currency string `json:"currency"`
				Issuer   string `json:"issuer"`
				Value    string `json:"value"`
			} `json:"delivered_amount"`
		} `json:"meta"`
		Status    string `json:"status"`
		Validated bool   `json:"validated"`
		Warnings  []struct {
			ID      int    `json:"id"`
			Message string `json:"message"`
		} `json:"warnings"`
	} `json:"result"`
}

type GetAccountResult struct {
	Result struct {
		Account        string `json:"account"`
		LedgerIndexMax int    `json:"ledger_index_max"`
		LedgerIndexMin int    `json:"ledger_index_min"`
		Limit          int    `json:"limit"`
		Marker         struct {
			Ledger int `json:"ledger"`
			Seq    int `json:"seq"`
		} `json:"marker"`
		Status       string `json:"status"`
		Transactions []struct {
			Meta struct {
				//AffectedNodes []struct {
				//	ModifiedNode struct {
				//		FinalFields struct {
				//			Account    string `json:"Account"`
				//			Balance    string `json:"Balance"`
				//			Flags      int    `json:"Flags"`
				//			OwnerCount int    `json:"OwnerCount"`
				//			Sequence   int    `json:"Sequence"`
				//		} `json:"FinalFields"`
				//		LedgerEntryType string `json:"LedgerEntryType"`
				//		LedgerIndex     string `json:"LedgerIndex"`
				//		PreviousFields  struct {
				//			Balance  string `json:"Balance"`
				//			Sequence int    `json:"Sequence"`
				//		} `json:"PreviousFields"`
				//		PreviousTxnID     string `json:"PreviousTxnID"`
				//		PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
				//	} `json:"ModifiedNode,omitempty"`
				//	CreatedNode struct {
				//		LedgerEntryType string `json:"LedgerEntryType"`
				//		LedgerIndex     string `json:"LedgerIndex"`
				//		NewFields       struct {
				//			Account  string `json:"Account"`
				//			Balance  string `json:"Balance"`
				//			Sequence int    `json:"Sequence"`
				//		} `json:"NewFields"`
				//	} `json:"CreatedNode,omitempty"`
				//} `json:"AffectedNodes"`
				TransactionIndex  int    `json:"TransactionIndex"`
				TransactionResult string `json:"TransactionResult"`
				DeliveredAmount   string `json:"delivered_amount"`
			} `json:"meta,omitempty"`
			Tx struct {
				Account         string `json:"Account"`
				Amount          string `json:"Amount"`
				Destination     string `json:"Destination"`
				DestinationTag  int    `json:"DestinationTag"`
				Fee             string `json:"Fee"`
				Flags           int64  `json:"Flags"`
				Sequence        int    `json:"Sequence"`
				SigningPubKey   string `json:"SigningPubKey"`
				TransactionType string `json:"TransactionType"`
				TxnSignature    string `json:"TxnSignature"`
				Date            int    `json:"date"`
				Hash            string `json:"hash"`
				InLedger        int    `json:"inLedger"`
				LedgerIndex     int    `json:"ledger_index"`
			} `json:"tx,omitempty"`
			Validated bool `json:"validated"`
		}
	} `json:"result"`
}
