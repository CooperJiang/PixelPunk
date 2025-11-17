package middleware

import (
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/auth"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextPayloadKey = "payload"    // JWT解析出的claims
	AuthErrorKey      = "auth_error" // 认证错误信息的键
)

/* checkUserActive 检查用户是否处于活跃状态（未被禁用）
 * 优先使用Redis黑名单，Redis不可用时降级为数据库查询
 * 返回true表示用户正常，false表示用户已禁用
 */
func checkUserActive(claims *auth.JWTClaims) bool {
	if claims == nil {
		return true
	}

	// 优先使用Redis黑名单检查（性能最优，<1ms）
	redisClient := cache.GetRedisClient()
	if redisClient != nil {
		ctx := cache.GetRedisContext()
		ns := cache.GetNamespace()
		key := ns + ":disabled_users"

		isDisabled, err := redisClient.SIsMember(ctx, key, claims.UserID).Result()
		if err != nil {
			logger.Warn("checkUserActive: Redis查询失败，降级DB: %v", err)
		} else if isDisabled {
			return false
		} else {
			return true
		}
	}

	// Redis不可用或查询失败，降级为数据库查询
	db := database.GetDB()
	var user models.User
	if err := db.Select("status").First(&user, claims.UserID).Error; err != nil {
		logger.Warn("checkUserActive: 用户状态查询失败: userID=%d, error=%v", claims.UserID, err)
		return true // 查询失败默认放行
	}

	if user.Status != common.UserStatusNormal {
		return false
	}

	return true
}

func GetCurrentUser(c *gin.Context) *auth.JWTClaims {
	value, _ := c.Get(ContextPayloadKey)
	claims, _ := value.(*auth.JWTClaims)
	return claims
}

func GetCurrentUserID(c *gin.Context) uint {
	claims := GetCurrentUser(c)
	if claims == nil {
		return 0
	}
	return claims.UserID
}

func GetCurrentUsername(c *gin.Context) string {
	claims := GetCurrentUser(c)
	if claims == nil {
		return ""
	}
	return claims.Username
}

func GetCurrentUserRole(c *gin.Context) int {
	claims := GetCurrentUser(c)
	if claims == nil {
		return 0
	}
	return claims.Role
}

func IsCurrentUserAdmin(c *gin.Context) bool {
	role := GetCurrentUserRole(c)
	return role == common.UserRoleAdmin || role == common.UserRoleSuperAdmin
}

func IsCurrentUserSuperAdmin(c *gin.Context) bool {
	role := GetCurrentUserRole(c)
	return role == common.UserRoleSuperAdmin
}

func CanUserAccessProtectedFile(c *gin.Context, imageUserID uint) bool {
	userID := GetCurrentUserID(c)

	if userID == 0 {
		return false
	}

	if userID == imageUserID {
		return true
	}

	return IsCurrentUserAdmin(c)
}

/* JWTAuth JWT解析中间件（验证token有效性和过期时间） */
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		jwtSecret := ""

		securitySettings, err := setting.GetSettingsByGroupAsMap("security")
		if err == nil {
			if val, ok := securitySettings.Settings["jwt_secret"]; ok {
				if secretStr, ok := val.(string); ok {
					jwtSecret = secretStr
				}
			}
		}
		if strings.TrimSpace(jwtSecret) == "" {
			c.Set(AuthErrorKey, "系统配置错误：JWT密钥未设置")
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString = c.Query("token")
			if tokenString == "" {
				tokenString, _ = c.Cookie("token")
			}
		} else {
			const prefix = "Bearer "
			if len(tokenString) > len(prefix) && tokenString[:len(prefix)] == prefix {
				tokenString = tokenString[len(prefix):]
			}
		}

		if tokenString == "" {
			c.Set(AuthErrorKey, "未登录")
			c.Next()
			return
		}

		claims, err := auth.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.Set(AuthErrorKey, "无效的认证令牌")
			c.Next()
			return
		}

		if claims.ExpiresAt.Unix() < auth.GetCurrentTimestamp() {
			c.Set(AuthErrorKey, "认证令牌已过期")
			c.Next()
			return
		}

		c.Set(ContextPayloadKey, claims)

		// 检查用户是否被禁用（在JWT解析后立即检查，覆盖所有需要认证的接口）
		if !checkUserActive(claims) {
			// 使用专用错误码，前端会跳转到refuse页面
			errors.HandleError(c, errors.New(errors.CodeUserDisabled, "用户账号已被禁用"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authError, exists := c.Get(AuthErrorKey); exists {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, authError.(string)))
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authError, exists := c.Get(AuthErrorKey); exists {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, authError.(string)))
			c.Abort()
			return
		}

		claims := GetCurrentUser(c)
		if claims == nil {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户认证信息无效"))
			c.Abort()
			return
		}

		if claims.Role != common.UserRoleAdmin && claims.Role != common.UserRoleSuperAdmin {
			errors.HandleError(c, errors.New(errors.CodeForbidden, "需要管理员权限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authError, exists := c.Get(AuthErrorKey); exists {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, authError.(string)))
			c.Abort()
			return
		}

		claims := GetCurrentUser(c)
		if claims == nil {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户认证信息无效"))
			c.Abort()
			return
		}

		if claims.Role != common.UserRoleSuperAdmin {
			errors.HandleError(c, errors.New(errors.CodeForbidden, "需要超级管理员权限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

/* RequireActiveUser 用户状态校验中间件（检查用户是否被禁用）
 * 使用Redis黑名单实现即时踢出，Redis不可用时降级为数据库查询
 * 建议用于：写操作、敏感接口、管理后台
 * 不建议用于：高频只读接口（如图片访问）
 */
func RequireActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetCurrentUser(c)
		if claims == nil {
			// 没有用户信息，跳过检查（由RequireAuth处理）
			c.Next()
			return
		}

		// 优先使用Redis黑名单检查（性能最优，<1ms）
		redisClient := cache.GetRedisClient()
		if redisClient != nil {
			ctx := cache.GetRedisContext()
			ns := cache.GetNamespace()
			key := ns + ":disabled_users"

			isDisabled, err := redisClient.SIsMember(ctx, key, claims.UserID).Result()
			if err != nil {
				logger.Warn("Redis黑名单查询失败，降级为数据库查询: %v", err)
			} else if isDisabled {
				errors.HandleError(c, errors.New(errors.CodeForbidden, "用户账号已被禁用"))
				c.Abort()
				return
			} else {
				// Redis检查通过，直接放行
				c.Next()
				return
			}
		}

		// Redis不可用或查询失败，降级为数据库查询
		db := database.GetDB()
		var user models.User
		if err := db.Select("status").First(&user, claims.UserID).Error; err != nil {
			logger.Warn("用户状态查询失败: userID=%d, error=%v", claims.UserID, err)
			c.Next()
			return
		}

		if user.Status != common.UserStatusNormal {
			var statusMsg string
			switch user.Status {
			case common.UserStatusDisabled:
				statusMsg = "用户账号已被禁用"
			case common.UserStatusDeleted:
				statusMsg = "用户账号已被删除"
			default:
				statusMsg = "用户账号状态异常"
			}
			errors.HandleError(c, errors.New(errors.CodeForbidden, statusMsg))
			c.Abort()
			return
		}

		c.Next()
	}
}
