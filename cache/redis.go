package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Valkey *redis.Client

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadConfig() *Config {
	db, _ := strconv.Atoi(getEnv("VALKEY_DB", "0"))
	return &Config{
		Host:     getEnv("VALKEY_HOST", "localhost"),
		Port:     getEnv("VALKEY_PORT", "6379"),
		Password: getEnv("VALKEY_PASSWORD", ""),
		DB:       db,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Connect() (*redis.Client, error) {
	config := LoadConfig()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	Valkey = client
	log.Println("Redis connected successfully")
	return client, nil
}

func Close() error {
	if Valkey != nil {
		return Valkey.Close()
	}
	return nil
}

// Helper functions
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Valkey.Set(ctx, key, value, expiration).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return Valkey.Get(ctx, key).Result()
}

func Del(ctx context.Context, keys ...string) error {
	return Valkey.Del(ctx, keys...).Err()
}

func Exists(ctx context.Context, keys ...string) (int64, error) {
	return Valkey.Exists(ctx, keys...).Result()
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Valkey.Expire(ctx, key, expiration).Err()
}
