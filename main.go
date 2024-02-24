package main

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cacher interface {
	Get(int) (string, bool)
	Set(int, string) error
	Remove(int) error
}

type NopCache struct{}

func (n *NopCache) Get(int) (string, bool) {
	return "", false
}

func (n *NopCache) Set(int, string) error {
	return nil
}

func (n *NopCache) Remove(int) error {
	return nil
}

type Store struct {
	data  map[int]string
	cache Cacher
}

func NewStore(c Cacher) *Store {
	data := map[int]string{
		1: "Elon Musk is the new owner of Twitter",
		2: "Foo is not Bar",
		3: "Bar is not Foo",
	}
	return &Store{data: data, cache: c}
}

func (s *Store) Get(key int) (string, error) {
	if val, ok := s.cache.Get(key); ok {
		fmt.Printf("get %d from cache value\n", key)
		return val, nil
	}
	if v, ok := s.data[key]; ok {
		if err := s.cache.Set(key, v); err != nil {
			fmt.Printf("failed to save %d key in cache: %v\n", key, err)
		}
		fmt.Printf("get %d from internal storage\n", key)
		return v, nil
	}
	return "", fmt.Errorf("key not found: %d", key)
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	s := NewStore(NewRedisCache(rdb, 4*time.Second))
	for i := 0; i < 3; i++ {
		key := 1
		val, err := s.Get(key)
		if err != nil {
			fmt.Printf("failed to get %d key: %v\n", key, err)
		}
		fmt.Println(val)
		// time.Sleep(5 * time.Second)
	}
}
