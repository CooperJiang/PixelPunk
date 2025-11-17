// Package sysinfo 提供系统信息获取工具
package sysinfo

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/utils"

	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetSystemInfo() models.SystemInfo {
	result := models.SystemInfo{}

	loadAvg, err := getSystemLoadAvg()
	if err != nil {
		loadAvg = []float64{0, 0, 0}
	}
	result.Load = models.SystemLoad{
		OneMin:     loadAvg[0],
		FiveMin:    loadAvg[1],
		FifteenMin: loadAvg[2],
	}

	memInfo, err := getSystemMemInfo()
	if err != nil {
		memInfo = map[string]uint64{
			"total": 0,
			"used":  0,
			"free":  0,
		}
	}

	// 计算内存使用百分比并保留两位小数
	memUsagePercent := float64(memInfo["used"]) / float64(memInfo["total"]) * 100
	memUsagePercent = float64(int(memUsagePercent*100)) / 100 // 保留两位小数

	result.Memory = models.SystemMemory{
		Total:          memInfo["total"],
		Used:           memInfo["used"],
		Free:           memInfo["free"],
		FormattedTotal: utils.FormatBytes(int64(memInfo["total"])),
		FormattedUsed:  utils.FormatBytes(int64(memInfo["used"])),
		FormattedFree:  utils.FormatBytes(int64(memInfo["free"])),
		UsagePercent:   memUsagePercent,
	}

	cpuInfo := getCPUInfo()

	// 保留CPU使用率两位小数
	cpuInfo.UsagePercent = float64(int(cpuInfo.UsagePercent*100)) / 100

	result.CPU = cpuInfo

	diskInfo := getDiskInfo()

	// 保留磁盘使用率两位小数
	diskInfo.UsagePercent = float64(int(diskInfo.UsagePercent*100)) / 100

	result.Disk = diskInfo

	uptime := getSystemUptime()
	result.Uptime = uptime

	basicInfo := getBasicSystemInfo()
	result.Info = basicInfo

	return result
}

// getSystemLoadAvg 获取系统负载
func getSystemLoadAvg() ([]float64, error) {
	loadAvg := []float64{0, 0, 0}

	loadFile, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		// 如果是Darwin (MacOS)
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("sysctl", "-n", "vm.loadavg")
			output, err := cmd.Output()
			if err != nil {
				return loadAvg, err
			}

			// 解析输出，格式通常是 "{ 0.00 0.00 0.00 }"
			parts := strings.Fields(string(output))
			if len(parts) >= 4 {
				for i := 0; i < 3; i++ {
					loadAvg[i], _ = strconv.ParseFloat(parts[i+1], 64)
				}
			}
			return loadAvg, nil
		}
		return loadAvg, err
	}

	// 解析Linux加载平均值
	parts := strings.Fields(string(loadFile))
	if len(parts) >= 3 {
		for i := 0; i < 3; i++ {
			loadAvg[i], _ = strconv.ParseFloat(parts[i], 64)
		}
	}

	return loadAvg, nil
}

// getSystemMemInfo 获取系统内存信息
func getSystemMemInfo() (map[string]uint64, error) {
	memInfo := map[string]uint64{
		"total": 0,
		"used":  0,
		"free":  0,
	}

	if runtime.GOOS == "linux" {
		// 读取 /proc/meminfo
		content, err := os.ReadFile("/proc/meminfo")
		if err != nil {
			return memInfo, err
		}

		lines := strings.Split(string(content), "\n")
		memTotal := uint64(0)
		memFree := uint64(0)
		memAvailable := uint64(0)

		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
					memTotal *= 1024 // Convert from KB to bytes
				}
			} else if strings.HasPrefix(line, "MemFree:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					memFree, _ = strconv.ParseUint(fields[1], 10, 64)
					memFree *= 1024 // Convert from KB to bytes
				}
			} else if strings.HasPrefix(line, "MemAvailable:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					memAvailable, _ = strconv.ParseUint(fields[1], 10, 64)
					memAvailable *= 1024 // Convert from KB to bytes
				}
			}
		}

		memInfo["total"] = memTotal
		memInfo["free"] = memAvailable
		memInfo["used"] = memTotal - memAvailable

	} else if runtime.GOOS == "darwin" {
		// 在MacOS上使用sysctl命令
		cmd := exec.Command("sysctl", "-n", "hw.memsize")
		output, err := cmd.Output()
		if err != nil {
			return memInfo, err
		}
		memTotal, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

		cmd = exec.Command("vm_stat")
		output, err = cmd.Output()
		if err != nil {
			return memInfo, err
		}

		lines := strings.Split(string(output), "\n")
		pageSize := uint64(4096) // 默认页大小
		freePages := uint64(0)

		for _, line := range lines {
			if strings.Contains(line, "Pages free:") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					pages, _ := strconv.ParseUint(strings.ReplaceAll(fields[2], ".", ""), 10, 64)
					freePages += pages
				}
			}
		}

		memFree := freePages * pageSize

		memInfo["total"] = memTotal
		memInfo["free"] = memFree
		memInfo["used"] = memTotal - memFree
	}

	return memInfo, nil
}

// getCPUInfo 获取CPU信息
func getCPUInfo() models.SystemCPU {
	info := models.SystemCPU{
		Cores:        runtime.NumCPU(),
		UsagePercent: 0.0,
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("top", "-bn1")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Cpu(s)") {
					fields := strings.Fields(line)
					for i, field := range fields {
						if strings.Contains(field, "id,") && i > 0 {
							idleStr := strings.Replace(fields[i-1], ",", ".", -1)
							idle, err := strconv.ParseFloat(idleStr, 64)
							if err == nil {
								info.UsagePercent = 100.0 - idle
							}
							break
						}
					}
					break
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		// MacOS上获取CPU使用率
		cmd := exec.Command("top", "-l", "1", "-n", "0")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "CPU usage") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						usageParts := strings.Split(parts[1], ",")
						for _, part := range usageParts {
							if strings.Contains(part, "idle") {
								fields := strings.Fields(part)
								if len(fields) > 0 {
									idleStr := strings.Replace(fields[0], "%", "", -1)
									idle, err := strconv.ParseFloat(idleStr, 64)
									if err == nil {
										info.UsagePercent = 100.0 - idle
									}
								}
							}
						}
					}
					break
				}
			}
		}
	}

	return info
}

// getDiskInfo 获取磁盘使用情况
func getDiskInfo() models.SystemDisk {
	info := models.SystemDisk{
		Total:          uint64(0),
		Used:           uint64(0),
		Free:           uint64(0),
		FormattedTotal: "0 B",
		FormattedUsed:  "0 B",
		FormattedFree:  "0 B",
		UsagePercent:   0.0,
	}

	pwd, err := os.Getwd()
	if err != nil {
		return info
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// 使用df命令获取磁盘信息
		cmd := exec.Command("df", "-k", pwd)
		output, err := cmd.Output()
		if err != nil {
			return info
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 6 {
				total, _ := strconv.ParseUint(fields[1], 10, 64)
				used, _ := strconv.ParseUint(fields[2], 10, 64)
				free, _ := strconv.ParseUint(fields[3], 10, 64)

				total *= 1024
				used *= 1024
				free *= 1024

				info.Total = total
				info.Used = used
				info.Free = free
				info.FormattedTotal = utils.FormatBytes(int64(total))
				info.FormattedUsed = utils.FormatBytes(int64(used))
				info.FormattedFree = utils.FormatBytes(int64(free))

				usageStr := strings.Replace(fields[4], "%", "", -1)
				usage, err := strconv.ParseFloat(usageStr, 64)
				if err == nil {
					info.UsagePercent = usage
				}
			}
		}
	}

	return info
}

// getSystemUptime 获取系统运行时间
func getSystemUptime() models.SystemUptime {
	info := models.SystemUptime{
		Seconds: int64(0),
		Text:    "未知",
	}

	if runtime.GOOS == "linux" {
		// 读取/proc/uptime文件
		content, err := os.ReadFile("/proc/uptime")
		if err == nil {
			fields := strings.Fields(string(content))
			if len(fields) > 0 {
				uptime, err := strconv.ParseFloat(fields[0], 64)
				if err == nil {
					info.Seconds = int64(uptime)
					info.Text = utils.FormatDuration(int64(uptime))
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		// 使用sysctl命令获取启动时间
		cmd := exec.Command("sysctl", "-n", "kern.boottime")
		output, err := cmd.Output()
		if err == nil {
			// 解析输出，格式通常是 "{ sec = 1234567890, usec = 123456 } Thu Jan 1 00:00:00 2020"
			parts := strings.Split(string(output), "sec = ")
			if len(parts) > 1 {
				secParts := strings.Split(parts[1], ",")
				if len(secParts) > 0 {
					bootTime, err := strconv.ParseInt(strings.TrimSpace(secParts[0]), 10, 64)
					if err == nil {
						now := time.Now().Unix()
						uptime := now - bootTime
						info.Seconds = uptime
						info.Text = utils.FormatDuration(uptime)
					}
				}
			}
		}
	}

	return info
}

// getBasicSystemInfo 获取系统基本信息
func getBasicSystemInfo() models.BasicSysInfo {
	info := models.BasicSysInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: "未知",
		Version:  "未知",
	}

	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	if runtime.GOOS == "linux" {
		// 尝试获取Linux发行版信息
		if content, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					versionStr := strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
					info.Version = versionStr
					break
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sw_vers", "-productVersion")
		if output, err := cmd.Output(); err == nil {
			info.Version = "macOS " + strings.TrimSpace(string(output))
		}
	}

	return info
}

// FormatDuration 将秒数格式化为人类可读的持续时间
func FormatDuration(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, secs)
	}
	return fmt.Sprintf("%d秒", secs)
}

// FormatTimeAgo 将时间格式化为"多久之前"的形式
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	months := days / 30
	years := days / 365

	switch {
	case years > 0:
		return fmt.Sprintf("%d年前", years)
	case months > 0:
		return fmt.Sprintf("%d个月前", months)
	case days > 0:
		return fmt.Sprintf("%d天前", days)
	case hours > 0:
		return fmt.Sprintf("%d小时前", hours)
	case minutes > 0:
		return fmt.Sprintf("%d分钟前", minutes)
	default:
		return "刚刚"
	}
}
