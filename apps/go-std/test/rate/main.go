package main

import (
	"fmt"
	"go-std/internal/utils"
	"time"
)

func main() {
	fw := utils.NewFixedWindowRateLimiter("test_fw", 5, time.Second*10)

	for i := 0; i < 3; i++ {
		allowed := fw.IsAllowed("user1")
		fmt.Printf("Request %d: %t\n", i+1, allowed)
	}

	for i := 0; i < 3; i++ {
		allowed := fw.IsAllowed("user1")
		allowed2 := fw.IsAllowed("user2")
		fmt.Printf("Request %d: %t, %t\n", i+1, allowed, allowed2)
	}

	time.Sleep(time.Second * 11)

	for i := 0; i < 3; i++ {
		allowed := fw.IsAllowed("user1")
		fmt.Printf("Request %d: %t\n", i+1, allowed)
	}
}
