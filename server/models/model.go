package models

import "time"

type User struct {
	UserID    uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Firstname string    `gorm:"not null"`
	Surname   string    `gorm:"not null"`
	Email     string    `gorm:"unique"`
	Password  string    `gorm:"not null"`
	Posts     []Post    `gorm:"foreignKey:UserID"`       // One-to-Many relationship (has many) | One user can have many posts
	Comments  []Comment `gorm:"many2many:comment_users"` // Many-to-Many relationship (has many) | One user can have many comments
}

type Post struct {
	PostID    uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID"` // One-to-One relationship (belongs to) | One post can have one user
}

type Comment struct {
	CommentID  uint      `gorm:"primaryKey"`
	CommentMSG string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoCreateTime"`
	Users      []User    `gorm:"many2many:comment_users"` // Many-to-Many relationship (belongs to) | Many users can have many comments which lead to CommentUser table
}

type CommentUser struct {
	CommentID uint `gorm:"primaryKey;autoIncrement:false"` // Composite primary key
	UserID    uint `gorm:"primaryKey;autoIncrement:false"` // Composite primary key
}
