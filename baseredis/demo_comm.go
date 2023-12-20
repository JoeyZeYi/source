package cache

import "github.com/go-redis/redis/v8"

type DemoComm struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type IDemoCommCache interface {
	ICommCache[DemoComm]
	GetTest() //额外接口
}

type demoCommCache struct {
	ICommCache[DemoComm]
}

func NewDemoCommCache(client *redis.Client) IDemoCommCache {
	return &demoCommCache{ICommCache: NewCommCache[DemoComm](client)}
}

func (d *demoCommCache) GetTest() {
}
