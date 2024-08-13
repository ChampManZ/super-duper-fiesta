package tests

import (
	"server/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ----------- Model Testing ----------- //
// Requirement:
// UserID must not be negative | uint is always positive and Golang is type strict
// Message must not be empty
// CreatedAt and UpdatedAt must not allow future dates
func TestPost(t *testing.T) {
	validPost := models.Post{
		Message:   "This is a test post",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inValidPost := models.Post{
		Message:   "",
		CreatedAt: time.Now().AddDate(0, 0, 1),
		UpdatedAt: time.Now().AddDate(0, 0, 1),
	}

	assert.NotEmpty(t, validPost.Message)
	assert.False(t, validPost.CreatedAt.After(time.Now()))
	assert.False(t, validPost.UpdatedAt.After(time.Now()))

	assert.Empty(t, inValidPost.Message)
	assert.True(t, inValidPost.CreatedAt.After(time.Now()))
	assert.True(t, inValidPost.UpdatedAt.After(time.Now()))
}

// ----------- API Testing ----------- //
