package common

import (
	"errors"
	"testing"
)

func TestTTLCacheSetMethod(t *testing.T) {
	cache := NewTTLCache(10000)
	cache.Set("1", "2")
	if val, exist := cache.Get("1"); !exist || val == nil || val.(string) != "2" {
		t.Fatal("Failed to set and retrive value from cache", exist, val)
	}
}

func TestTTLCacheGetExpiredKeyMethod(t *testing.T) {
	cache := NewTTLCache(0)
	cache.Set("1", "2")
	if val, exist := cache.Get("1"); exist || val != nil {
		t.Fatal("Expected to get nothing when key expired", exist, val)
	}
}

func TestTTLCacheSetAndDeleteMethods(t *testing.T) {
	cache := NewTTLCache(10000)
	cache.Set("1", "2")
	cache.Delete("1")
	if val, exist := cache.Get("1"); exist || val != nil {
		t.Fatal("Failed to set, delete and retrive value from cache", exist, val)
	}
}

func TestTTLCacheGetOrSetMethod(t *testing.T) {
	cache := NewTTLCache(10000)

	val0, loaded0 := cache.GetOrSet("1", func() (interface{}, error) {
		return nil, errors.New("error")
	})
	if loaded0 || val0 != nil {
		t.Fatal("Expected to get nothing if set function emit error", loaded0, val0)
	}

	val1, loaded1 := cache.GetOrSet("1", func() (interface{}, error) {
		return "2", nil
	})
	if loaded1 || val1 == nil || val1.(string) != "2" {
		t.Fatal("Expected to get 2 by key 1", loaded1, val1)
	}

	val2, loaded2 := cache.GetOrSet("1", func() (interface{}, error) {
		return "3", nil
	})
	if !loaded2 || val2 == nil || val2.(string) != "2" {
		t.Fatal("Expected to get also 2 by key 1", loaded2, val2)
	}
}

func TestTTLCacheGC(t *testing.T) {
	cache := NewTTLCache(10)
	cache.Set("1", "10")
	cache.Set("2", "20")
	cache.Set("3", "30")
	cache.Set("4", "40")

	cache.(*ttlCache).gc(255, int64(0xFFFFFFFFFFFFFFF))

	if len(cache.(*ttlCache).items) != 0 {
		t.Fatal("Expected gc to clean items", cache.(*ttlCache).items)
	}
}
