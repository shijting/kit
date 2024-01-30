package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/shijting/kit/cache"
)

// GinLimiter gin 限流装饰器
func GinLimiter(cap, rate int64) func(handler gin.HandlerFunc) gin.HandlerFunc {
	bucket := NewBucket(cap, rate)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if !bucket.Accept() {
				ctx.AbortWithStatusJSON(429, gin.H{"message": "Too many requests"})
				return
			}
			handler(ctx)
			ctx.Next()
		}
	}
}

// GinQueryLimiter gin query
// key: query key example: /api?limit=xx key: limit 有值时才限流
func GinQueryLimiter(cap, rate int64, key string) func(handler gin.HandlerFunc) gin.HandlerFunc {
	bucket := NewBucket(cap, rate)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if ctx.Query(key) != "" {
				if !bucket.Accept() {
					ctx.AbortWithStatusJSON(429, gin.H{"message": "Too many requests"})
					return
				}
			}
			handler(ctx)
			ctx.Next()
		}
	}
}

// IPTokenBucketLimiter ip 限流
// 可以使用redis实现分布式限流
type IPTokenBucketLimiter struct {
	cache *cache.LRUCache[string, *Bucket]
}

// NewIPTokenBucketLimiter 创建 IPTokenBucketLimiter 实例。
// maxIPs: 最大允许的 IP 数量
func NewIPTokenBucketLimiter(maxIPs int) *IPTokenBucketLimiter {
	return &IPTokenBucketLimiter{cache: cache.NewLRUCache[string, *Bucket](maxIPs)}
}

// Build 返回一个 gin 中间件函数，用于限制 IP 请求频率。
// cap: 桶容量。
// rate: 每秒产生令牌数量。
func (i *IPTokenBucketLimiter) Build(cap, rate int64) func(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ip := ctx.ClientIP()

			item, ok := i.cache.Get(ip)

			var bucket *Bucket

			if !ok {
				bucket = NewBucket(cap, rate)
				i.cache.Set(ip, bucket, 0)
			} else {
				bucket = item.Value
			}

			if !bucket.Accept() {
				ctx.AbortWithStatusJSON(429, gin.H{"message": "Too many requests"})
				return
			}

			handler(ctx)
			ctx.Next()
		}
	}
}
