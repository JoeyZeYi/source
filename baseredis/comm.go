package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type ICommStringCache[T any] interface {
	Get(ctx context.Context, key string) (*T, error)
	Set(ctx context.Context, key string, t *T, expiration time.Duration) error
	GetSet(ctx context.Context, key string, value *T) (*T, error)
	MGet(ctx context.Context, keys ...string) ([]*T, error)
	MSet(ctx context.Context, maps map[string]*T) error
	Incr(ctx context.Context, key string) (int64, error)              //自增1
	IncrBy(ctx context.Context, key string, val int64) (int64, error) //自增N
	Decr(ctx context.Context, key string) (int64, error)              //自减1
	DecrBy(ctx context.Context, key string, val int64) (int64, error) //自减N
}

type ICommHashCache[T any] interface {
	HSet(ctx context.Context, key, field string, t *T) error
	HLen(ctx context.Context, key string) (int64, error)
	HDel(ctx context.Context, key string, field ...string) error
	HGet(ctx context.Context, key, field string) (*T, error)
	HGetAll(ctx context.Context, key string) (map[string]*T, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]*T, error)
	HMSet(ctx context.Context, key string, maps map[string]*T) error
}

type ICommCache[T any] interface {
	Expire(ctx context.Context, key string, expiration time.Duration) error //设置过期时间
	TTL(ctx context.Context, key string) (time.Duration, error)             //剩余过期时间
	Exists(ctx context.Context, key string) (bool, error)                   //key是否存在
	Del(ctx context.Context, keys ...string) error
	String() ICommStringCache[T]
	Hash() ICommHashCache[T]
	Geo() ICommGeoCache[T]
	GetClient() *redis.Client
}

type ICommGeoCache[T any] interface {
	GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) error
	// GeoQuery 获取成员的经纬度信息
	GeoQuery(ctx context.Context, key string, members ...string) (map[string]*redis.GeoPos, error)
	// GeoDist 获取两个成员之间的距离 unit:[m|km|ft|mi] 米|千米|英里|英尺
	GeoDist(ctx context.Context, key, member1, member2, unit string) (float64, error)
	// GeoRadius 查询输入经纬度周边的用户信息 Radius半径 WithDist是否输出距离
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error)
	// GeoDel 删除成员的经纬度信息
	GeoDel(ctx context.Context, key, member string) error
}

type CommCache[T any] struct {
	client *redis.Client
}

func NewCommCache[T any](client *redis.Client) *CommCache[T] {
	c := &CommCache[T]{client: client}
	return c
}

func (i *CommCache[T]) Geo() ICommGeoCache[T] {
	return i
}

func (i *CommCache[T]) String() ICommStringCache[T] {
	return i
}
func (i *CommCache[T]) Hash() ICommHashCache[T] {
	return i
}

func (i *CommCache[T]) GeoDel(ctx context.Context, key, member string) error {
	return i.client.ZRem(ctx, key, member).Err()
}

func (i *CommCache[T]) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return i.client.GeoRadius(ctx, key, longitude, latitude, query).Result()
}

func (i *CommCache[T]) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) error {
	return i.client.GeoAdd(ctx, key, geoLocation...).Err()
}

func (i *CommCache[T]) GeoDist(ctx context.Context, key, member1, member2, unit string) (float64, error) {
	return i.client.GeoDist(ctx, key, member1, member2, unit).Result()
}

func (i *CommCache[T]) GeoQuery(ctx context.Context, key string, members ...string) (map[string]*redis.GeoPos, error) {
	result, err := i.client.GeoPos(ctx, key, members...).Result()
	if err != nil {
		return nil, err
	}
	if len(result) != len(members) {
		return nil, errors.New("数据查询异常")
	}
	geoPosMaps := make(map[string]*redis.GeoPos)
	for index, k := range members {
		if result[index] == nil {
			continue
		}
		geoPosMaps[k] = result[index]
	}
	return geoPosMaps, nil
}

func (i *CommCache[T]) TTL(ctx context.Context, key string) (time.Duration, error) {
	return i.client.TTL(ctx, key).Result()
}

func (i *CommCache[T]) Exists(ctx context.Context, key string) (bool, error) {
	result, err := i.client.Exists(ctx, key).Result()
	return result == 1, err
}

func (i *CommCache[T]) Get(ctx context.Context, key string) (*T, error) {
	t := new(T)
	bytes, err := i.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (i *CommCache[T]) Set(ctx context.Context, key string, t *T, expiration time.Duration) error {

	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return i.client.Set(ctx, key, bytes, expiration).Err()
}
func (i *CommCache[T]) Del(ctx context.Context, key ...string) error {
	return i.client.Del(ctx, key...).Err()
}

func (i *CommCache[T]) GetClient() *redis.Client {
	return i.client
}

func (i *CommCache[T]) Incr(ctx context.Context, key string) (int64, error) {
	return i.client.Incr(ctx, key).Result()
}

func (i *CommCache[T]) IncrBy(ctx context.Context, key string, val int64) (int64, error) {
	return i.client.IncrBy(ctx, key, val).Result()
}

func (i *CommCache[T]) Decr(ctx context.Context, key string) (int64, error) {
	return i.client.Decr(ctx, key).Result()
}

func (i *CommCache[T]) DecrBy(ctx context.Context, key string, val int64) (int64, error) {
	return i.client.DecrBy(ctx, key, val).Result()
}

func (i *CommCache[T]) GetSet(ctx context.Context, key string, value *T) (*T, error) {
	reqBytes, _ := json.Marshal(value)

	resBytes, err := i.client.GetSet(ctx, key, reqBytes).Bytes()
	if err != nil {
		return nil, err
	}
	result := new(T)
	if err = json.Unmarshal(resBytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *CommCache[T]) MGet(ctx context.Context, keys ...string) ([]*T, error) {
	if len(keys) == 0 {
		return nil, errors.New("keys len is 0")
	}

	results, err := i.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	list := make([]*T, 0)
	for _, result := range results {
		if s, ok := result.(string); ok {
			t := new(T)
			if err = json.Unmarshal([]byte(s), t); err != nil {
				return nil, err
			}

			list = append(list, t)
		}

	}

	return list, nil
}

func (i *CommCache[T]) MSet(ctx context.Context, maps map[string]*T) error {
	values := make([]interface{}, 0)
	for k, val := range maps {
		bytes, _ := json.Marshal(val)
		values = append(values, k, bytes)
	}

	return i.client.MSet(ctx, values).Err()
}

func (i *CommCache[T]) HSet(ctx context.Context, key, field string, t *T) error {
	bytes, _ := json.Marshal(t)
	return i.client.HSet(ctx, key, field, bytes).Err()
}

func (i *CommCache[T]) HLen(ctx context.Context, key string) (int64, error) {
	return i.client.HLen(ctx, key).Result()
}
func (i *CommCache[T]) HDel(ctx context.Context, key string, field ...string) error {
	return i.client.HDel(ctx, key, field...).Err()
}

func (i *CommCache[T]) HGet(ctx context.Context, key, field string) (*T, error) {
	bytes, err := i.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		return nil, err
	}
	t := new(T)
	if err = json.Unmarshal(bytes, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (i *CommCache[T]) HGetAll(ctx context.Context, key string) (map[string]*T, error) {
	results, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	maps := make(map[string]*T)
	for k, v := range results {

		t := new(T)
		if err = json.Unmarshal([]byte(v), t); err != nil {
			return nil, err
		}

		maps[k] = t
	}

	return maps, nil
}

func (i *CommCache[T]) HMGet(ctx context.Context, key string, fields ...string) ([]*T, error) {
	results, err := i.client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	list := make([]*T, 0, len(results))
	for _, v := range results {

		if s, ok := v.(string); ok {
			t := new(T)
			if err = json.Unmarshal([]byte(s), t); err != nil {
				return nil, err
			}

			list = append(list, t)
		}

	}

	return list, nil
}

func (i *CommCache[T]) HMSet(ctx context.Context, key string, maps map[string]*T) error {
	values := make(map[string]interface{})
	for k, v := range maps {

		bytes, _ := json.Marshal(v)

		values[k] = bytes
	}

	return i.client.HMSet(ctx, key, values).Err()
}

func (i *CommCache[T]) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return i.client.Expire(ctx, key, expiration).Err()
}
