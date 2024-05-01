package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Dependency Injection of Golang and Redis with unit tests
// https://www.razvanh.com/blog/testing-golang-redis-dependency
// https://redis.io/docs/connect/clients/go/

// RedisClient interface defines the methods for interacting with Redis
type RedisClient interface {
	Close() error

	Set(ctx context.Context, key string, value interface{}) error
	SetEx(ctx context.Context, key string, member interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetEx(ctx context.Context, key string, expiration time.Duration) (string, error)

	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error

	Del(ctx context.Context, keys ...string) error

	Expire(ctx context.Context, key string, sessionExpiryInMinutes time.Duration) error

	TTL(ctx context.Context, key string) (time.Duration, error)
}

// RedisClientImpl is the concrete implementation of RedisClient
type RedisClientImpl struct {
	client *redis.Client
}

// NewRedisClient creates a new instance of RedisClientImpl
func NewRedisClient() RedisClient {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// initialize redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr, // redis server address
		Password: "",        // No Password
		DB:       0,         // Default DB
	})

	// check if Redis is reachable via Ping command
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis server: %v\n", err)
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

func (rc *RedisClientImpl) SetEx(ctx context.Context, key string, member interface{}, expiration time.Duration) error {
	return rc.client.SetEx(ctx, key, member, expiration).Err()
}

func (rc *RedisClientImpl) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

func (rc *RedisClientImpl) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return rc.client.GetEx(ctx, key, expiration).Result()
}

func (rc *RedisClientImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rc.client.SAdd(ctx, key, members...).Err()
}

func (rc *RedisClientImpl) SRem(ctx context.Context, key string, members ...interface{}) error {
	return rc.client.SRem(ctx, key, members...).Err()
}

func (rc *RedisClientImpl) Del(ctx context.Context, keys ...string) error {
	return rc.client.Del(ctx, keys...).Err()
}

func (rc *RedisClientImpl) Expire(ctx context.Context, key string, sessionExpiryInMinutes time.Duration) error {
	return rc.client.Expire(ctx, key, sessionExpiryInMinutes).Err()
}

// (used in DEVELOPMENT) For testing Redis string expiry to ensure key expiry is set correctly
func (rc *RedisClientImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rc.client.TTL(ctx, key).Result()
}
