package rlock

import "time"

// RetryStrategy 重试策略
type RetryStrategy interface {
	// Next 返回下一次重试的间隔，如果不需要继续重试，那么第二参数返回 false
	Next() (time.Duration, bool)
}

// FixIntervalRetry 固定间隔重试
type FixIntervalRetry struct {
	// 重试间隔
	Interval time.Duration
	// 最大次数
	Max int
	// 当前次数
	cnt int
}

func NewFixIntervalRetry(interval time.Duration, max int, cnt int) *FixIntervalRetry {
	return &FixIntervalRetry{Interval: interval, Max: max, cnt: cnt}
}

func (f FixIntervalRetry) Next() (time.Duration, bool) {
	f.cnt++
	return f.Interval, f.cnt <= f.Max
}
