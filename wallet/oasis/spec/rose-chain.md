## 1.document

- official document：https://docs.oasis.dev/general
- official website：https://oasisprotocol.org/
- Paratime explorer：https://explorer.emerald.oasis.dev
- Paratime explorer：https://testnet.explorer.emerald.oasis.dev/
- Oasis explorer：https://www.oasisscan.com/
- Rosetta api: https://www.rosetta-api.org/docs/AccountApi.html

## 2. transaction fee rule

### 2.1. Calculation of handling fees

At present, staking.tranfer transfers transaction transfer gas, but it does not actually consume
A word of caution: staking your ROSE is a different transaction than sending them! When you stake your tokens (stake.escrow transaction), you can withdraw them at any time. staking.Transfer On the other hand, sending your tokens (transaction) means that the recipient will own the tokens and you cannot get them back yourself.

### 2.2. About the role of fee and Gas

fee: An optional fee that the caller promises to pay for executing the transaction. Consensus layer transfer can choose to fill in 0. The transfer between the consensus layer and ParaTime needs to consume fee and gas. The deposit (transfer to paratime) transaction is currently free, and a handling fee may be incurred later, and the withdrawal (withdrawal from paratime) transaction is not free, and the handling fee will be deducted from your paratime balance. Similar to the ETH consumption model here, a transaction that burns tokens needs to be sent.

## 3. wallet api

### 3.1. query account balance
oasis api
rosetta api
```azure
curl --location --request POST 'https://rosetta.oasis.dev/account/balance' \
--header 'Content-Type: application/json' \
--header 'Cookie: __cflb=02DiuDkfxbUyQCDbAMscAxT996N9Xw8YZHJhXe8Qi2J6g' \
--data-raw '{
    "network_identifier": {
        "blockchain": "Oasis",
        "network": "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535"
    },
    "account_identifier": {
        "address": "oasis1qpg2xuz46g53737343r20yxeddhlvc2ldqsjh70p"
    }
}'
```
return
```azure
{
    "block_identifier": {
        "index": 14170142,
        "hash": "53408fda9b0a8dde2e69f51dfb782b7aba314dc32662b450e87133186c67e2a3"
    },
    "balances": [
        {
            "value": "1661899675965425099",
            "currency": {
                "symbol": "ROSE",
                "decimals": 9
            }
        }
    ],
    "metadata": {
        "nonce": 111
    }
}
```

### 3.2. query account nonce
oasis api
/oasis-core.Consensus/GetSignerNonce
rosetta api

please see balance api
```azure
 "metadata": {
    "nonce": 111
  }
```


### 3.3. query estimate gas fee
oasis api
/oasis-core.Consensus/EstimateGas
rosetta api


### 3.4. query chain context
oasis api
oasis-core.Consensus/GetChainContext
rosetta api
```azure
curl --location --request POST 'https://rosetta.oasis.dev/network/list' \
--header 'Content-Type: application/json' \
--header 'Cookie: __cflb=02DiuDkfxbUyQCDbAMscAxT996N9Xw8YZHJhXe8Qi2J6g' \
--data-raw '{
    "metadata": {}
}'
```
chainid 

### 3.5. send raw transaction
oasis api
/oasis-core.RuntimeClient/SubmitTx
rosetta api
```azure
curl --location --request POST 'https://rosetta.oasis.dev/construction/submit' \
--header 'Content-Type: application/json' \
--header 'Cookie: __cflb=02DiuDkfxbUyQCDbAMscAxT996N9Xw8YZHJhXe8Qi2J6g' \
--data-raw '{
    "network_identifier": {
        "blockchain": "Oasis",
        "network": "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535"
    },
    "signed_transaction": "0x00000000000000"
}'
```

### 3.6. query transaction recordlist
use explorer api
please look it from explorer

### 3.7. query transaction detail
use explorer api
please look it from explorer

### 3.8. rosetta api for other

#### 3.8.1 fetch support network list

```azure
curl --location --request POST 'https://rosetta.oasis.dev/network/list' \
--header 'Content-Type: application/json' \
--header 'Cookie: __cflb=02DiuDkfxbUyQCDbAMscAxT996N9Xw8YZHJhXe8Qi2J6g' \
--data-raw '{
 "metadata": {}
}'

```




