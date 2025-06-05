package pokecache

import (
	"sync"
	"time"
)

// CacheEntry represents a single entry in the cache, containing the creation time and the value.
type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

// Cache is a simple in-memory cache that stores entries with a key and a value.
type Cache struct {
	Entries  map[string]CacheEntry
	MU       sync.Mutex
	Interval time.Duration
}

// NewCache initializes a new Cache with a specified interval for reaping old entries.
func NewCache(interval time.Duration) *Cache {
	c := Cache{}
	c.Entries = make(map[string]CacheEntry)
	c.Interval = interval
	ticker := time.NewTicker(interval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				c.ReapLoop()
			}
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	return &c
}

// Add adds a new entry to the cache with the current time as the creation time.
// If the key already exists, it will overwrite the existing entry.
func (c *Cache) Add(s string, v []byte) {
	c.MU.Lock()

	defer c.MU.Unlock()
	ce := CacheEntry{
		CreatedAt: time.Now(),
		Val:       v,
	}
	c.Entries[s] = ce
}

// Get returns our cache entry if found
func (c *Cache) Get(s string) ([]byte, bool) {
	c.MU.Lock()

	defer c.MU.Unlock()
	_, ok := c.Entries[s]
	if !ok {
		return []byte{}, false
	}

	return c.Entries[s].Val, true
}

// ReapLoop checks each entry in the cache and removes those that are older than the specified interval.
// It is called periodically based on the interval set during cache initialization.
func (c *Cache) ReapLoop() {
	c.MU.Lock()
	defer c.MU.Unlock()

	// if now is after the creation date adding the interval delete the entry
	now := time.Now()
	for k, v := range c.Entries {
		if now.After(v.CreatedAt.Add(c.Interval)) {
			delete(c.Entries, k)
		}
	}
}
