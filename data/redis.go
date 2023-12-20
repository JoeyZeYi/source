package data

import (
	"context"
	"crypto/tls"
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

type RedisOption struct {
	Network            string
	Addr               string
	Dialer             func(ctx context.Context, network, addr string) (net.Conn, error)
	OnConnect          func(ctx context.Context, cn *redis.Conn) error
	Username           string
	Password           string
	DB                 int
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolFIFO           bool
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	readOnly           bool
	TLSConfig          *tls.Config
	Limiter            redis.Limiter
}

type IRedisOption interface {
	apply(*RedisOption)
}

type redisOptionFunc func(*RedisOption)

func (r redisOptionFunc) apply(option *RedisOption) {
	r(option)
}

// RedisNetwork 网络类型，一般为 "tcp"
func RedisNetwork(network string) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.Network = network
	})
}

// RedisAddr Redis 服务器地址，格式为 "host:port"。
func RedisAddr(addr string) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.Addr = addr
	})
}

func RedisDialer(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.Dialer = dialer
	})
}

func RedisMaxRetries(maxRetries int) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.MaxRetries = maxRetries
	})
}

func RedisMinRetryBackoff(minRetryBackoff time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.MinRetryBackoff = minRetryBackoff
	})
}
func RedisMaxRetryBackoff(maxRetryBackoff time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.MaxRetryBackoff = maxRetryBackoff
	})
}

func RedisDialTimeout(dialTimeout time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.DialTimeout = dialTimeout
	})
}
func RedisReadTimeout(readTimeout time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.ReadTimeout = readTimeout
	})
}
func RedisWriteTimeout(writeTimeout time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.WriteTimeout = writeTimeout
	})
}

func RedisPoolFIFO(poolFIFO bool) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.PoolFIFO = poolFIFO
	})
}

func RedisPoolSize(poolSize int) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.PoolSize = poolSize
	})
}

func RedisMinIdleConns(minIdleConns int) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.MinIdleConns = minIdleConns
	})
}

func RedisMaxConnAge(maxConnAge time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.MaxConnAge = maxConnAge
	})
}
func RedisPoolTimeout(poolTimeout time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.PoolTimeout = poolTimeout
	})
}
func RedisIdleTimeout(idleTimeout time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.IdleTimeout = idleTimeout
	})
}
func RedisIdleCheckFrequency(idleCheckFrequency time.Duration) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.IdleCheckFrequency = idleCheckFrequency
	})
}

func RedisTLSConfig(tLSConfig *tls.Config) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.TLSConfig = tLSConfig
	})
}

func RedisLimiter(Limiter redis.Limiter) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.Limiter = Limiter
	})
}

// RedisUserPwdDB 用户名+密码，如果不需要密码则为空字符串。db
func RedisUserPwdDB(userName, password string, db int) IRedisOption {
	return redisOptionFunc(func(option *RedisOption) {
		option.Username = userName
		option.Password = password
		option.DB = db
	})
}

func NewRedis(options ...IRedisOption) (*redis.Client, error) {
	option := &RedisOption{}

	for _, redisOption := range options {
		redisOption.apply(option)
	}

	rdb := redis.NewClient(&redis.Options{
		Network:            option.Network,
		Addr:               option.Addr,
		Dialer:             option.Dialer,
		OnConnect:          option.OnConnect,
		Username:           option.Username,
		Password:           option.Password,
		DB:                 option.DB,
		MaxRetries:         option.MaxRetries,
		MinRetryBackoff:    option.MinRetryBackoff,
		MaxRetryBackoff:    option.MaxRetryBackoff,
		DialTimeout:        option.DialTimeout,
		ReadTimeout:        option.ReadTimeout,
		WriteTimeout:       option.WriteTimeout,
		PoolFIFO:           option.PoolFIFO,
		PoolSize:           option.PoolSize,
		MinIdleConns:       option.MinIdleConns,
		MaxConnAge:         option.MaxConnAge,
		PoolTimeout:        option.PoolTimeout,
		IdleTimeout:        option.IdleTimeout,
		IdleCheckFrequency: option.IdleCheckFrequency,
		TLSConfig:          option.TLSConfig,
		Limiter:            option.Limiter,
	})
	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		log.Error("NewRedis", zap.Error(err), zap.Any("option", option))
	}
	return rdb, err
}
