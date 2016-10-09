package cache

import (
	"log"
	"runtime"
	"time"
)

// this is for gc stop clear goroutine trick
type TTLCache struct {
	*ttlCache
}
type ttlCache struct {
	clearDuration time.Duration
	stopGc        chan struct{}
	CMap
}

func NewCache(shared int) *TTLCache {
	cmap := NewCMap(shared)
	tc := &ttlCache{CMap: *cmap, clearDuration: 100 * time.Millisecond, stopGc: make(chan struct{})}
	rtc := &TTLCache{tc}
	go tc.loopClear()
	runtime.SetFinalizer(rtc, clearTTL)
	return rtc
}

func clearTTL(tc *TTLCache) {
	tc.stopGc <- struct{}{}
}


func (c *ttlCache) Set(key string, value interface{}, exp time.Duration) {
	c.CMap.Set(key, &CacheValue{Value: value, Exp: time.Now().Add(exp).UnixNano()})
}
func (c *ttlCache) Get(key string) (interface{}, bool) {
	cv, ok := c.CMap.Get(key)
	if !ok {
		return nil, false
	}
	if cv.Exp > time.Now().UnixNano() {
		c.CMap.Del(key)
		return nil, false
	}
	return cv.Value, true
}

func (c *ttlCache) Clear(key string) {
	c.CMap.Del(key)
}

func (c *ttlCache) loopClear() {
	ticker := time.NewTicker(c.clearDuration)
	for {
		select {
		case <-ticker.C:
			log.Println("ticker")
			c.clearExpire()
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}
func (c *ttlCache) clearExpire() {
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
