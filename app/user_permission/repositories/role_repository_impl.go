// 角色仓库的具体实现
package repositories

import (
	"context"
	"jade-mes/app/infrastructure/persistence/database"
	"jade-mes/app/user_permission/model"
)

// CreateRole
func CreateRole(ctx context.Context, role *model.Role) error {
	return database.GetDB().Table("roles").Create(role).Error
}

// FindRoleByID 实现根据ID查找角色
func FindRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	role := &model.Role{}
	err := database.GetDB().First(role, "id = ?", id).Error
	return role, err
}

// FindRoleByName 实现根据名称查找角色
func FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	role := &model.Role{}
	err := database.GetDB().Find(role, "role_name = ?", name).Error
	return role, err
}
