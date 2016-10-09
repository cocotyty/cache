package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, exp time.Duration)
	Get(key string) (interface{}, bool)
	Clear(key string)
}
