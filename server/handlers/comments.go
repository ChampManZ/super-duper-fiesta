package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetComments godoc
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {array} models.Comment
// @Failure 500 {object} map[string]string "Failed to retrieve comments"
// @Router /api/v1/admin/comments [get]
func GetComments(c echo.Context) error {
	postID := c.Param("pid")

	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Post does not exist"})
	}

	var comments []models.GetCommentRequest
	if err := config.DB.Table("comments").Select("users.username, comments.comment_msg").Joins("inner join comment_users on comment_users.comment_id = comments.comment_id").Joins("inner join users on users.user_id = comment_users.user_id").Where("comments.post_id = ?", postID).Scan(&comments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get comments"})
	}

	return c.JSON(http.StatusOK, comments)
}

// PostComment godoc
// @Summary Create a comment
// @Description Create a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment object that needs to be created"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "Invalid input or failed to create comment"
// @Router /api/v1/restricted/comments [post]
func CreateComment(c echo.Context) error {
	request := new(models.CreateCommentRequest)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	var post models.Post
	if result := config.DB.First(&post, request.PostID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Post does not exist"})
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(*models.JWTClaims)

	userID := claims.UserID

	comment := models.Comment{
		PostID:     request.PostID,
		CommentMSG: request.CommentMSG,
	}

	if result := config.DB.Create(&comment); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create comment"})
	}

	commentUser := models.CommentUser{
		CommentID: comment.CommentID,
		UserID:    userID,
	}

	if result := config.DB.Create(&commentUser); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create comment and user data"})
	}

	return c.JSON(http.StatusCreated, comment)
}
