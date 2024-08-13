package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetPosts godoc
// @Summary Retrieve all posts
// @Description Get all posts with associated user details (username, firstname, surname)
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {array} models.GetPublicPostsRequest "List of posts with user details"
// @Failure 500 {object} map[string]string "Failed to retrieve posts"
// @Router /api/v1/posts [get]
func GetPosts(c echo.Context) error {
	postID := c.QueryParam("pid")

	if postID != "" {
		postID, err := strconv.Atoi(postID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		var post models.Post
		if result := config.DB.First(&post, postID); result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Post not found"})
		}

		return c.JSON(http.StatusOK, post)
	}

	var posts []models.GetPublicPostsRequest
	if result := config.DB.Table("posts").Select("posts.post_id, users.username, users.firstname, users.surname, posts.message, posts.created_at, posts.updated_at").Joins("inner join users on users.user_id = posts.user_id").Scan(&posts); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get posts"})
	}

	return c.JSON(http.StatusOK, posts)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post by an authenticated user
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post object that needs to be created"
// @Success 201 {object} models.Post "Newly created post"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to create post"
// @Router /api/v1/restricted/posts [post]
func CreatePost(c echo.Context) error {
	request := new(models.Post)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	// Reference: https://github.com/labstack/echo/issues/1504
	user := c.Get(("user"))
	token := user.(*jwt.Token)
	claims := token.Claims.(*models.JWTClaims)

	userID := claims.UserID

	post := models.Post{
		Message: request.Message,
		UserID:  userID,
	}

	if result := config.DB.Create(&post); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create post. Please try again"})
	}

	return c.JSON(http.StatusCreated, post)
}
