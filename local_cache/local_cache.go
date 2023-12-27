package local_cache

import (
	"context"
	"github.com/JoeyZeYi/source/util"
	"time"
)

// IRedisCache 实现该接口，可同步到redis接口操作
type IRedisCache[V any] interface {
	Get(ctx context.Context, key string) (*V, error)
	Set(ctx context.Context, key string, v *V, expire time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type Data[V any] struct {
	expireTime int64
	v          *V
}

type LocalCache[V any] struct {
	m             util.Map[string, Data[V]]
	timer         time.Duration //定时处理过期key
	defaultExpire time.Duration //默认过期时间
	redisCache    IRedisCache[V]
}

func NewLocalCache[V any](timer, defaultExpire time.Duration, redisCache IRedisCache[V]) *LocalCache[V] {
	l := &LocalCache[V]{
		m:             util.Map[string, Data[V]]{},
		timer:         timer,
		defaultExpire: defaultExpire,
		redisCache:    redisCache,
	}
	go l.expireKeyTimer()
	return l
}

func (l *LocalCache[V]) Set(key string, v *V, duration time.Duration) {
	d := &Data[V]{
		expireTime: time.Now().Add(duration).Unix(),
		v:          v,
	}
	l.m.Set(key, d)
	if l.redisCache != nil {
		_ = l.redisCache.Set(context.TODO(), key, v, duration)
	}
}

func (l *LocalCache[V]) Get(key string) *V {
	data := l.m.Get(key)
	if data != nil {
		if data.expireTime > time.Now().Unix() {
			return data.v
		}
	}
	if l.redisCache != nil {
		val, err := l.redisCache.Get(context.TODO(), key)
		if err != nil {
			return nil
		}
		data = &Data[V]{
			expireTime: time.Now().Add(l.defaultExpire).Unix(),
			v:          val,
		}
		l.m.Set(key, data)
		return val
	}

	return nil
}

func (l *LocalCache[V]) Delete(key string) {
	l.m.Delete(key)
	if l.redisCache != nil {
		_ = l.redisCache.Del(context.TODO(), key)
	}
}

func (l *LocalCache[V]) expireKeyTimer() {
	for {
		time.Sleep(l.timer)
		for k, data := range l.m.Maps() {
			if data.expireTime < time.Now().Unix() {
				l.m.Delete(k)
			}
		}
	}
}
