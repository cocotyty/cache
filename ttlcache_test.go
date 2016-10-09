package cache

import (
	"testing"
	"time"
	"strconv"
	"github.com/patrickmn/go-cache"
)

func BenchmarkRealTTLCache_Set(b *testing.B) {
	ttlcache := NewCache(256)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ttlcache.Set(strconv.Itoa(int(time.Now().UnixNano())), "sometext", 10*time.Millisecond)
		}
	})
}
func BenchmarkPatrickmnGoCache(b *testing.B) {
	ttlcache := cache.New(time.Second,100*time.Millisecond)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ttlcache.Set(strconv.Itoa(int(time.Now().UnixNano())), "sometext", 10*time.Millisecond)
		}
	})
}
func TestNewCache(t *testing.T) {
	c:=NewCache(32)
	c.Set("ABC","1",100*time.Millisecond)
	if v,ok:=c.Get("ABC");ok{
		if v != "1"{
			t.Fatal("not equal!")
		}
	}else{
		t.Fatal("cache miss")
	}
	time.Sleep(100*time.Millisecond)
	if _,ok:=c.Get("ABC");ok{
		t.Fatal("cache should miss")
	}
	t.Log("success")
}