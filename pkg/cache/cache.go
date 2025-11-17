package cache

import (
	"fmt"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/logger"
	"strings"
	"time"
)

// Cache 定义缓存接口
type Cache interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Exists(key string) bool
	TTL(key string) (time.Duration, error)
	Expire(key string, expiration time.Duration) error
	Close() error
}

var (
	defaultCache Cache
	namespace    string
)

// IsRedisEnabled 检查Redis是否启用
func IsRedisEnabled() bool {
	return redisCache != nil && defaultCache == redisCache
}

func InitCache() {

	initNamespace()

	if err := InitRedis(); err != nil {
		logger.Error("Redis初始化失败: %v，将使用内存缓存", err)
		defaultCache = InitMemCache()
	} else {
	}

}

// initNamespace 初始化命名空间
func initNamespace() {
	cfg := config.GetConfig()
	ns := strings.TrimSpace(cfg.App.Namespace)

	if ns == "" {
		ns = "pixelpunk"
		logger.Warn("未配置APP_NS，使用默认命名空间: %s", ns)
	}

	namespace = ns
}

// buildKey 构建带命名空间前缀的缓存键
func buildKey(key string) string {
	if namespace == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", namespace, key)
}

func GetCache() Cache {
	if defaultCache == nil {
		InitCache()
	}
	return defaultCache
}

// Set 设置缓存（自动添加命名空间前缀）
func Set(key string, value string, expiration time.Duration) error {
	return GetCache().Set(buildKey(key), value, expiration)
}

// Get 获取缓存（自动添加命名空间前缀）
func Get(key string) (string, error) {
	return GetCache().Get(buildKey(key))
}

// Del 删除缓存（自动添加命名空间前缀）
func Del(key string) error {
	return GetCache().Del(buildKey(key))
}

// Exists 检查键是否存在（自动添加命名空间前缀）
func Exists(key string) bool {
	return GetCache().Exists(buildKey(key))
}

// TTL 获取过期时间（自动添加命名空间前缀）
func TTL(key string) (time.Duration, error) {
	return GetCache().TTL(buildKey(key))
}

// Expire 设置过期时间（自动添加命名空间前缀）
func Expire(key string, expiration time.Duration) error {
	return GetCache().Expire(buildKey(key), expiration)
}

// Close 关闭缓存连接
func Close() error {
	if defaultCache != nil {
		return defaultCache.Close()
	}
	return nil
}

func GetNamespace() string {
	return namespace
}

// ClearNamespaceCache 清空当前命名空间的所有缓存（慎用）
func ClearNamespaceCache() error {
	if redisCache, ok := defaultCache.(*RedisCache); ok {
		pattern := fmt.Sprintf("%s:*", namespace)
		ctx := redisCache.ctx

		keys, err := redisCache.client.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			return redisCache.client.Del(ctx, keys...).Err()
		}
	}

	logger.Warn("当前缓存实现不支持批量清空，请手动清理或重启Redis")
	return nil
}
