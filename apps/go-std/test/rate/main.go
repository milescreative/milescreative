package main

import (
	"fmt"
	"go-std/internal/utils"
	"math/rand"
	"time"
)

func main() {
	// Shared parameters
	bucketCapacity := 5
	rate := 1.0 // Requests/tokens per second

	// Create the rate limiters
	tokenBucket := utils.NewTokenBucketRateLimiter("test_token", bucketCapacity, rate)
	leakyBucket := utils.NewLeakyBucketRateLimiter("test_leaky", bucketCapacity, rate)

	fmt.Println("Running test to demonstrate DIFFERENT behavior with multiple users...")
	fmt.Printf("Bucket Capacity: %d, Rate: %.2f requests/second\n", bucketCapacity, rate)

	users := []string{"user1", "user2"}

	// --- Phase 1: Initial Burst with Variability ---
	fmt.Println("\n--- Phase 1: Initial Burst with Variability ---")
	for i := 0; i < 10; i++ { // Send 10 requests in a row for each user
		for _, user := range users {
			tokenAllowed := tokenBucket.IsAllowed(user)
			leakyAllowed := leakyBucket.IsAllowed(user)

			fmt.Printf("Burst Request %d for %s: Token = %t, Leaky = %t\n", i+1, user, tokenAllowed, leakyAllowed)
		}
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Random delay
	}

	// --- Phase 2: Sustained Traffic with Variability ---
	fmt.Println("\n--- Phase 2: Sustained Traffic with Variability (1 request per second) ---")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond) // Random delay between 0.5 to 1.5 seconds

		for _, user := range users {
			tokenAllowed := tokenBucket.IsAllowed(user)
			leakyAllowed := leakyBucket.IsAllowed(user)

			fmt.Printf("Sustained Request %d for %s: Token = %t, Leaky = %t\n", i+1, user, tokenAllowed, leakyAllowed)
		}
	}

	// --- Phase 3: Post-Burst Recovery with Variability ---
	fmt.Println("\n--- Phase 3: Post-Burst Recovery with Variability ---")
	time.Sleep(5 * time.Second) // Wait to allow recovery

	for i := 0; i < 10; i++ { // Send another 10 requests in a row for each user
		for _, user := range users {
			tokenAllowed := tokenBucket.IsAllowed(user)
			leakyAllowed := leakyBucket.IsAllowed(user)

			fmt.Printf("Recovery Burst Request %d for %s: Token = %t, Leaky = %t\n", i+1, user, tokenAllowed, leakyAllowed)
		}
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Random delay
	}
}
