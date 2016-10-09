package cache

import (
	"runtime"
	"time"
)

// when it  finalized ,program will automatic  stop realTTLCache.loopClear goroutine,because
// this trick is  learning from https://github.com/patrickmn/go-cache
type ttlCache struct {
	*realTTLCache
}

// real cache struct
type realTTLCache struct {
	clearDuration time.Duration
	stopGc        chan struct{}
	CMap
}

func NewCache(shared int) Cache {
	concurrentMap := NewCMap(shared)
	realCache := &realTTLCache{CMap: *concurrentMap, clearDuration: 100 * time.Millisecond, stopGc: make(chan struct{})}
	go realCache.loopClear()
	cache := &ttlCache{realCache}
	runtime.SetFinalizer(cache, clearTTL)
	return cache
}

func clearTTL(tc *ttlCache) {
	tc.stopGc <- struct{}{}
}

func (c *realTTLCache) Set(key string, value interface{}, exp time.Duration) {
	c.CMap.Set(key, &CacheValue{Value: value, Exp: time.Now().Add(exp).UnixNano()})
}
func (c *realTTLCache) Get(key string) (interface{}, bool) {
	cv, ok := c.CMap.Get(key)
	if !ok {
		return nil, false
	}
	if cv.Exp < time.Now().UnixNano() {
		c.CMap.Del(key)
		return nil, false
	}
	return cv.Value, true
}

func (c *realTTLCache) Clear(key string) {
	c.CMap.Del(key)
}

func (c *realTTLCache) loopClear() {
	ticker := time.NewTicker(c.clearDuration)
	for {
		select {
		case <-ticker.C:
			c.clearExpire()
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}
func (c *realTTLCache) clearExpire() {
	for _, en := range c.CMapEntries {
		now := time.Now().UnixNano()
		en.Lock()
		for k, v := range en.entry {
			if v.Exp < now {
				delete(en.entry, k)
			}
		}
		en.Unlock()
	}
}
