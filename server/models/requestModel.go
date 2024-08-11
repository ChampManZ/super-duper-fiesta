package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// Request Model
// Question: Does JWT Claims consider as a request model in terms of dealing file for project structure?
type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Firstname string `json:"firstname" validate:"required"`
	Surname   string `json:"surname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Firstname string `json:"firstname" validate:"required"`
	Surname   string `json:"surname" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,min=8"`
}

type LoginUserRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type JWTClaims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	Admin    string `json:"admin"`
	jwt.RegisteredClaims
}
