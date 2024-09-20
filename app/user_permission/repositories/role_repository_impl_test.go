package repositories

import (
	"context"
	"testing"

	"jade-mes/app/user_permission/model"
)

func TestFindByName(t *testing.T) {
	// Create a new role repository
	repo := NewRoleRepository()

	// Define the test cases
	tests := []struct {
		name        string
		roleName    string
		expectedRole *model.Role
	}{
		{
			name:     "Role found",
			roleName: "admin",
			expectedRole: &model.Role{
				ID:       1,
				RoleName: "admin",
				RoleDesc: "Administrator role",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, _ := repo.FindByName(context.Background(), tt.roleName)
			t.Logf("Role: %v", role)
		})
	}
}