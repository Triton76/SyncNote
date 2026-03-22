package auth

import "github.com/golang-jwt/jwt/v4"

// Claims 自定义的 JWT 载荷结构
// 必须嵌入 jwt.RegisteredClaims 以包含标准字段 (exp, iat, iss 等)
type Claims struct {
	UserID string `json:"userId"` // 核心：用户唯一标识
	// 你可以在这里扩展其他字段，比如 Role, Username 等
	// Role   string `json:"role"`
	jwt.RegisteredClaims
}
