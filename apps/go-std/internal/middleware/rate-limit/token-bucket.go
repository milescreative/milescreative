package ratelimit

import (
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	mu              sync.RWMutex
	memoryStore     = make(map[string]interface{})
	cleanupInterval = 1 * time.Hour
)

type TokenBucketRateLimiter struct {
	mu             sync.Mutex
	bucketCapacity int
	refillRate     float64
	prefix         string
	store          BucketStore
}

type BucketStore struct {
	getCurrentState func(keyLastRefill string, keyCount string) (int64, int)
	setCurrentState func(keyLastRefill string, keyCount string, lastRefill int64, count int)
}

var rateLimitStorage = "memory"

func init() {
	log.Printf("Starting rate limit cleanup goroutine with interval: %v", cleanupInterval)
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("Running rate limit store cleanup. Current size: %d", len(memoryStore))
			cleanupStore()
			log.Printf("Cleanup complete. New size: %d", len(memoryStore))
		}
	}()
}

func cleanupStore() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UnixMilli()
	threshold := now - (24 * time.Hour).Milliseconds() // Maybe make this shorter for testing
	removed := 0

	for key, value := range memoryStore {
		if strings.HasSuffix(key, ":tb_lastRefill") {
			lastRefill := value.(int64)
			if lastRefill < threshold {
				// Remove both the lastRefill and count entries
				countKey := strings.TrimSuffix(key, ":tb_lastRefill") + ":tb_count"
				delete(memoryStore, key)
				delete(memoryStore, countKey)
				removed++
			}
		}
	}

	if removed > 0 {
		log.Printf("Removed %d expired entries from rate limit store", removed)
	}
}

func NewTokenBucketRateLimiter(bucketCapacity int, refillRate float64, prefix string) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		bucketCapacity: bucketCapacity,
		refillRate:     refillRate,
		prefix:         prefix,
		store:          setupStore(bucketCapacity),
	}
}

func setupStore(bucketCapacity int) BucketStore {
	if rateLimitStorage == "memory" {
		return BucketStore{
			getCurrentState: func(keyLastRefill string, keyCount string) (int64, int) {
				mu.RLock()
				defer mu.RUnlock()

				lastRefill, ok := memoryStore[keyLastRefill]
				if !ok {
					return time.Now().UnixMilli(), bucketCapacity
				}
				count, ok := memoryStore[keyCount]
				if !ok {
					return lastRefill.(int64), bucketCapacity
				}
				return lastRefill.(int64), count.(int)
			},
			setCurrentState: func(keyLastRefill string, keyCount string, lastRefill int64, count int) {
				mu.Lock()
				defer mu.Unlock()
				memoryStore[keyLastRefill] = lastRefill
				memoryStore[keyCount] = count
			},
		}
	}
	return BucketStore{}
}

func (rl *TokenBucketRateLimiter) formatKeys(key string) (string, string) {
	prefix := rl.prefix
	if prefix != "" {
		prefix = prefix + ":"
	}
	return prefix + key + ":tb_lastRefill", prefix + key + ":tb_count"
}

func (rl *TokenBucketRateLimiter) IsAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	keyLastRefill, keyCount := rl.formatKeys(key)
	lastRefill, count := rl.store.getCurrentState(keyLastRefill, keyCount)

	currentTime := time.Now().UnixMilli()
	elapsedTimeMs := currentTime - lastRefill
	elapsedTimeSecs := float64(elapsedTimeMs) / 1000.0
	tokensToAdd := int(math.Floor(elapsedTimeSecs * rl.refillRate))

	// Only update lastRefill if we're actually adding tokens
	newLastRefill := lastRefill
	if tokensToAdd > 0 {
		newLastRefill = currentTime
		count = int(math.Min(float64(count+tokensToAdd), float64(rl.bucketCapacity)))
	}

	isAllowed := count >= 0
	if isAllowed {
		count--
	}

	// Store the updated state
	rl.store.setCurrentState(keyLastRefill, keyCount, newLastRefill, count)
	return isAllowed
}

func GetStoreSize() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(memoryStore)
}

func GetStoreKeys() []string {
	mu.RLock()
	defer mu.RUnlock()

	keys := make([]string, 0, len(memoryStore))
	for k := range memoryStore {
		keys = append(keys, k)
	}
	return keys
}
