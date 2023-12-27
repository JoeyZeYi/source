package local_cache

import (
	cache "github.com/JoeyZeYi/source/baseredis"
	"github.com/go-redis/redis/v8"
	"reflect"
	"testing"
	"time"
)

type UserInfo struct {
}

func TestNewLocalCache(t *testing.T) {
	type args[V any] struct {
		timer         time.Duration
		defaultExpire time.Duration
		redisCache    IRedisCache[V]
	}
	type testCase[V any] struct {
		name string
		args args[V]
		want *LocalCache[V]
	}
	//已经初始化
	var redisClient *redis.Client
	userRedisCache := cache.NewCommCache[UserInfo](redisClient).String()
	tests := []testCase[UserInfo]{
		{
			name: "",
			args: args[UserInfo]{
				timer:         time.Second,
				defaultExpire: time.Minute * 30,
				redisCache:    userRedisCache,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocalCache(tt.args.timer, tt.args.defaultExpire, tt.args.redisCache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocalCache() = %v, want %v", got, tt.want)
				got.Set("Test", &UserInfo{}, time.Minute)
			}
		})
	}
}
