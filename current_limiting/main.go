package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	rate       float64   // 令牌生成速率（每秒生成的令牌数）
	capacity   float64   // 令牌桶容量
	tokens     float64   // 当前令牌数量
	lastAccess time.Time // 上次访问时间
	mu         sync.Mutex
}

func NewTokenBucket(rate, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:     rate,
		capacity: capacity,
		tokens:   capacity,
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastAccess).Seconds()
	tb.tokens += elapsed * tb.rate
	tb.lastAccess = now

	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

func main() {
	// 创建一个每秒生成5个令牌，令牌桶容量为10的令牌桶
	tokenBucket := NewTokenBucket(5, 10)

	// 模拟一些请求
	for i := 0; i < 15; i++ {
		if tokenBucket.Allow() {
			fmt.Printf("Request %d allowed\n", i+1)
		} else {
			fmt.Printf("Request %d denied\n", i+1)
		}

		time.Sleep(1 * time.Millisecond) // 模拟请求之间的间隔
	}
}
