package cache

import (
	"testing"
	"time"
	"github.com/streamrail/concurrent-map"
	"github.com/dgryski/go-farm"
	"hash/fnv"
)

func BenchmarkCMap_Set(b *testing.B) {
	cmap:=NewCMap(32)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			cmap.Set(time.Now().String(),nil)
		}
	})
}
func BenchmarkCt(b *testing.B) {
	cmap:=cmap.New()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			cmap.Set(time.Now().String(),"ww")
		}
	})
}
var buf = []byte("RMVx)@MLxH9M.WeGW-ktWwR3Cy1XS.,K~i@n-Y+!!yx4?AB%cM~l/#0=2:BOn7HPipG&o/6Qe<hU;$w1-~bU4Q7N&yk/8*Zz.Yg?zl9bVH/pXs6Bq^VdW#Z)NH!GcnH-UesRd@gDij?luVQ3;YHaQ<~SBm17G9;RWvGlsV7tpe*RCe=,?$nE1u9zvjd+rBMu7_Rg4)2AeWs^aaBr&FkC#rcwQ.L->I+Da7Qt~!C^cB2wq(^FGyB?kGQpd(G8I.A7")
var res32 uint32

func BenchmarkHash32(b *testing.B) {
	var r uint32
	for i := 0; i < b.N; i++ {
		// record the result to prevent the compiler eliminating the function call
		r = farm.Hash32(buf)
	}
	// store the result to a package level variable so the compiler cannot eliminate the Benchmark itself
	res32 = r
}
func BenchmarkHashByFvn(b *testing.B) {
	var r uint32
	for i := 0; i < b.N; i++ {
		hasher:=fnv.New32()
		hasher.Write(buf)
		r=hasher.Sum32()
	}
	res32=r
}