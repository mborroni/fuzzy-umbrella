package cache

import (
    "github.com/patrickmn/go-cache"

    "sync"
    "time"
)

type goCache struct {
    mutex         *sync.Mutex
    goCacheClient *cache.Cache
}

func NewGoCache() *goCache {
    cache_ := cache.New(1*time.Minute, 5*time.Minute)
    goCache := &goCache{
        goCacheClient: cache_,
        mutex:         &sync.Mutex{},
    }
    return goCache
}

func (cache *goCache) Get(key string) ([]byte, error) {
    cache.mutex.Lock()
    defer cache.mutex.Unlock()
    value, ok := cache.goCacheClient.Get(key)
    if !ok {
        return nil, nil
    }
    return value.([]byte), nil
}

func (cache *goCache) Save(key string, value []byte, timeExpiration ...int32) error {
    cache.mutex.Lock()
    defer cache.mutex.Unlock()
    te := 60 * time.Second
    if len(timeExpiration) > 0 {
        te = time.Duration(timeExpiration[0]) * time.Second
    }
    cache.goCacheClient.Set(key, value, te)
    return nil
}