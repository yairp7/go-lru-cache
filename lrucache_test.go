package lrucache

import (
	"sync"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func TestSetAndGet(t *testing.T) {
	cache := NewLRUCache[string](10)

	cache.Set("a", "b")
	cache.Set("b", "c")
	cache.Set("c", "d")

	require.Equal(t, "d", cache.Get("c"))
	require.Equal(t, "c", cache.Get("b"))
	require.Equal(t, "b", cache.Get("a"))

	newKeysOrder := []string{"a", "b", "c"}
	ptr := cache.head
	i := 0
	for ptr != nil {
		require.Equal(t, newKeysOrder[i], ptr.key)
		ptr = ptr.next
	}
}

func TestSetAndGetLRU(t *testing.T) {
	cache := NewLRUCache[string](3)

	cache.Set("a", "a")
	cache.Set("b", "b")
	cache.Set("c", "c")
	cache.Set("d", "d")
	cache.Set("e", "e")

	require.Equal(t, 3, cache.size)

	require.Equal(t, "e", cache.head.key)
	require.Equal(t, "c", cache.tail.key)

	cache.Get("d")
	require.Equal(t, "d", cache.head.key)
	require.Equal(t, "c", cache.tail.key)

	cache.Get("c")
	require.Equal(t, "c", cache.head.key)
	require.Equal(t, "e", cache.tail.key)
}

func TestSetRace(t *testing.T) {
	cache := NewLRUCache[int](100000)
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 0; i < 10000; i++ {
			cache.Set(i, i)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 10000; i < 20000; i++ {
			cache.Set(i, i)
		}
	}()

	waitGroup.Wait()
}

func TestGetRace(t *testing.T) {
	cache := NewLRUCache[int](100000)
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		cache.Set(i, i)
	}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 0; i < 10000; i++ {
			cache.Get(i)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for i := 0; i < 10000; i++ {
			cache.Get(i)
		}
	}()

	waitGroup.Wait()
}

func BenchmarkSet(b *testing.B) {
	cache := NewLRUCache[int](b.N)

	for i := 0; i < b.N; i++ {
		cache.Set(i, i)
	}
}

func BenchmarkSetWithLRU(b *testing.B) {
	cache := NewLRUCache[int](10)

	for i := 0; i < b.N; i++ {
		cache.Set(i, i)
	}
}
