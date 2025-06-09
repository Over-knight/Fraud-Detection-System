package storage

import (
	"github.com/redis/go-redis/v9"
)

type NewRedisClientFunc func(string) *redis.Client

func NewRedisClient() (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
	})
}

