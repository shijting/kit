package lockx

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shijting/kit/option"
	"time"
)

var (
	//go:embed script/lua/lock.lua
	luaLock string
	//go:embed script/lua/unlock.lua
	luaUnlock string
	//go:embed script/lua/refresh.lua
	luaRefresh string
	// ErrReleaseLock 释放锁错误
	ErrReleaseLock = errors.New("release lock error")
	// ErrNotHoldingLock 未持有锁
	ErrNotHoldingLock = errors.New("not holding lock")
	// ErrGetLockFailed 获取锁失败
	ErrGetLockFailed = errors.New("get lock failed")
	// ErrorNotGetLock 未抢到锁错误
	ErrorNotGetLock = errors.New("not get lock")
)

type Client struct {
	client redis.Cmdable
	// value setnx 的值
	value string
	// expiration 锁的过期时间
	expiration time.Duration
	// timeout 调用redis的超时时间
	timeout time.Duration
	retry   RetryStrategy
}

func NewClient(client redis.Cmdable, opts ...option.Option[Client]) *Client {
	cli := &Client{
		client:     client,
		value:      uuid.New().String(),
		expiration: 30 * time.Second,
		timeout:    3 * time.Second,
	}
	option.Options[Client](opts).Apply(cli)
	return cli
}

// WithValue 设置 setnx 的值
func WithValue(value string) option.Option[Client] {
	return func(t *Client) {
		t.value = value
	}
}

// WithExpiration 设置锁的过期时间
func WithExpiration(expiration time.Duration) option.Option[Client] {
	return func(t *Client) {
		t.expiration = expiration
	}
}

// WithTimeout 设置锁的超时时间
func WithTimeout(timeout time.Duration) option.Option[Client] {
	return func(t *Client) {
		t.timeout = timeout
	}
}

// WithRetry 设置重试策略
func WithRetry(retry RetryStrategy) option.Option[Client] {
	return func(t *Client) {
		t.retry = retry
	}
}

func (cli *Client) Lock(ctx context.Context, key string) (*Lock, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	for {
		rct, cancel := context.WithTimeout(ctx, cli.timeout)
		res, err := cli.client.Eval(rct, luaLock, []string{key}, cli.value, cli.expiration.Seconds()).Result()
		cancel()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		if res == "OK" {
			return newLock(cli.client, key, cli.value, cli.expiration, cli.timeout), nil
		}

		if cli.retry != nil {
			nextInterval, canRetry := cli.retry.Next()
			if !canRetry {
				if err != nil {
					err = fmt.Errorf("最后一次重试错误: %w", err)
				} else {
					err = fmt.Errorf("未抢到锁: %w", ErrorNotGetLock)
				}
				return nil, fmt.Errorf("重试机会耗尽，%w", err)
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(nextInterval):
				continue
			}
		}
		return nil, ErrGetLockFailed
	}
}

// TryLock 尝试获取锁, 不一定能获取到
func (cli *Client) TryLock(ctx context.Context, key string, expire, timeout time.Duration) (*Lock, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	rct, cancel := context.WithTimeout(ctx, cli.timeout)
	ok, err := cli.client.SetNX(rct, key, cli.value, cli.expiration).Result()
	cancel()
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrGetLockFailed
	}
	return newLock(cli.client, key, cli.value, expire, timeout), nil

}

type Lock struct {
	client     redis.Cmdable
	key        string
	value      string
	expiration time.Duration
	timeout    time.Duration
}

func newLock(client redis.Cmdable, key, value string, expiration, timeout time.Duration) *Lock {
	return &Lock{
		client:     client,
		key:        key,
		value:      value,
		expiration: expiration,
		timeout:    timeout,
	}
}

// Unlock 释放锁
func (l *Lock) Unlock(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()

	if err == redis.Nil || res != 1 {
		return ErrNotHoldingLock
	}
	if err != nil {
		// TODO: 其他未知错误，要不要重试？
		return err
	}
	return nil
}

// Refresh 刷新锁的过期时间
func (l *Lock) Refresh(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.value, l.expiration.Seconds()).Int64()

	if err == redis.Nil || res == 0 {
		return ErrNotHoldingLock
	}
	if err != nil {
		return err
	}
	return nil
}

// AutoRefresh 自动续约
func (l *Lock) AutoRefresh(interval time.Duration, stop <-chan struct{}, maxRetry int) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		//ticker := time.NewTicker(l.expiration / 2)
		ticker := time.NewTicker(interval)
		defer func() {
			ticker.Stop()
			close(errCh)
		}()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
				err := l.Refresh(ctx)
				cancel()
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						//	超时进行重试
						if maxRetry > 0 {
							for i := maxRetry; i > 0; i-- {
								ctx, cancel = context.WithTimeout(context.Background(), l.timeout)
								err = l.Refresh(ctx)
								cancel()
								// 重试成功或者没必要重试的错误，退出重试
								if err == nil || errors.Is(err, ErrNotHoldingLock) || err == redis.Nil {
									break
								}
							}
						}
					}
					// 重试后还是失败，返回错误并退出
					if err != nil {
						errCh <- err
						return
					}
				}
			}
		}
	}()
	return errCh
}
