package model

// CREATE TABLE `role_permissions` (
// 	`role_id` int(11) NOT NULL,
// 	`permission_id` int(11) NOT NULL,
// 	PRIMARY KEY (`role_id`,`permission_id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       int `gorm:"column:role_id"`       // 角色ID
	PermissionID int `gorm:"column:permission_id"` // 权限ID
}
