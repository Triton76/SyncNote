package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// 建议：在生产环境中，这个 Key 应该放在配置文件中，并且要足够复杂
	// 可以在 init() 函数中从 config 加载，或者作为参数传入
	SecretKey      = "syncnote-secret-key-change-me-in-prod-8888"
	TokenExpireDur = time.Hour * 24 // Token 有效期：24小时
)

// GenerateToken 生成 JWT Token
// 输入：用户 ID
// 输出：token 字符串，错误信息
func GenerateToken(userID string) (string, error) {
	now := time.Now()
	expireTime := now.Add(TokenExpireDur)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),        // 签发时间
			Issuer:    "SyncNote-API",                 // 签发者
			Subject:   "user-token",                   // 主题
		},
	}

	// 创建 token，使用 HS256 算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并返回字符串
	return token.SignedString([]byte(SecretKey))
}

// ParseToken 解析 JWT Token
// 输入：token 字符串
// 输出：解析后的 Claims 对象，错误信息
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 断言 claims 类型并检查有效性
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
