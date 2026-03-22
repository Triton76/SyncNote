package util

import (
	"math/rand"
	"time"
)

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	adjs  = []string{"happy", "clever", "brave", "swift", "calm", "bright", "gentle", "wise"}
	nouns = []string{"panda", "tiger", "eagle", "wolf", "fox", "bear", "lion", "deer"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 方案 1: 纯随机字符
func GenerateRandomUsername(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// 方案 2: 形容词 + 名词 + 数字（可读性好）✅ 推荐
func GenerateReadableUsername() string {
	adj := adjs[rand.Intn(len(adjs))]
	noun := nouns[rand.Intn(len(nouns))]
	num := rand.Intn(9999)
	return adj + noun + string(rune(num))
}

// 方案 3: UUID 前缀
func GenerateUUIDUsername() string {
	return "user_" + GenerateRandomUsername(8)
}
