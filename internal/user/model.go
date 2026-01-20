// Package user 定义用户数据模型
package user

import (
	"time"

	"gorm.io/gorm"
)

// User 表示系统中的用户实体
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`                 // 用户唯一标识
	Name         string         `gorm:"not null" json:"name"`                 // 用户姓名
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`    // 用户邮箱（唯一）
	PasswordHash string         `gorm:"not null" json:"-"`                    // 密码哈希（不返回给客户端）
	Roles        []Role         `gorm:"many2many:user_roles;" json:"-"`       // 用户角色列表（多对多关系）
	CreatedAt    time.Time      `json:"created_at"`                           // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                           // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                       // 软删除时间
}

// TableName 指定用户模型对应的数据库表名
func (User) TableName() string {
	return "users"
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// GetRoleNames 返回用户的所有角色名称列表
func (u *User) GetRoleNames() []string {
	roleNames := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roleNames[i] = role.Name
	}
	return roleNames
}
