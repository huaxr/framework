// Author: huaxinrui@tal.com
// Time: 2022-12-15 11:27
// Git: huaxr

package apicache

import (
	"sync"

	"github.com/huaxr/framework/component/kv"
	"github.com/huaxr/framework/pkg/confutil"
)

type apiCache struct {
	sync.RWMutex
	// key: uri, val: response, expire: sec
	apiCacheSet kv.Cache
	// uri:sec
	keyExpireMap map[string]int
}

var cache = &apiCache{
	keyExpireMap: make(map[string]int),
	apiCacheSet:  kv.InitExpireCache(),
}

func GetApiCache() *apiCache {
	return cache
}

func (a *apiCache) GetCache() kv.Cache {
	return a.apiCacheSet
}

func (a *apiCache) GetExpSec(path string) (int, bool) {
	a.RLock()
	defer a.RUnlock()
	v, ok := a.keyExpireMap[path]
	return v, ok
}

func (a *apiCache) UpdateFromTcm(configs *confutil.DynamicConfig) {
	a.Lock()
	defer a.Unlock()

	var m = make(map[string]int)
	for _, cache := range configs.Caches {
		m[cache.Path] = cache.Duration
	}

	a.keyExpireMap = m
}

func (a *apiCache) Size() int {
	return a.apiCacheSet.KVSize()
}
