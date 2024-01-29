package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/shijting/kit/syncx"
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
		}
	}
}

// GinQueryLimiter gin query 限流装饰器
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
		}
	}
}

// IPBucketCache ip 限流
// 可以使用redis实现分布式限流
var IPBucketCache = syncx.NewMap[string, *Bucket]()

// IPLimiter ip 限流装饰器
func IPLimiter(cap, rate int64) func(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ip := ctx.ClientIP()
			bucket, ok := IPBucketCache.Load(ip)
			if !ok {
				bucket = NewBucket(cap, rate)
				IPBucketCache.Store(ip, bucket)
			}
			if !bucket.Accept() {
				ctx.AbortWithStatusJSON(429, gin.H{"message": "Too many requests"})
				return
			}
			handler(ctx)
		}
	}
}
