package models

import (
	"time"
)

// GORM Model
// In the end, I decide to go with SQL from Scratch instead of GORM AutoMigrate() Function
// In real-world projects, we have a separate datacenter for the database that use SQL script and
// store procedures to manage the database schema. This is because we want to have full control

// User represents a user in the system
// @Description Represents a user with associated posts and comments
type User struct {
	UserID      uint      `gorm:"primaryKey" json:"uid"`
	Username    string    `gorm:"unique" json:"username"`
	Firstname   string    `gorm:"not null" json:"firstname"`
	Surname     string    `gorm:"not null" json:"surname"`
	Email       string    `gorm:"unique" json:"email"`
	Password    string    `gorm:"not null" json:"-"` // Hide password from JSON response
	IsAdmin     string    `gorm:"type:VARCHAR(1);default:'0'" json:"is_admin"`
	CookieToken string    `gorm:"unique;not null" json:"-"` // Hide cookie token from JSON response
	Posts       []Post    `gorm:"foreignKey:UserID"`        // One-to-Many relationship (has many) | One user can have many posts
	Comments    []Comment `gorm:"many2many:CommentUser"`    // Many-to-Many relationship (has many) | One user can have many comments
}

// Post represents a post in the system
// @Description Represents a post created by a user
type Post struct {
	PostID    uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID;references:UserID"` // One-to-One relationship (belongs to) | One post can have one user
}

// Comment represents a comment in the system
// @Description Represents a comment made by users
type Comment struct {
	CommentID  uint      `gorm:"primaryKey"`
	CommentMSG string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoCreateTime"`
	Users      []User    `gorm:"many2many:CommentUser"` // Many-to-Many relationship (belongs to) | Many users can have many comments which lead to CommentUser table
}

// CommentUser represents the relationship between comments and users
// @Description Represents the many-to-many relationship between comments and users
type CommentUser struct {
	CommentID uint `gorm:"primaryKey;autoIncrement:false"` // Composite primary key
	UserID    uint `gorm:"primaryKey;autoIncrement:false"` // Composite primary key
}
