package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

type CacheItem struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type Cache interface {
	Get(ctx context.Context, key string) (*CacheItem, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Stats(ctx context.Context) (*CacheStats, error)
}

type CacheStats struct {
	TotalKeys    int64 `json:"total_keys"`
	HitCount     int64 `json:"hit_count"`
	MissCount    int64 `json:"miss_count"`
	HitRatio     float64 `json:"hit_ratio"`
	MemoryUsage  int64 `json:"memory_usage"`
	Connections  int   `json:"connections"`
}

func NewRedisCache(addr, password string, db int, prefix string, ttl time.Duration) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 100,
	})

	return &RedisCache{
		client: rdb,
		prefix: prefix,
		ttl:    ttl,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (*CacheItem, error) {
	fullKey := r.getFullKey(key)
	
	val, err := r.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var item CacheItem
	if err := json.Unmarshal([]byte(val), &item); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if time.Now().After(item.ExpiresAt) {
		r.Delete(ctx, key)
		return nil, ErrCacheExpired
	}

	return &item, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = r.ttl
	}

	item := CacheItem{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	fullKey := r.getFullKey(key)
	return r.client.Set(ctx, fullKey, data, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.getFullKey(key)
	return r.client.Del(ctx, fullKey).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.getFullKey(key)
	count, err := r.client.Exists(ctx, fullKey).Result()
	return count > 0, err
}

func (r *RedisCache) Clear(ctx context.Context) error {
	pattern := r.getFullKey("*")
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (r *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	fullPattern := r.getFullKey(pattern)
	keys, err := r.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = key[len(r.prefix)+1:]
	}
	return result, nil
}

func (r *RedisCache) Stats(ctx context.Context) (*CacheStats, error) {
	_, err := r.client.Info(ctx, "memory", "stats").Result()
	if err != nil {
		return nil, err
	}

	stats := &CacheStats{
		TotalKeys:   0,
		HitCount:    0, 
		MissCount:   0, 
		HitRatio:    0, 
		MemoryUsage: 0, 
		Connections: 0, 
	}

	// Info string'ini parse et ve stats'ı doldur
	// Bu kısım Redis info formatına göre implement edilecek

	return stats, nil
}

func (r *RedisCache) getFullKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

var (
	ErrCacheMiss    = fmt.Errorf("cache miss")
	ErrCacheExpired = fmt.Errorf("cache expired")
)

type DistributedCache struct {
	primary   Cache
	secondary Cache
	strategy  CacheStrategy
}

type CacheStrategy string

const (
	WriteThrough CacheStrategy = "write_through"
	WriteBack    CacheStrategy = "write_back"
	WriteAround  CacheStrategy = "write_around"
)

func NewDistributedCache(primary, secondary Cache, strategy CacheStrategy) *DistributedCache {
	return &DistributedCache{
		primary:   primary,
		secondary: secondary,
		strategy:  strategy,
	}
}

func (d *DistributedCache) Get(ctx context.Context, key string) (*CacheItem, error) {
	item, err := d.primary.Get(ctx, key)
	if err == nil {
		return item, nil
	}

	item, err = d.secondary.Get(ctx, key)
	if err == nil {
		d.primary.Set(ctx, key, item.Value, time.Until(item.ExpiresAt))
		return item, nil
	}

	return nil, err
}

func (d *DistributedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	switch d.strategy {
	case WriteThrough:
		err1 := d.primary.Set(ctx, key, value, ttl)
		err2 := d.secondary.Set(ctx, key, value, ttl)
		if err1 != nil {
			return err1
		}
		return err2
	case WriteBack:
		err := d.primary.Set(ctx, key, value, ttl)
		go d.secondary.Set(context.Background(), key, value, ttl)
		return err
	case WriteAround:
		return d.primary.Set(ctx, key, value, ttl)
	default:
		return d.primary.Set(ctx, key, value, ttl)
	}
}