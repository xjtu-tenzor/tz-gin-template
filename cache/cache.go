package cache

import (
	"context"
	"encoding/json"
	"template/config"
	"time"
)

// SetCache 为常用的页面或查询设置缓存
func SetCache(key string, value interface{}, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return config.RDB.Set(ctx, key, val, expiry).Err()
}

// GetCache
func GetCache(key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := config.RDB.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}
