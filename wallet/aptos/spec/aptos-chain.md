## 1. docs

developer docs: https://aptos.dev/
api docs: https://aptos.dev/nodes/aptos-api-spec
api explorer: https://fullnode.mainnet.aptoslabs.com/v1/spec#/

## 2. wallet api

### 2.1. get sequence_number for address

api name: v1/accounts/address

example:
```curl
https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12
```

return
```json
{
"sequence_number": "7",
"authentication_key": "0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12"
}
```

### 2.2. Estimate gas price

api name: v1/estimate_gas_price
example:
```curl
https://fullnode.mainnet.aptoslabs.com/v1/estimate_gas_price
```

return
```json
{
    "deprioritized_gas_estimate": 100,
    "gas_estimate": 100,
    "prioritized_gas_estimate": 150
}
```

### 2.3. get account balance


### 2.4. send raw transaction
example:
```curl
curl --request POST \
  --url https://fullnode.mainnet.aptoslabs.com/v1/transactions \
  --header 'Accept: application/json, application/x-bcs' \
  --header 'Content-Type: application/json' \
  --data '{
  "sender": "0x88fbd33f54e1126269769780feb24480428179f552e2313fbe571b72e62a1ca1 ",
  "sequence_number": "32425224034",
  "max_gas_amount": "32425224034",
  "gas_unit_price": "32425224034",
  "expiration_timestamp_secs": "32425224034",
  "payload": {
    "type": "entry_function_payload",
    "function": "0x1::aptos_coin::transfer",
    "type_arguments": [
      "string"
    ],
    "arguments": [
      null
    ]
  },
  "signature": {
    "type": "ed25519_signature",
    "public_key": "0x88fbd33f54e1126269769780feb24480428179f552e2313fbe571b72e62a1ca1 ",
    "signature": "0x88fbd33f54e1126269769780feb24480428179f552e2313fbe571b72e62a1ca1 "
  }
}'
```
Please read the document directly, it is relatively simple, so I won’t explain it here


### 2.5. get transaction by address
api name: v1/accounts/address/transactions
param: address, limit, start
Please read the document directly, it is relatively simple, so I won’t explain it here

### 2.6. get transaction by hash
api name: 1/transactions/by_hash/txn_hash
Please read the document directly, it is relatively simple, so I won’t explain it here

### 2.7. get transaction by version
api name: v1/transactions/by_version/txn_version
Please read the document directly, it is relatively simple, so I won’t explain it here