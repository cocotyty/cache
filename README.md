# Memory Cache  In Go

It has not bad performance:

> BenchmarkRealTTLCache_Set        5000000	       232 ns/op
> BenchmarkPatrickmnGoCache        2000000	       713 ns/op

Easy To Use
```go
c:=cache.NewCache(32)
c.Set("FOO","BAR",100*time.Millisecond)
if value,ok:=c.Get("FOO");ok{
    fmt.Println(value)
}
c.Clear("FOO")
```
