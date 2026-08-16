// Sharded Cache: сегментированный in-memory кеш.
//
// Один мьютекс на всю мапу — точка конкуренции: под нагрузкой все горутины
// сериализуются. Решение — N независимых шардов, каждый со своей мапой и
// RWMutex. Шард выбирается по FNV-хешу ключа; N — степень двойки, чтобы
// вместо деления брать hash & (N-1). Конкуренция падает в ~N раз.
package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Cache interface {
	Set(k string, v string)
	Get(k string) (string, bool)
}

const shardCount = 256 // степень двойки

type shard struct {
	mu   sync.RWMutex
	data map[string]string
}

type ShardedCache struct {
	shards [shardCount]*shard
}

var _ Cache = (*ShardedCache)(nil)

func NewShardedCache() *ShardedCache {
	c := &ShardedCache{}
	for i := range c.shards {
		c.shards[i] = &shard{data: make(map[string]string)}
	}
	return c
}

func (c *ShardedCache) getShard(k string) *shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(k))
	return c.shards[h.Sum32()&(shardCount-1)]
}

func (c *ShardedCache) Set(k, v string) {
	s := c.getShard(k)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[k] = v
}

func (c *ShardedCache) Get(k string) (string, bool) {
	s := c.getShard(k)

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[k]
	return v, ok
}

func main() {
	cache := NewShardedCache()

	// конкурентная запись/чтение
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("key-%d", i)
			cache.Set(k, fmt.Sprintf("value-%d", i))
			if v, ok := cache.Get(k); !ok || v == "" {
				panic("lost write")
			}
		}(i)
	}
	wg.Wait()

	fmt.Println(cache.Get("key-42"))  // value-42 true
	fmt.Println(cache.Get("missing")) // "" false
}
