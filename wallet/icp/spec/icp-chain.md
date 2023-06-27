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
Request example:
```curl
curl --location --request POST 'https://rosetta-api.internetcomputer.org/account/balance' \
--header 'Content-Type: application/json' \
--data-raw '{
    "network_identifier": {
        "blockchain": "Internet Computer",
        "network": "00000000000000020101"
    },
    "account_identifier": {
        "address": "220c3a33f90601896e26f76fa619fe288742df1fa75426edfaf759d39f2455a5"
    }
}'
```
Return:
```json
{
    "block_identifier": {
        "index": 6363276,
        "hash": "ee6fb57a4bd8ea32a53a70d2d3fa8876e14cbdbbf3cbfdb7f44f0141324e4d0e"
    },
    "balances": [
        {
            "value": "161276942136012",
            "currency": {
                "symbol": "ICP",
                "decimals": 8
            }
        }
    ]
}
```

### 3.2. query account nonce
Request example:

Return:

### 3.3. query estimate gas fee
Request example:
```curl
curl --location --request POST 'https://rosetta-api.internetcomputer.org/construction/metadata' \
--header 'Content-Type: application/json' \
--data-raw '{
    "network_identifier": {
        "blockchain": "Internet Computer",
        "network": "00000000000000020101"
    }
}'
```

Return:
```json
{
    "metadata": {},
    "suggested_fee": [
        {
            "value": "10000",
            "currency": {
                "symbol": "ICP",
                "decimals": 8
            }
        }
    ]
}
```

### 3.4. query chain context
Request example:

Return:

### 3.5. send raw transaction
Request example:
```curl
curl --location --request POST 'https://rosetta-api.internetcomputer.org/construction/submit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "network_identifier": {
        "blockchain": "Internet Computer",
        "network": "00000000000000020101"
    },
    "transaction_identifier": {
        "signed_transaction": "8888"
    }
}'
```

Return:
```json
{
    "code": 701,
    "message": "Invalid request",
    "retriable": false,
    "details": {
        "error_message": "Json deserialize error: missing field `signed_transaction` at line 9 column 1"
    }
}
```

### 3.6. query transaction recordlist
use explorer api
please look it from explorer

### 3.7. query transaction detail
Request example:
```curl
curl --location --request POST 'https://rosetta-api.internetcomputer.org/search/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
    "network_identifier": {
        "blockchain": "Internet Computer",
        "network": "00000000000000020101"
    },
    "transaction_identifier": {
        "hash": "82ea41477a0d646bc53dc5cdaf0e232f87d32ea4772e41660d10e892c19c885e"
    }
}'
```

Return:
```json
{
    "transactions": [
        {
            "block_identifier": {
                "index": 6365500,
                "hash": "4e4756719d7f49630afdf8a0ec03ab0cba173492996b0178426764f36db5ca6c"
            },
            "transaction": {
                "transaction_identifier": {
                    "hash": "82ea41477a0d646bc53dc5cdaf0e232f87d32ea4772e41660d10e892c19c885e"
                },
                "operations": [
                    {
                        "operation_identifier": {
                            "index": 0
                        },
                        "type": "TRANSACTION",
                        "status": "COMPLETED",
                        "account": {
                            "address": "e3d9cc1cb18253ba58aeff4fe2922527f788461190a6ef1d7eaf719a4cb7f135"
                        },
                        "amount": {
                            "value": "-292499021",
                            "currency": {
                                "symbol": "ICP",
                                "decimals": 8
                            }
                        }
                    },
                    {
                        "operation_identifier": {
                            "index": 1
                        },
                        "type": "TRANSACTION",
                        "status": "COMPLETED",
                        "account": {
                            "address": "bb3357cba483f268d044d4bbd4639f82c16028a76eebdf62c51bc11fc918d278"
                        },
                        "amount": {
                            "value": "292499021",
                            "currency": {
                                "symbol": "ICP",
                                "decimals": 8
                            }
                        }
                    },
                    {
                        "operation_identifier": {
                            "index": 2
                        },
                        "type": "FEE",
                        "status": "COMPLETED",
                        "account": {
                            "address": "e3d9cc1cb18253ba58aeff4fe2922527f788461190a6ef1d7eaf719a4cb7f135"
                        },
                        "amount": {
                            "value": "-10000",
                            "currency": {
                                "symbol": "ICP",
                                "decimals": 8
                            }
                        }
                    }
                ],
                "metadata": {
                    "block_height": 6365500,
                    "memo": 37093,
                    "timestamp": 1687831749746816914
                }
            }
        }
    ],
    "total_count": 1
}
```


### 3.8. fetch support network list

Request example:
```curl
curl --location --request POST 'https://rosetta-api.internetcomputer.org/network/list' \
--header 'Content-Type: application/json' \
--data-raw '{
    "metadata": {}
}'
```

Return
```json
{
    "network_identifiers": [
        {
            "blockchain": "Internet Computer",
            "network": "00000000000000020101"
        }
    ]
}
```






