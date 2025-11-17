package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/* JWTClaims 自定义JWT声明结构 */
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Role     int    `json:"role"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const (
	defaultJWTSecret    = "defaultSecretKey" // 默认密钥（只在无法获取设置时使用）
	defaultExpiresHours = 24                 // 默认24小时
)

/* GetCurrentTimestamp 获取当前时间戳 */
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

/* GenerateToken 生成JWT令牌 */
func GenerateToken(userID uint, username string, role int, jwtSecret string, expiresHours int) (string, error) {
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
	}

	if expiresHours <= 0 {
		expiresHours = defaultExpiresHours
	}

	expirationTime := time.Now().Add(time.Duration(expiresHours) * time.Hour)

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

/* ParseToken 解析JWT令牌 */
func ParseToken(tokenString string, jwtSecret string) (*JWTClaims, error) {
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

/* VerifyTokenValid 验证token是否有效且未过期，返回claims */
func VerifyTokenValid(tokenString string, jwtSecret string) (*JWTClaims, error) {
	claims, err := ParseToken(tokenString, jwtSecret)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt.Unix() < GetCurrentTimestamp() {
		return nil, fmt.Errorf("token已过期")
	}

	return claims, nil
}
