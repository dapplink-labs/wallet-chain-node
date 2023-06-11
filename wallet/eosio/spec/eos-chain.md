## 1.document

- official document：https://developers.eos.io
- official website：https://eos.io/
- blockchain explorer：https://bloks.io/
- api: https://developers.eos.io/welcome/latest/reference/index
- github: https://github.com/eosio

endpoint:
api.eosnewyork.io
eos.greymass.com
eos.greymass.com:443 (secured endpoint)
public.eosinfra.io

## 2. transaction fee rule
use eos powerup

## 3. wallet api

### 3.1.query chain info
v1/chain/get_info

### 3.2. query account info
/v1/chain/get_account
=> getBalance 
=> done

### 3.3. send raw transaction
/v1/chain/push_transaction
=> SendTx
=> 

### 3.4. query transaction recordlist
v1/history/get_actions
=> GetTxByAddress

### 3.5. query transaction detail
/v1/history/get_transaction
=> GetTxByHash

### 3.6. abi_json_to_bin and abi_bin_to_json
v1/chain/abi_json_to_bin
v1/chain/abi_bin_to_json

### 3.7. account active
https://github.com/EOSIO/eosjs/blob/master/docs/how-to-guides/05_how-to-create-an-account.md
api detail: https://developers.eos.io/manuals/eos/latest/nodeos/plugins/chain_api_plugin/api-reference/index#operation/abi_bin_to_json


