package global

import (
	"sync"
	"time"
)

/*
全局本地缓存结构体,作为redis的备用场景
*/
type Cache struct {
	store sync.Map
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	if _, exist := c.Get(key); exist {
		c.store.Delete(key)
	}
	c.store.Store(key, value)
	if duration > 0 {
		go func() {
			<-time.After(duration)
			c.store.Delete(key)
		}()
	}
}
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.store.Load(key)
}

var Store = new(Cache)
