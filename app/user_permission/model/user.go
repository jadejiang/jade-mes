package model

import "time"

// User 是对数据库中用户表的一个映射
type User struct {
    ID           string    `gorm:"column:id"`
    Username     string    `gorm:"column:name"`
    Password     string    `gorm:"column:password_hash"`
    PasswordSalt string    `gorm:"column:password_salt"`
    Email        string    `db:"email"`
    Mobile       string    `db:"mobile"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
    // 其他字段...
}
