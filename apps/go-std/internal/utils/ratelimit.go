package utils

import (
	"math"
	"time"
)

type TokenBucketCache struct {
	count      int
	lastRefill int64
}

type TokenBucketRateLimiter struct {
	cache          Cache
	cachePrefix    string
	bucketCapacity int
	refillRate     float64
}

func NewTokenBucketRateLimiter(identifier string, bucketCapacity int, refillRate float64) *TokenBucketRateLimiter {
	storage := NewTTLMap(10000, time.Minute*15, "token_bucket_rate_limiter:"+identifier)
	cachePrefix := "tb_rate_limit:" + identifier + ":"

	return &TokenBucketRateLimiter{
		cache:          storage,
		cachePrefix:    cachePrefix,
		bucketCapacity: bucketCapacity,
		refillRate:     refillRate,
	}

}

func (tb *TokenBucketRateLimiter) IsAllowed(key string) bool {
	cacheKey := tb.cachePrefix + key
	var currentTime int64 = time.Now().UnixMilli()
	//current state
	results, ok := tb.cache.Get(cacheKey)
	if !ok {
		results = &TokenBucketCache{
			count:      tb.bucketCapacity,
			lastRefill: currentTime,
		}
	}
	bucket := results.(*TokenBucketCache)
	tokenCount, lastRefillTime := bucket.count, bucket.lastRefill

	elapsedTimeMs := currentTime - lastRefillTime
	elapsedTimeSeconds := float64(elapsedTimeMs) / 1000.0
	tokensToAdd := int(tb.refillRate * elapsedTimeSeconds)
	tokenCount = int(math.Min(float64(tb.bucketCapacity), float64(tokenCount+tokensToAdd)))

	isAllowed := tokenCount > 0
	if isAllowed {
		tokenCount--
	}

	tb.cache.Put(cacheKey, &TokenBucketCache{
		count:      tokenCount,
		lastRefill: currentTime,
	})

	return isAllowed
}

type FixedWindowCache struct {
	count int
}

type FixedWindowRateLimiter struct {
	cache       Cache
	cachePrefix string
	limit       int
	windowSize  time.Duration
}

func NewFixedWindowRateLimiter(identifier string, limit int, windowSize time.Duration) *FixedWindowRateLimiter {
	storage := NewTTLMap(10000, windowSize, "fixed_window_rate_limiter:"+identifier)
	cachePrefix := "fw_rate_limit:" + identifier + ":"

	return &FixedWindowRateLimiter{
		cache:       storage,
		cachePrefix: cachePrefix,
		limit:       limit,
		windowSize:  windowSize,
	}
}
func (fw *FixedWindowRateLimiter) IsAllowed(key string) bool {
	cacheKey := fw.cachePrefix + key
	results, ok := fw.cache.Get(cacheKey)
	if !ok {
		results = &FixedWindowCache{
			count: 0,
		}
	}
	bucket := results.(*FixedWindowCache)
	count := bucket.count

	isAllowed := count < fw.limit
	if isAllowed {
		count++
		fw.cache.Put(cacheKey, &FixedWindowCache{
			count: count,
		})
	}

	return isAllowed

}
