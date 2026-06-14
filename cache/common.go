package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RDB struct {
	*redis.Client
}

func (r *RDB) StoreVerificationCode(email string, code string) error {
	if r == nil || r.Client == nil {
		return errors.New("redis client is not initialized")
	}
	//todo 生产环境建议对key进行哈希处理
	key := "verify:email:" + email
	//设置验证码，5分钟过期
	err := r.Set(context.Background(), key, code, 5*time.Minute).Err()
	return err
}
func (r *RDB) GetVerificationCode(email string) (string, error) {
	if r == nil || r.Client == nil {
		return "", errors.New("redis client is not initialized")
	}
	key := "verify:email:" + email
	val, err := r.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
func NewRDB(ctx context.Context) *RDB {
	return &RDB{NewRedisClient(ctx)}
}
