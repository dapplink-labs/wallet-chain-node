package cache

import lru "github.com/hashicorp/golang-lru"

var txCache *lru.ARCCache
var balanceCache *lru.ARCCache

func init() {
	txCache, _ = lru.NewARC(1000)
	balanceCache, _ = lru.NewARC(1000)
}

func GetTxCache() *lru.ARCCache {
	return txCache
}

func GetBalanceCache() *lru.ARCCache {
	return balanceCache
}
