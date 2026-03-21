package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对明文密码进行 bcrypt 加密
// 输入：用户输入的原始密码
// 输出：加密后的字符串 (存入数据库)，错误信息
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 通常是 10，安全性与性能的平衡
	// 如果需要更高安全性，可以增加 cost (范围 4-31)，但会变慢
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证用户输入的密码是否与数据库中的哈希值匹配
// 输入：用户输入的明文密码，数据库中取出的哈希密码
// 输出：bool (true 表示匹配，false 表示不匹配)
func CheckPassword(password, hash string) bool {
	// CompareHashAndPassword 如果匹配则返回 nil，否则返回 error
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
