package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCaching(t *testing.T) {
	c := NewCache(time.Millisecond * 5000)

	c.Add("test1", []byte("hello world1"))
	c.Add("test2", []byte("hello world"))
	c.Add("test3", []byte("hello world"))
	c.Add("test4", []byte("hello world"))

	want := []byte("hello world1")
	got, found := c.Get("test1")

	if !found {
		t.Errorf("entry not found in cache")
	}

	if string(got) != string(want) {
		t.Errorf("wanted %v got %v", got, want)
	}

}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://google.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com",
			val: []byte("abc123"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)

	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to NOT find key")
		return
	}
}
