/**
 * 性能相关类型定义
 * 包含性能监控、指标统计等类型

/* 性能指标 */
export interface PerformanceMetrics {
  load_time: number
  render_time: number
  memory_usage: number
  cpu_usage?: number
  network_latency?: number
  error_rate?: number
}

/* 性能监控配置 */
export interface PerformanceConfig {
  enable_monitoring: boolean
  sample_rate: number
  report_interval: number
  thresholds: {
    load_time: number
    memory_usage: number
    error_rate: number
  }
}
