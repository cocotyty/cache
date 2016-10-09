package cache

import (
	"github.com/dgryski/go-farm"
	"sync"
	"unsafe"
)

func NewCMap(size int) *CMap {
	cmap := &CMap{CMapEntries: make([]*CMapEntry, size), Size: uint32(size)}
	for ; size > 0; size-- {
		cmap.CMapEntries[size-1] = &CMapEntry{entry: map[string]*CacheValue{}}
	}
	return cmap
}

type CMap struct {
	CMapEntries []*CMapEntry
	Size        uint32
}
type CacheValue struct {
	Exp   int64
	Value interface{}
}
type CMapEntry struct {
	entry map[string]*CacheValue
	sync.RWMutex
}

func (c *CMapEntry) copyMap() map[string]*CacheValue {
	m := map[string]*CacheValue{}
	c.RLock()
	for k, v := range c.entry {
		m[k] = v
	}
	c.RUnlock()
	return m
}
func (c *CMap) Set(key string, value *CacheValue) {
	entry := c.CMapEntries[farm.Hash32(str2bytes(key))%c.Size]
	entry.Lock()
	entry.entry[key] = value
	entry.Unlock()
}

func (c *CMap) Del(key string) {
	p := int(farm.Hash32(str2bytes(key)) % c.Size)
	entry := c.CMapEntries[p]
	entry.Lock()
	delete(entry.entry, key)
	entry.Unlock()
}
func (c *CMap) Get(key string) (v *CacheValue, ok bool) {
	p := int(farm.Hash32(str2bytes(key)) % c.Size)
	entry := c.CMapEntries[p]
	entry.RLock()
	v, ok = entry.entry[key]
	entry.RUnlock()
	return
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
