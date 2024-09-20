// 用户仓库的具体实现
package repositories

import (
	"context"
	"jade-mes/app/infrastructure/persistence/database"
	"jade-mes/app/user_permission/model"
)

// CreateUser
func CreateUser(ctx context.Context, user *model.User) error {
	return database.GetDB().Table("users").Create(user).Error
}

// FindUserByID 实现根据ID查找用户
func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	err := database.GetDB().First(user, " id = ?", id).Error
	return user, err
}

// FindByUsername 实现根据用户名查找用户
func FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := database.GetDB().Find(user, "username = ?", username).Error
	return user, err
}

// 其他方法实现，例如 Update 或 Delete，可以按照类似模式
