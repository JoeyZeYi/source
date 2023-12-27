package cache

import "github.com/go-redis/redis/v8"

type DemoComm struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type IDemoCommCache interface {
	ICommCache[DemoComm]
	ICommStringCache[DemoComm]
}

type demoCommCache struct {
	ICommCache[DemoComm]
	ICommStringCache[DemoComm]
}

func NewDemoCommCache(client *redis.Client) IDemoCommCache {
	commCache := NewCommCache[DemoComm](client)
	return &demoCommCache{
		ICommCache:       NewCommCache[DemoComm](client),
		ICommStringCache: commCache.String(),
	}
}

func (d *demoCommCache) GetTest() {
}
