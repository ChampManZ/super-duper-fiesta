package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Request Model
// Question: Does JWT Claims consider as a request model in terms of dealing file for project structure?

// CreateUserRequest represents the data needed to create a user
// @Description Request model for creating a user
type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Firstname string `json:"firstname" validate:"required"`
	Surname   string `json:"surname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents the data needed to update user information
// @Description Request model for updating user information
type UpdateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Firstname string `json:"firstname" validate:"required"`
	Surname   string `json:"surname" validate:"required"`
}

// UpdateUserPasswordRequest represents the data needed to update a user's password
// @Description Request model for updating a user's password
type UpdateUserPasswordRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserRequest represents the data needed for user login
// @Description Request model for user login
type LoginUserRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

// JWTClaims represents the JWT claims
// @Description JWT claims used for authentication and authorization
type JWTClaims struct {
	UserID    uint   `json:"uid"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Admin     string `json:"admin"`
	jwt.RegisteredClaims
}

// GetPublicPostsRequest represents the data for retrieving public posts
// @Description Response model for retrieving public posts
type GetPublicPostsRequest struct {
	PostID    uint      `json:"post_id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Surname   string    `json:"surname"`
	Message   string    `json:"post_message"`
	CreatedAt time.Time `json:"post_created_at"`
	UpdatedAt time.Time `json:"post_updated_at"`
}

// GetMigrationListRequest represents the data for retrieving migration information
// @Description Response model for retrieving migration information
type GetMigrationListRequest struct {
	MigrationID string `json:"migration_id"`
	Title       string `json:"migration_title"`
}

// RunMigrationRequest represents the data needed to run a migration
// @Description Request model for running a migration
type RunMigrationRequest struct {
	MigrationID string `json:"migration_id"`
}
