package tests

import (
	"server/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Requirement:
// UserID must not be negative | uint is always positive and Golang is type strict
// CommentMSG must not be empty
// CreatedAt and UpdatedAt must not allow future dates
func TestComment(t *testing.T) {
	validComment := models.Comment{
		CommentMSG: "This is a test comment",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	assert.NotEmpty(t, validComment.CommentMSG)
	assert.False(t, validComment.CreatedAt.After(time.Now()))
	assert.False(t, validComment.UpdatedAt.After(time.Now()))
}
