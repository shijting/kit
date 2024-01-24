package kit

import (
	"math/rand"
	"time"
)

// RetryStrategy 重试策略
type RetryStrategy interface {
	// Next 返回下一次重试的间隔，如果不能继续重试，第二参数返回 false
	Next() (time.Duration, bool)
}

// FixIntervalRetry 固定间隔重试
type FixIntervalRetry struct {
	// 重试间隔
	Interval time.Duration
	// 最大次数
	MaxAttempts int
	// 当前次数
	CurrentAttempt int
}

func NewFixIntervalRetry(interval time.Duration, max int) *FixIntervalRetry {
	return &FixIntervalRetry{Interval: interval, MaxAttempts: max, CurrentAttempt: 0}
}

func (f FixIntervalRetry) Next() (time.Duration, bool) {
	f.CurrentAttempt++
	return f.Interval, f.CurrentAttempt <= f.MaxAttempts
}

// ExponentialBackoffRetry 指数退避重试
type ExponentialBackoffRetry struct {
	// 初始间隔
	InitialInterval time.Duration
	// 最大间隔
	MaxInterval time.Duration
	// 最大次数
	MaxAttempts int
	// 当前次数，从 0 开始，表示第一次重试，第二次重试，第三次重试，以此类推。
	CurrentAttempt int
}

func NewExponentialBackoffRetry(initialInterval, maxInterval time.Duration, maxAttempts int) *ExponentialBackoffRetry {
	return &ExponentialBackoffRetry{
		InitialInterval: initialInterval,
		MaxInterval:     maxInterval,
		MaxAttempts:     maxAttempts,
		CurrentAttempt:  0,
	}
}

func (e *ExponentialBackoffRetry) Next() (time.Duration, bool) {
	if e.CurrentAttempt >= e.MaxAttempts {
		return 0, false
	}

	interval := e.InitialInterval * time.Duration(1<<uint(e.CurrentAttempt))
	if interval > e.MaxInterval {
		interval = e.MaxInterval
	}

	e.CurrentAttempt++
	return interval, true
}

// RandomizedRetry 随机间隔重试
type RandomizedRetry struct {
	MaxInterval    time.Duration
	MaxAttempts    int
	CurrentAttempt int
}

func NewRandomizedRetry(maxInterval time.Duration, maxAttempts int) *RandomizedRetry {
	return &RandomizedRetry{
		MaxInterval:    maxInterval,
		MaxAttempts:    maxAttempts,
		CurrentAttempt: 0,
	}
}

func (r *RandomizedRetry) Next() (time.Duration, bool) {
	if r.CurrentAttempt >= r.MaxAttempts {
		return 0, false
	}

	// 生成一个随机的等待时间
	randomInterval := time.Duration(rand.Int63n(int64(r.MaxInterval)))

	r.CurrentAttempt++
	return randomInterval, true
}
