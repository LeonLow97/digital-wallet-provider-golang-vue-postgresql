package infrastructure

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Dependency Injection of Golang and Redis with unit tests
// https://www.razvanh.com/blog/testing-golang-redis-dependency
// https://redis.io/docs/connect/clients/go/

// RedisClient interface defines the methods for interacting with Redis
type RedisClient interface {
	Close() error

	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)

	SAdd(ctx context.Context, key string, members ...interface{}) error
}

// RedisClientImpl is the concrete implementation of RedisClient
type RedisClientImpl struct {
	client *redis.Client
}

// NewRedisClient creates a new instance of RedisClientImpl
func NewRedisClient() RedisClient {
	// initialize redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379", // redis server address
		Password: "",                  // No Password
		DB:       0,                   // Default DB
	})

	// check if Redis is reachable via Ping command
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis server", err)
	}
	log.Println("Connected to Redis:", pong)

	return &RedisClientImpl{
		client: redisClient,
	}
}

func (rc *RedisClientImpl) Close() error {
	return rc.client.Close()
}

func (rc *RedisClientImpl) Set(ctx context.Context, key string, value interface{}) error {
	return rc.client.Set(ctx, key, value, 0).Err()
}

func (rc *RedisClientImpl) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

func (rc *RedisClientImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rc.client.SAdd(ctx, key, members).Err()
}
