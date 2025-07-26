package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config Redis配置
type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Redis Redis缓存
type Redis struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedis 创建Redis缓存
func NewRedis(config Config) (*Redis, error) {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 创建上下文
	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set 设置缓存
func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 设置缓存
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (r *Redis) Get(key string, value interface{}) error {
	// 获取缓存
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}

	// 反序列化值
	return json.Unmarshal(data, value)
}

// Delete 删除缓存
func (r *Redis) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists 检查缓存是否存在
func (r *Redis) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// Expire 设置缓存过期时间
func (r *Redis) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// Close 关闭Redis连接
func (r *Redis) Close() error {
	return r.client.Close()
}

// SetNX 设置缓存（如果不存在）
func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	// 设置缓存
	return r.client.SetNX(r.ctx, key, data, expiration).Result()
}

// Incr 自增
func (r *Redis) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Decr 自减
func (r *Redis) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// HSet 设置哈希表字段
func (r *Redis) HSet(key, field string, value interface{}) error {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 设置哈希表字段
	return r.client.HSet(r.ctx, key, field, data).Err()
}

// HGet 获取哈希表字段
func (r *Redis) HGet(key, field string, value interface{}) error {
	// 获取哈希表字段
	data, err := r.client.HGet(r.ctx, key, field).Bytes()
	if err != nil {
		return err
	}

	// 反序列化值
	return json.Unmarshal(data, value)
}

// HDelete 删除哈希表字段
func (r *Redis) HDelete(key, field string) error {
	return r.client.HDel(r.ctx, key, field).Err()
}

// HExists 检查哈希表字段是否存在
func (r *Redis) HExists(key, field string) (bool, error) {
	return r.client.HExists(r.ctx, key, field).Result()
}