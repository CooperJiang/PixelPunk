package health

import (
	"runtime"
	"sync"
	"time"
)

// Status 表示健康检查状态
type Status string

const (
	// StatusUp 表示组件正常运行
	StatusUp Status = "up"
	// StatusDown 表示组件不可用
	StatusDown Status = "down"
	// StatusDegraded 表示组件性能下降
	StatusDegraded Status = "degraded"
)

// CheckType 定义健康检查类型
type CheckType string

const (
	// CheckTypeBasic 仅检查关键组件
	CheckTypeBasic CheckType = "basic"
	// CheckTypeComplete 全面检查所有组件
	CheckTypeComplete CheckType = "complete"
)

// Checker 定义健康检查接口
type Checker interface {
	// Name 返回检查项名称
	Name() string
	// Check 执行健康检查
	Check() (Status, map[string]interface{})
	// Type 返回检查类型
	Type() CheckType
}

// Result 表示单个健康检查结果
type Result struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	CheckTime time.Time              `json:"check_time"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// SystemHealth 表示系统健康状态
type SystemHealth struct {
	Status     Status         `json:"status"`
	Version    string         `json:"version"`
	CheckTime  time.Time      `json:"check_time"`
	Components []Result       `json:"components"`
	SystemInfo map[string]any `json:"system_info,omitempty"`
	Uptime     time.Duration  `json:"uptime"`
	StartTime  time.Time      `json:"start_time"`
}

var (
	// 健康检查项列表
	checkers []Checker
	// 防止并发注册检查项
	registerMutex sync.Mutex
	// 应用启动时间
	startTime = time.Now()
	// 应用版本
	appVersion = "1.0.0"
)

// RegisterChecker 注册健康检查项
func RegisterChecker(checker Checker) {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	for _, c := range checkers {
		if c.Name() == checker.Name() {
			return
		}
	}

	checkers = append(checkers, checker)
}

func SetVersion(version string) {
	appVersion = version
}

// Check 执行健康检查
func Check(checkType CheckType) SystemHealth {
	now := time.Now()
	results := make([]Result, 0)

	// 系统总体状态默认为正常
	overallStatus := StatusUp

	for _, checker := range checkers {
		if checkType != CheckTypeComplete && checker.Type() != CheckTypeBasic {
			continue
		}

		status, details := checker.Check()
		results = append(results, Result{
			Name:      checker.Name(),
			Status:    status,
			CheckTime: now,
			Details:   details,
		})

		// 如果任一组件状态为 Down，则系统状态为 Down
		if status == StatusDown {
			overallStatus = StatusDown
		} else if status == StatusDegraded && overallStatus != StatusDown {
			// 如果任一组件状态为 Degraded 且当前系统状态不是 Down，则系统状态为 Degraded
			overallStatus = StatusDegraded
		}
	}

	systemInfo := getSystemInfo()

	return SystemHealth{
		Status:     overallStatus,
		Version:    appVersion,
		CheckTime:  now,
		Components: results,
		SystemInfo: systemInfo,
		Uptime:     time.Since(startTime),
		StartTime:  startTime,
	}
}

// getSystemInfo 获取系统信息
func getSystemInfo() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"go_version":   runtime.Version(),
		"goroutines":   runtime.NumGoroutine(),
		"cpu_num":      runtime.NumCPU(),
		"memory_alloc": memStats.Alloc,
		"memory_total": memStats.TotalAlloc,
		"memory_sys":   memStats.Sys,
		"gc_num":       memStats.NumGC,
		"os":           runtime.GOOS,
		"arch":         runtime.GOARCH,
	}
}
