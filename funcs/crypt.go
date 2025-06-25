package funcs

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// 密码哈希加密
func HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 比对哈希加密密码
// hash: 密码哈希加密后的字符串
// password: 明文密码
func CompareHashPassword(ctx context.Context, hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
