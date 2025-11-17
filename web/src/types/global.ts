/**
 * 全局类型定义 - 统一导出文件
 * ⚠️ 此文件已重构为模块化结构，保留此文件是为了向后兼容
 * 新代码请直接从以下模块导入：
 * - @/types/base      - 基础类型
 * - @/types/api       - API相关
 * - @/types/business  - 业务领域
 * - @/types/ui        - UI组件
 * - @/types/form      - 表单相关
 * - @/types/notification - 通知相关
 * - @/types/event     - 事件相关
 * - @/types/performance - 性能相关
 * - @/types/utility   - 工具类型
 * 或统一从 @/types 导入所有类型

/* 重新导出所有模块化类型 */
export * from './base'
export * from './api'
export * from './business'
export * from './ui'
export * from './form'
export * from './notification'
export * from './event'
export * from './performance'
export * from './utility'

/* 环境变量声明已移至 env.d.ts，会自动加载 */
