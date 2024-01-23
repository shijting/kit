package limiter

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shijting/kit/option"
	"log"
	"net/http"
	"time"
)

type IpLimiter struct {
	prefix   string
	cli      redis.Cmdable
	interval time.Duration
	rate     int
}

//go:embed script/slide_window.lua
var luaScript string

func NewIpLimiter(cli redis.Cmdable, opts ...option.Option[IpLimiter]) *IpLimiter {
	limiter := &IpLimiter{
		cli:      cli,
		prefix:   "ip-limiter",
		interval: time.Second,
		rate:     200,
	}

	option.Options[IpLimiter](opts).Apply(limiter)

	return limiter
}

func WithInterval(interval time.Duration) option.Option[IpLimiter] {
	return func(t *IpLimiter) {
		t.interval = interval
	}
}

func WithRate(rate int) option.Option[IpLimiter] {
	return func(t *IpLimiter) {
		t.rate = rate
	}
}

func WithPrefix(prefix string) option.Option[IpLimiter] {
	return func(t *IpLimiter) {
		t.prefix = prefix
	}
}

func (b *IpLimiter) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *IpLimiter) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cli.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
