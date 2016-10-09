package cache

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

var wg = sync.WaitGroup{}

func fuck() {
	tc := NewCache()
	tc.Set("sb", "sb", time.Second)
	time.Sleep(10 * time.Second)
	fmt.Println("finish")
	wg.Done()
}
func TestTTLCache_Set(t *testing.T) {
	wg.Add(1)
	go fuck()
	wg.Wait()
	runtime.GC()
	time.Sleep(3 * time.Second)
	fmt.Println("ok")
}
func BenchmarkNewTTLCache(b *testing.B) {
	ttlcache := NewCache()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ttlcache.Set(string(time.Now().UnixNano()), "sb", 10*time.Millisecond)
		}
	})
}
