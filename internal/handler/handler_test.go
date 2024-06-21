package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/rostis232/prmv/models"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) AddPost(post models.Post) (models.Post, error) {
	args := m.Called(post)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *MockService) GetAllPosts() ([]models.Post, error) {
	args := m.Called()
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockService) UpdatePost(post models.Post) (models.Post, error) {
	args := m.Called(post)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *MockService) GetPost(id int) (models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *MockService) DeletePost(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAddPost(t *testing.T) {
	mockService := new(MockService)
	h := NewHandler(mockService)
	e := echo.New()
	reqBody := `{"title":"Test Post","content":"Test Content"}`

	mockService.On("AddPost", mock.Anything).Return(models.Post{ID: 1, Title: "Test Post", Content: "Test Content"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.AddPost(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := models.Post{}

	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, "Test Post", resp.Title)
	assert.Equal(t, "Test Content", resp.Content)

	mockService.AssertExpectations(t)
}

func TestGetAllPosts(t *testing.T) {
	mockService := new(MockService)
	h := NewHandler(mockService)
	e := echo.New()

	mockPosts := []models.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, Title: "Post 2", Content: "Content 2"},
	}
	mockService.On("GetAllPosts").Return(mockPosts, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetAllPosts(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := []models.Post{}
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		assert.NoError(t, err)
	}

	for i, post := range resp {
		assert.Equal(t, mockPosts[i].Title, post.Title)
		assert.Equal(t, mockPosts[i].Content, post.Content)
	}

	mockService.AssertExpectations(t)
}

func TestUpdatePost(t *testing.T) {
	mockService := new(MockService)
	h := NewHandler(mockService)
	e := echo.New()
	reqBody := `{"title":"Updated Post","content":"Updated Content"}`

	mockService.On("UpdatePost", mock.Anything).Return(models.Post{ID: 1, Title: "Updated Post", Content: "Updated Content"}, nil)

	req := httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.UpdatePost(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := models.Post{}

	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, "Updated Post", resp.Title)
	assert.Equal(t, "Updated Content", resp.Content)

	mockService.AssertExpectations(t)
}

func TestGetPost(t *testing.T) {
	mockService := new(MockService)
	h := NewHandler(mockService)
	e := echo.New()

	mockService.On("GetPost", 1).Return(models.Post{ID: 1, Title: "Test Post", Content: "Test Content"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.GetPost(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := models.Post{}

	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, "Test Post", resp.Title)
	assert.Equal(t, "Test Content", resp.Content)

	mockService.AssertExpectations(t)
}

func TestDeletePost(t *testing.T) {
	mockService := new(MockService)
	h := NewHandler(mockService)
	e := echo.New()

	mockService.On("DeletePost", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.DeletePost(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	mockService.AssertExpectations(t)
}
