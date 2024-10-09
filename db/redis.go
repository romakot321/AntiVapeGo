package db

import (
  "os"
  "github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
  addr := os.Getenv("REDIS_ADDRESS")
  if len(addr) == 0 {
    addr = "localhost:6379"
  }
  client := redis.NewClient(&redis.Options{
    Addr:     addr,
    DB:       0,              
  })
  return client
}
