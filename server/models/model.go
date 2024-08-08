package models

import "time"

type User struct {
	UserID    uint   `gorm:"primarykey"`
	Username  string `gorm:"unique"`
	Firstname string `gorm:"not null"`
	Surname   string `gorm:"not null"`
	Email     string `gorm:"unique"`
	Password  string `gorm:"not null"`
	Posts     []Post // One-to-Many relationship (has many) | One user can have many posts
}

type Post struct {
	PostID    uint `gorm:"primarykey"`
	UserID    uint
	User      User
	Message   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
