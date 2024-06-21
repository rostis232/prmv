package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rostis232/prmv/models"
	"net/http"
	"strconv"
)

type Handler struct {
	Service Service
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
		Service: service,
	}
}

func (h Handler) AddPost(c echo.Context) error {
	var post models.Post
	err := c.Bind(&post)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	newPost, err := h.Service.AddPost(post)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, newPost)
}

func (h Handler) GetAllPosts(c echo.Context) error {
	posts, err := h.Service.GetAllPosts()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, posts)
}

func (h Handler) UpdatePost(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	var post models.Post
	err = c.Bind(&post)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	post.ID = idInt
	updatedPost, err := h.Service.UpdatePost(post)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, updatedPost)
}

func (h Handler) GetPost(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	post, err := h.Service.GetPost(idInt)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, post)
}

func (h Handler) DeletePost(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.Service.DeletePost(idInt)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
