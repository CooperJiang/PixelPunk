/**
 * 活动动作枚举
 * 与后端 internal/constants/activity_action.go 保持一致
 */
export enum ActivityAction {
  CREATE = 1, // 创建
  DELETE = 2, // 删除
  UPDATE = 3, // 更新
  ENABLE = 4, // 启用
  DISABLE = 5, // 禁用
  MOVE = 6, // 移动
  RENAME = 7, // 重命名
  UPLOAD = 8, // 上传
  LOGIN = 9, // 登录
  REGISTER = 10, // 注册
  CHANGE = 11, // 修改（通用）
  EXPIRE = 12, // 过期
  CLEANUP = 13, // 清理
  REGENERATE = 14, // 重新生成
  TOGGLE = 15, // 切换状态（通用）
}

/**
 * 访问级别枚举
 */
export enum AccessLevel {
  PRIVATE = 'private', // 私有
  PUBLIC = 'public', // 公开
  PROTECTED = 'protected', // 受保护
  UNLISTED = 'unlisted', // 不公开列表
}

/**
 * 帮助函数：获取动作名称（调试用）
 */
export function getActionName(action: ActivityAction): string {
  const actionNames: Record<ActivityAction, string> = {
    [ActivityAction.CREATE]: 'create',
    [ActivityAction.DELETE]: 'delete',
    [ActivityAction.UPDATE]: 'update',
    [ActivityAction.ENABLE]: 'enable',
    [ActivityAction.DISABLE]: 'disable',
    [ActivityAction.MOVE]: 'move',
    [ActivityAction.RENAME]: 'rename',
    [ActivityAction.UPLOAD]: 'upload',
    [ActivityAction.LOGIN]: 'login',
    [ActivityAction.REGISTER]: 'register',
    [ActivityAction.CHANGE]: 'change',
    [ActivityAction.EXPIRE]: 'expire',
    [ActivityAction.CLEANUP]: 'cleanup',
    [ActivityAction.REGENERATE]: 'regenerate',
    [ActivityAction.TOGGLE]: 'toggle',
  }
  return actionNames[action] || 'unknown'
}
