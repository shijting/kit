package limiter

import (
	"sync"
	"time"
)

// Bucket 令牌桶限流器
type Bucket struct {
	// cap 桶容量
	cap int64
	// tokens 桶中的令牌数
	tokens int64
	// rate 令牌产生速率/s
	rate int64
	// 上次执行rate计算的时间
	last int64
	mu   sync.Mutex
}

func NewBucket(cap, rate int64) *Bucket {
	if cap <= 0 || rate <= 0 {
		panic("cap and rate must be greater than 0")
	}
	return &Bucket{cap: cap, tokens: cap, rate: rate}
}

// add 添加令牌
//func (b *Bucket) add() {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//
//	if b.tokens+b.rate < b.cap {
//		b.tokens += b.rate
//	} else {
//		b.tokens = b.cap
//	}
//}

// Accept 是否允许通过
func (b *Bucket) Accept() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now().Unix()
	// 计算当前应该添加的令牌数
	b.tokens += (now - b.last) * b.rate
	if b.tokens > b.cap {
		b.tokens = b.cap
	}
	b.last = now

	if b.tokens <= 0 {
		return false
	}

	b.tokens--
	return true
}
