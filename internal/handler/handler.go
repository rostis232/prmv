package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rostis232/prmv/models"
	"net/http"
	"strconv"
)

type Handler struct {
	Service  Service
	validate *validator.Validate
}

type Service interface {
	AddPost(models.Post) (models.Post, error)
	GetAllPosts() ([]models.Post, error)
	UpdatePost(models.Post) (models.Post, error)
	GetPost(id int) (models.Post, error)
	DeletePost(id int) error
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service:  service,
		validate: validator.New(),
	}
}

type postData struct {
	Title   string `db:"title" json:"title" validate:"required,min=3,max=100"`
	Content string `db:"content" json:"content" validate:"required,min=3"`
}

// AddPost godoc
// @Summary Add a new post
// @Description Add a new post with the input payload
// @Tags posts
// @Accept  json
// @Produce  json
// @Param post body postData true "Post Data"
// @Success 200 {object} models.Post
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts [post]
func (h *Handler) AddPost(c echo.Context) error {
	var post postData

	err := c.Bind(&post)
	if err != nil {
		return newErrorResponse(c, http.StatusBadRequest, "invalid post data")
	}

	err = h.validate.Struct(post)
	if err != nil {
		return newErrorResponse(c, http.StatusBadRequest, "invalid post data")
	}

	newPost, err := h.Service.AddPost(models.Post{
		Title:   post.Title,
		Content: post.Content,
	})
	if err != nil {
		log.Errorf("error adding post: %v", err)
		return newErrorResponse(c, http.StatusInternalServerError, "error adding post")
	}

	return c.JSON(http.StatusCreated, newPost)
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Get a list of all posts
// @Tags posts
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Post
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func (h *Handler) GetAllPosts(c echo.Context) error {
	posts, err := h.Service.GetAllPosts()
	if err != nil {
		log.Errorf("error getting all posts: %v", err)
		return newErrorResponse(c, http.StatusInternalServerError, "error getting all posts")
	}

	return c.JSON(http.StatusOK, posts)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post with the given id
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Param post body postData true "Post Data"
// @Success 200 {object} models.Post
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [put]
func (h *Handler) UpdatePost(c echo.Context) error {
	idStr := c.Param("id")

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("error converting id to int: %v", err)
		return newErrorResponse(c, http.StatusBadRequest, "invalid post id")
	}

	if idInt < 1 {
		return newErrorResponse(c, http.StatusBadRequest, "invalid post id")
	}

	var post postData

	err = c.Bind(&post)
	if err != nil {
		log.Errorf("error unmarshalling post: %v", err)
		return newErrorResponse(c, http.StatusBadRequest, "invalid post data")
	}

	if post.Title == "" && post.Content == "" {
		return newErrorResponse(c, http.StatusBadRequest, "invalid post data")
	}

	updatedPost, err := h.Service.UpdatePost(models.Post{
		ID:      idInt,
		Title:   post.Title,
		Content: post.Content,
	})
	if err != nil {
		log.Errorf("error updating post: %v", err)
		return newErrorResponse(c, http.StatusInternalServerError, "error updating post")
	}

	return c.JSON(http.StatusOK, updatedPost)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get a single post by its ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [get]
func (h *Handler) GetPost(c echo.Context) error {
	idStr := c.Param("id")

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("error converting id to int: %v", err)
		return newErrorResponse(c, http.StatusBadRequest, "invalid post id")
	}

	if idInt < 1 {
		return newErrorResponse(c, http.StatusBadRequest, "invalid post id")
	}

	post, err := h.Service.GetPost(idInt)
	if err != nil {
		log.Errorf("error getting post: %v", err)
		return newErrorResponse(c, http.StatusInternalServerError, "error getting post")
	}

	return c.JSON(http.StatusOK, post)
}

// DeletePost godoc
// @Summary Delete a post by ID
// @Description Delete a single post by its ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [delete]
func (h *Handler) DeletePost(c echo.Context) error {
	idStr := c.Param("id")

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("error converting id to int: %v", err)
		return newErrorResponse(c, http.StatusBadRequest, "invalid post id")
	}

	if idInt < 1 {
		return newErrorResponse(c, http.StatusBadRequest, "invalid post id")
	}

	err = h.Service.DeletePost(idInt)
	if err != nil {
		log.Errorf("error deleting post: %v", err)
		return newErrorResponse(c, http.StatusInternalServerError, "error deleting post")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Home(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
}
