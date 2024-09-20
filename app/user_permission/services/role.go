package services

import (
	"context"
	"jade-mes/app/user_permission/model"
	"jade-mes/app/user_permission/repositories"
)

func CreateRole(ctx context.Context, id int64, name, desc string) (error) {
	role := &model.Role{
		ID: id,
		RoleName: name,
		RoleDesc: desc,
	}
	return repositories.CreateRole(ctx, role)
}

func FindRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	return repositories.FindRoleByID(ctx, id)
}

func FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	return repositories.FindRoleByName(ctx, name)
}