package model

// Role 角色实体
type Role struct {
	ID       int64  `gorm:"column:id"`        // 角色的唯一标识符
	RoleName string `gorm:"column:role_name"` // 角色的名称
	RoleDesc string `gorm:"column:role_desc"` // 角色的描述
}
