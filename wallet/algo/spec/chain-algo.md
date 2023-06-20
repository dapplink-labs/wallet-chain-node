
### 1. document and sdk

go-sdk: https://github.com/algorand/go-algorand-sdk
open node：https://mainnet-api.algonode.cloud， 
node operator：https://algonode.io/api/#free-as-in--algorand-api-access
docs for open node: https://developer.algorand.org/docs/get-started/devenv/
explorer：https://algoexplorer.io/
api docs：https://developer.algorand.org/docs/rest-apis/algod/v2/


### 2. api
2.1. get account balance
```curl
curl https://mainnet-api.algonode.cloud/v2/accounts/ZW3ISEHZUHPO7OZGMKLKIIMKVICOUDRCERI454I3DB2BH52HGLSO67W754/
```

```python
def get_account(address):
   path = "/v2/accounts/%s"% address
   info = await self.request_rest(path, method='GET')
   return info['balance']
```


2.2.get transaction sign params

```curl
curl https://mainnet-api.algonode.cloud/v2/transactions/params
```
```python  
 def get_transaction_params():
   path = '/v1/transactions/params'
   return await self.request_rest(path, method='GET')
```

2.3. send raw transaction
please see /v2/transactions for detail

```python 
def send_transaction(tx):
   path = '/v1/transactions'
   return await self.request_rest(path, method='POST', rawtxn=bytes.fromhex(tx))
```python 


2.4. 根据地址获取交易列表
https://algoexplorer.io/api-dev/indexer-v2


2.5.根据 Hash 获取交易详情

https://algoexplorer.io/api-dev/indexer-v2

