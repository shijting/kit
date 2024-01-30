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

// SlideWindowIPLimiter 滑动窗口IP 限流器
type SlideWindowIPLimiter struct {
	prefix   string
	cli      redis.Cmdable
	interval time.Duration
	rate     int
}

//go:embed script/slide_window.lua
var luaScript string

func NewSlideWindowIPLimiter(cli redis.Cmdable, opts ...option.Option[SlideWindowIPLimiter]) *SlideWindowIPLimiter {
	limiter := &SlideWindowIPLimiter{
		cli:      cli,
		prefix:   "ip-limiter",
		interval: time.Second,
		rate:     200,
	}

	option.Options[SlideWindowIPLimiter](opts).Apply(limiter)

	return limiter
}

func WithInterval(interval time.Duration) option.Option[SlideWindowIPLimiter] {
	return func(t *SlideWindowIPLimiter) {
		t.interval = interval
	}
}

func WithRate(rate int) option.Option[SlideWindowIPLimiter] {
	return func(t *SlideWindowIPLimiter) {
		t.rate = rate
	}
}

func WithPrefix(prefix string) option.Option[SlideWindowIPLimiter] {
	return func(t *SlideWindowIPLimiter) {
		t.prefix = prefix
	}
}

func (b *SlideWindowIPLimiter) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.accept(ctx)
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

func (b *SlideWindowIPLimiter) accept(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cli.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
