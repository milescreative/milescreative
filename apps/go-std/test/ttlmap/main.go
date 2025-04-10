package main

import (
	"log"
	"time"

	"go-std/internal/utils" // Update this import path to match your module structure
)

func main() {
	m := utils.NewTTLMap(100, time.Second*100, "test-map-1")
	g := utils.NewTTLMap(5, time.Second*100, "test-map-2")

	m.Put("test", "test")
	value, ok := m.Get("test")
	if !ok {
		log.Println("Value not found")
	}
	log.Println(value)

	value, ok = g.Get("test")
	if !ok {
		log.Println("Value not found- test passed")
	}
	log.Println(value)

	// Add multiple items to trigger the warning
	m.Put("test2", "test2")
	m.Put("test3", "test3")
	m.Put("test4", "test4")
	m.Put("test5", "test5")
	m.Put("test6", "test6")
	m.Put("test7", "test7")
	m.Put("test8", "test8")
	m.Put("test9", "test9")

	// Wait to see both the TTL cleanup and the warning
	time.Sleep(15 * time.Second)

	g.Put("g-test", "g-test")
	value, ok = m.Get("test")
	if !ok {
		log.Println("Value not found")
	}
	log.Println(value)
	log.Println("Map 1 name:", m.GetName())

	value, ok = g.Get("g-test")
	if !ok {
		log.Println("Value not found")
	}
	log.Println(value)
	log.Println("Map 2 name:", g.GetName())
}
