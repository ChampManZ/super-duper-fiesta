package models

import (
	"time"
)

// User represents a user in the system
// @Description Represents a user with associated posts and comments
type User struct {
	UserID      uint   `gorm:"primaryKey" json:"uid"`
	Username    string `gorm:"unique;not null" json:"username" validate:"required,min=3,max=32"`
	Firstname   string `gorm:"not null" json:"firstname" validate:"required"`
	Surname     string `gorm:"not null" json:"surname" validate:"required"`
	Email       string `gorm:"unique;not null" validate:"required,email"`
	Password    string `gorm:"not null" json:"-" validate:"required,min=8"`
	IsAdmin     string `gorm:"default:'0'" json:"is_admin"`
	CookieToken string `json:"-"`
	Posts       []Post
	Comments    []CommentUser
}

// Post represents a post in the system
// @Description Represents a post created by a user
type Post struct {
	PostID    uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	Message   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Comment represents a comment in the system
// @Description Represents a comment made by users
type Comment struct {
	CommentID  uint   `gorm:"primaryKey"`
	PostID     uint   `gorm:"not null"`
	CommentMSG string `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CommentUser represents the relationship between comments and users
// @Description Represents the many-to-many relationship between comments and users
type CommentUser struct {
	ID        uint `gorm:"primaryKey"`
	CommentID uint
	UserID    uint
	Comment   Comment `gorm:"constraint:OnDelete:CASCADE"`
	User      User    `gorm:"constraint:OnDelete:CASCADE"`
}
