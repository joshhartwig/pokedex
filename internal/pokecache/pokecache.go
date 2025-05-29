package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	Entries  map[string]CacheEntry
	MU       sync.Mutex
	Interval time.Duration
}

// create a new cache
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

// locks the cache and add a new cache entry with the value of v
func (c *Cache) Add(s string, v []byte) {
	c.MU.Lock()

	defer c.MU.Unlock()
	ce := CacheEntry{
		CreatedAt: time.Now(),
		Val:       v,
	}
	c.Entries[s] = ce
}

// returns our cache entry if found
func (c *Cache) Get(s string) ([]byte, bool) {
	c.MU.Lock()

	defer c.MU.Unlock()
	_, ok := c.Entries[s]
	if !ok {
		return []byte{}, false
	}

	return c.Entries[s].Val, true
}

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
