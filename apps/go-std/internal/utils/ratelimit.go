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

func NewTokenBucketRateLimiter(prefix string, bucketCapacity int, refillRate float64) *TokenBucketRateLimiter {
	storage := NewTTLMap(10000, time.Minute*15)
	cachePrefix := "tb_rate_limit:" + prefix + ":"

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

type LeakyBucketCache struct {
	count    int
	lastLeak int64
}

type LeakyBucketRateLimiter struct {
	cache          Cache
	cachePrefix    string
	bucketCapacity int
	leakRate       float64
}

func NewLeakyBucketRateLimiter(prefix string, bucketCapacity int, leakRate float64) *LeakyBucketRateLimiter {
	storage := NewTTLMap(10000, time.Minute*15)
	cachePrefix := "lb_rate_limit:" + prefix + ":"

	return &LeakyBucketRateLimiter{
		cache:          storage,
		cachePrefix:    cachePrefix,
		bucketCapacity: bucketCapacity,
		leakRate:       leakRate,
	}

}

func (lb *LeakyBucketRateLimiter) IsAllowed(key string) bool {
	cacheKey := lb.cachePrefix + key
	var currentTime int64 = time.Now().UnixMilli()

	results, ok := lb.cache.Get(cacheKey)
	if !ok {
		results = &LeakyBucketCache{
			count:    0,
			lastLeak: currentTime,
		}
	}
	bucket := results.(*LeakyBucketCache)
	requestCount, lastLeakTime := bucket.count, bucket.lastLeak

	elapsedTimeMs := currentTime - lastLeakTime
	elapsedTimeSeconds := float64(elapsedTimeMs) / 1000.0
	requestsToLeak := int(lb.leakRate * elapsedTimeSeconds)
	requestCount = int(math.Max(0.0, float64(requestCount-requestsToLeak)))

	isAllowed := requestCount < lb.bucketCapacity
	if isAllowed {
		requestCount++
	}

	lb.cache.Put(cacheKey, &LeakyBucketCache{
		count:    requestCount,
		lastLeak: currentTime,
	})

	return isAllowed
}
