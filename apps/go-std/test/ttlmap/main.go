package main

import (
	"log"
	"sync"
	"time"
)

// item is a struct that holds the value and the last access time
type item struct {
	value      interface{}
	lastAccess int64
}

// You can have a single map for an application or few maps for different purposes
type TTLMap struct {
	m map[string]*item
	// For safe access to the map
	mu sync.Mutex
}

func New(size int, maxTTL int) (m *TTLMap) {
	// map is created with the given length
	m = &TTLMap{m: make(map[string]*item, size)}

	// this goroutine will clean up the map from old itemsa
	go func() {
		// You can adjust this ticker to be more or less frequent
		for now := range time.Tick(time.Second) {
			m.mu.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.mu.Unlock()
		}
	}()

	return
}

// Put adds a new item to the map or updates the existing one
func (m *TTLMap) Put(k string, v interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	it, ok := m.m[k]
	if !ok {
		it = &item{
			value: v,
		}
	}
	it.value = v
	it.lastAccess = time.Now().Unix()
	m.m[k] = it
}

// Get returns the value of the given key if it exists
func (m *TTLMap) Get(k string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if it, ok := m.m[k]; ok {
		it.lastAccess = time.Now().Unix()
		return it.value, true
	}

	return nil, false
}

// Delete removes the item from the map
func (m *TTLMap) Delete(k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, k)
}

func main() {
	m := New(100, 10)
	g := New(5, 10)
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

	time.Sleep(15 * time.Second)
	g.Put("g-test", "g-test")
	value, ok = m.Get("test")
	if !ok {
		log.Println("Value not found")
	}
	log.Println(value)
	log.Println("m", m.m)
	value, ok = g.Get("g-test")
	if !ok {
		log.Println("Value not found")
	}
	log.Println(value)
	log.Println("g", g.m)
}
