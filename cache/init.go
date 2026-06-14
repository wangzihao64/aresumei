package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

func InitRedis(cfg RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}
	redisClient = client
	return nil
}
func NewRedisClient(_ context.Context) *redis.Client {
	return redisClient
}
func CloseRedis() error {
	if redisClient == nil {
		return nil
	}
	return redisClient.Close()
}
