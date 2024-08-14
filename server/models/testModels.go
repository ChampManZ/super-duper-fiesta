package models

import "time"

// For testing only, the main model have cascading references while this focus on unit test
type MockPost struct {
	PostID    uint      `gorm:"column:post_id;primaryKey;autoIncrement:true" json:"post_id"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	Message   string    `gorm:"column:message;not null" json:"message"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type MockUser struct {
	UserID      uint   `gorm:"column:user_id;primaryKey;autoIncrement:true" json:"user_id"`
	Username    string `gorm:"column:username;not null" json:"username"`
	Firstname   string `gorm:"column:firstname;not null" json:"firstname"`
	Surname     string `gorm:"column:surname;not null" json:"surname"`
	Email       string `gorm:"column:email;not null" json:"email"`
	Password    string `gorm:"column:password;not null" json:"password"`
	IsAdmin     string `gorm:"column:is_admin" json:"is_admin"`
	CookieToken string `gorm:"column:cookie_token" json:"cookie_token"`
}

type MockCommentUser struct {
	CommentID uint `gorm:"column:comment_id;primaryKey" json:"comment_id"`
	UserID    uint `gorm:"column:user_id;primaryKey" json:"user_id"`
}

type MockComment struct {
	CommentID uint      `gorm:"column:comment_id;primaryKey;autoIncrement:true" json:"comment_id"`
	PostID    uint      `gorm:"column:post_id;not null" json:"post_id"`
	Message   string    `gorm:"column:comment_msg;not null" json:"comment_msg"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
