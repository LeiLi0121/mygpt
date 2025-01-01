package model

import "time"

// User 用户模型
type User struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"unique;not null" json:"username" binding:"required"`
	Password  string `gorm:"not null" json:"password" binding:"required"`
	Email     string `gorm:"unique;not null" json:"email" binding:"required,email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
