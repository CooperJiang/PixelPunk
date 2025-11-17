package health

import (
	"pixelpunk/pkg/cache"
	"time"
)

// CacheChecker 缓存健康检查器
type CacheChecker struct{}

// Name 返回检查项名称
func (c *CacheChecker) Name() string {
	return "cache"
}

// Check 执行健康检查
func (c *CacheChecker) Check() (Status, map[string]interface{}) {
	cacheInstance := cache.GetCache()
	if cacheInstance == nil {
		return StatusDown, map[string]interface{}{
			"error": "缓存未初始化",
		}
	}

	testKey := "health_check_test_" + time.Now().Format("20060102150405")
	testValue := "ok"

	start := time.Now()
	err := cacheInstance.Set(testKey, testValue, 5*time.Second)
	if err != nil {
		return StatusDown, map[string]interface{}{
			"error": "缓存写入失败: " + err.Error(),
		}
	}
	setTime := time.Since(start)

	start = time.Now()
	value, err := cacheInstance.Get(testKey)
	if err != nil {
		return StatusDown, map[string]interface{}{
			"error": "缓存读取失败: " + err.Error(),
		}
	}
	getTime := time.Since(start)

	if value != testValue {
		return StatusDown, map[string]interface{}{
			"error": "缓存值不一致, 期望值:" + testValue + ", 实际值:" + value,
		}
	}

	_ = cacheInstance.Del(testKey)

	details := map[string]interface{}{
		"set_time_ms": setTime.Milliseconds(),
		"get_time_ms": getTime.Milliseconds(),
		"type":        "memory",
	}

	if cache.IsRedisEnabled() {
		details["type"] = "redis"
	}

	// 如果响应时间超过50ms, 视为性能下降
	if setTime > 50*time.Millisecond || getTime > 50*time.Millisecond {
		return StatusDegraded, details
	}

	return StatusUp, details
}

// Type 返回检查类型
func (c *CacheChecker) Type() CheckType {
	return CheckTypeBasic
}

// init 注册缓存健康检查器
func init() {
	RegisterChecker(&CacheChecker{})
}
