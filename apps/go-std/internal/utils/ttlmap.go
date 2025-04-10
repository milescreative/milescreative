package utils

import (
	"log"
	"sync"
	"time"
)

type Cache interface {
	Put(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
}

// item is a struct that holds the value and the last access time
type item struct {
	value      interface{}
	lastAccess int64
}

// You can have a single map for an application or few maps for different purposes
type TTLMap struct {
	m map[string]*item
	// For safe access to the map
	mu          sync.Mutex
	name        string
	lastWarning int64
}

func NewTTLMap(size int, maxTTL time.Duration, name string) (m *TTLMap) {
	// map is created with the given length
	m = &TTLMap{m: make(map[string]*item, size), name: name, lastWarning: 0}

	// this goroutine will clean up the map from old items
	go func() {
		// You can adjust this ticker to be more or less frequent
		for now := range time.Tick(time.Second) {
			m.mu.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(maxTTL.Seconds()) {
					delete(m.m, k)
				}
			}

			m.mu.Unlock()
		}
	}()

	//check if map is over threshold every 24 hours
	go func() {
		for range time.Tick(time.Hour * 24) {
			m.HandleTTLMapLimitWarning()
		}
	}()

	return
}

func (m *TTLMap) GetName() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.name
}

// Put adds a new item to the map or updates the existing one
func (m *TTLMap) Put(k string, v interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	it, ok := m.m[k]
	if !ok {
		it = &item{}
		m.m[k] = it
	}
	it.value = v
	it.lastAccess = time.Now().Unix()
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

func (m *TTLMap) HandleTTLMapLimitWarning() {
	m.mu.Lock()
	//if map is over threshold, log warning
	warningThreshold := 1000
	// every 24 hours, log warning
	warningWindow := (time.Hour * 24).Milliseconds()
	count := len(m.m)
	if count > warningThreshold && (time.Now().UnixMilli()-m.lastWarning) > warningWindow {
		log.Println("TLLMap "+m.name+" size of ", count, " is above warning threshold of ", warningThreshold, ".")
		m.lastWarning = time.Now().UnixMilli()
	}
	m.mu.Unlock()
}
