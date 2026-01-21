// Package user 提供密码强度验证功能
package user

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"

	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

var (
	// ErrPasswordTooShort 密码太短
	ErrPasswordTooShort = errors.New("password is too short")
	// ErrPasswordMissingUppercase 密码缺少大写字母
	ErrPasswordMissingUppercase = errors.New("password must contain at least one uppercase letter")
	// ErrPasswordMissingLowercase 密码缺少小写字母
	ErrPasswordMissingLowercase = errors.New("password must contain at least one lowercase letter")
	// ErrPasswordMissingNumber 密码缺少数字
	ErrPasswordMissingNumber = errors.New("password must contain at least one digit")
	// ErrPasswordMissingSpecial 密码缺少特殊字符
	ErrPasswordMissingSpecial = errors.New("password must contain at least one special character")
)

// PasswordValidator 密码验证器
type PasswordValidator struct {
	minLength        int
	requireUppercase bool
	requireLowercase bool
	requireNumber    bool
	requireSpecial   bool
}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator(cfg *config.SecurityConfig) *PasswordValidator {
	return &PasswordValidator{
		minLength:        cfg.PasswordMinLength,
		requireUppercase: cfg.PasswordRequireUppercase,
		requireLowercase: cfg.PasswordRequireLowercase,
		requireNumber:    cfg.PasswordRequireNumber,
		requireSpecial:   cfg.PasswordRequireSpecial,
	}
}

// Validate 验证密码强度
func (v *PasswordValidator) Validate(password string) error {
	// 检查最小长度
	if len(password) < v.minLength {
		return fmt.Errorf("%w: minimum length is %d characters", ErrPasswordTooShort, v.minLength)
	}

	// 检查大写字母
	if v.requireUppercase && !containsUppercase(password) {
		return ErrPasswordMissingUppercase
	}

	// 检查小写字母
	if v.requireLowercase && !containsLowercase(password) {
		return ErrPasswordMissingLowercase
	}

	// 检查数字
	if v.requireNumber && !containsNumber(password) {
		return ErrPasswordMissingNumber
	}

	// 检查特殊字符
	if v.requireSpecial && !containsSpecial(password) {
		return ErrPasswordMissingSpecial
	}

	return nil
}

// containsUppercase 检查是否包含大写字母
func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// containsLowercase 检查是否包含小写字母
func containsLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// containsNumber 检查是否包含数字
func containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// containsSpecial 检查是否包含特殊字符
func containsSpecial(s string) bool {
	// 定义特殊字符正则表达式
	specialChars := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>[\]\\/_+=~` + "`" + `-]`)
	return specialChars.MatchString(s)
}

// GetPasswordRequirements 获取密码要求说明
func (v *PasswordValidator) GetPasswordRequirements() string {
	requirements := fmt.Sprintf("Password must be at least %d characters", v.minLength)
	
	var additional []string
	if v.requireUppercase {
		additional = append(additional, "one uppercase letter")
	}
	if v.requireLowercase {
		additional = append(additional, "one lowercase letter")
	}
	if v.requireNumber {
		additional = append(additional, "one number")
	}
	if v.requireSpecial {
		additional = append(additional, "one special character")
	}
	
	if len(additional) > 0 {
		requirements += " and contain at least "
		for i, req := range additional {
			if i > 0 {
				if i == len(additional)-1 {
					requirements += " and "
				} else {
					requirements += ", "
				}
			}
			requirements += req
		}
	}
	
	return requirements + "."
}
