package constants

/* ActivityAction 活动动作枚举 - 用于活动日志的 action 字段 */
type ActivityAction int

const (
	ActionCreate     ActivityAction = 1  // 创建
	ActionDelete     ActivityAction = 2  // 删除
	ActionUpdate     ActivityAction = 3  // 更新
	ActionEnable     ActivityAction = 4  // 启用
	ActionDisable    ActivityAction = 5  // 禁用
	ActionMove       ActivityAction = 6  // 移动
	ActionRename     ActivityAction = 7  // 重命名
	ActionUpload     ActivityAction = 8  // 上传
	ActionLogin      ActivityAction = 9  // 登录
	ActionRegister   ActivityAction = 10 // 注册
	ActionChange     ActivityAction = 11 // 修改（通用）
	ActionExpire     ActivityAction = 12 // 过期
	ActionCleanup    ActivityAction = 13 // 清理
	ActionRegenerate ActivityAction = 14 // 重新生成
	ActionToggle     ActivityAction = 15 // 切换状态（通用）
)

/* String 转换为字符串（用于日志/调试） */
func (a ActivityAction) String() string {
	switch a {
	case ActionCreate:
		return "create"
	case ActionDelete:
		return "delete"
	case ActionUpdate:
		return "update"
	case ActionEnable:
		return "enable"
	case ActionDisable:
		return "disable"
	case ActionMove:
		return "move"
	case ActionRename:
		return "rename"
	case ActionUpload:
		return "upload"
	case ActionLogin:
		return "login"
	case ActionRegister:
		return "register"
	case ActionChange:
		return "change"
	case ActionExpire:
		return "expire"
	case ActionCleanup:
		return "cleanup"
	case ActionRegenerate:
		return "regenerate"
	case ActionToggle:
		return "toggle"
	default:
		return "unknown"
	}
}

/* Int 转换为 int（用于存储到数据库） */
func (a ActivityAction) Int() int {
	return int(a)
}

/* FromInt 从 int 创建 ActivityAction */
func FromInt(value int) ActivityAction {
	return ActivityAction(value)
}
