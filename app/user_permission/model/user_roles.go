package model

// CREATE TABLE `user_roles` (
// 	`user_id` int(11) NOT NULL,
// 	`role_id` int(11) NOT NULL,
// 	PRIMARY KEY (`user_id`,`role_id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

// UserRole 用户角色关联表
type UserRole struct {
	UserID int `gorm:"column:user_id"` // 用户ID
	RoleID int `gorm:"column:role_id"` // 角色ID
}
