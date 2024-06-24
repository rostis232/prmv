package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	testCases := []struct {
		reqBody      string
		post         models.Post
		status       int
		errorExpects bool
	}{
		{
			reqBody:      `{"title":"Test Post","content":"Test Content"}`,
			post:         models.Post{ID: 1, Title: "Test Post", Content: "Test Content"},
			status:       http.StatusCreated,
			errorExpects: false,
		},
		{
			reqBody:      `{"title":"Test Post","content":""}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
		},
		{
			reqBody:      `{"title":"","content":"Content"}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
		},
		{
			reqBody:      `{"title":"","content":"Content"}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
		},
		{
			reqBody:      `{"field":"value"}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
		},
		{
			reqBody:      `{"title":123,"content":456}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
		},
	}

	for i, tc := range testCases {
		mockService := new(MockService)
		h := NewHandler(mockService)
		e := echo.New()

		if !tc.errorExpects {
			mockService.On("AddPost", mock.Anything).Return(tc.post, nil)
		}

		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.AddPost(c)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("case %d", i))

		if tc.errorExpects {
			resp := ErrorResponse{}

			err = json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				assert.NoError(t, err, fmt.Sprintf("case %d", i))
			}

			assert.Equal(t, "invalid post data", resp.Error, fmt.Sprintf("case %d", i))
		} else {
			resp := models.Post{}

			err = json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.Equal(t, "Test Post", resp.Title, fmt.Sprintf("case %d", i))
			assert.Equal(t, "Test Content", resp.Content, fmt.Sprintf("case %d", i))
		}

		mockService.AssertExpectations(t)
	}

}

func TestGetAllPosts(t *testing.T) {
	testCases := [][]models.Post{
		{
			{ID: 1, Title: "Post 1", Content: "Content 1"},
			{ID: 2, Title: "Post 2", Content: "Content 2"},
		},
		{
			{ID: 1, Title: "Post 1", Content: "Content 1"},
		},
		{},
	}

	for i, tc := range testCases {
		mockService := new(MockService)
		h := NewHandler(mockService)
		e := echo.New()

		mockService.On("GetAllPosts").Return(tc, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetAllPosts(c)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, http.StatusOK, rec.Code, fmt.Sprintf("case %d", i))

		resp := []models.Post{}
		err = json.Unmarshal([]byte(rec.Body.String()), &resp)
		if err != nil {
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
		}

		assert.Len(t, resp, len(tc), fmt.Sprintf("case %d", i))

		for j, post := range resp {
			assert.Equal(t, tc[j].Title, post.Title, fmt.Sprintf("case %d", i))
			assert.Equal(t, tc[j].Content, post.Content, fmt.Sprintf("case %d", i))
		}

		mockService.AssertExpectations(t)
	}

}

func TestUpdatePost(t *testing.T) {
	testCases := []struct {
		id           string
		reqBody      string
		post         models.Post
		status       int
		errorExpects bool
		errorMessage string
	}{
		{
			id:           "1",
			reqBody:      `{"title":"Updated Post","content":"Updated Content"}`,
			post:         models.Post{ID: 1, Title: "Updated Post", Content: "Updated Content"},
			status:       http.StatusOK,
			errorExpects: false,
		},
		{
			reqBody:      `{"title":"Updated Post","content":"Updated Content"}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
			errorMessage: "invalid post id",
		},
		{
			id:           "1",
			reqBody:      `{"content":"Updated Content"}`,
			post:         models.Post{ID: 1, Title: "Post", Content: "Updated Content"},
			status:       http.StatusOK,
			errorExpects: false,
		},
		{
			id:           "1",
			reqBody:      `{"title":"Updated Post"}`,
			post:         models.Post{ID: 1, Title: "Updated Post", Content: "Content"},
			status:       http.StatusOK,
			errorExpects: false,
		},
		{
			id:           "1",
			reqBody:      `{"title111":"Updated Post111"}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
			errorMessage: "invalid post data",
		},
		{
			id:           "1",
			reqBody:      `{}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
			errorMessage: "invalid post data",
		},
		{
			id:           "1",
			reqBody:      `{"title":"","content":""}`,
			post:         models.Post{},
			status:       http.StatusBadRequest,
			errorExpects: true,
			errorMessage: "invalid post data",
		},
	}

	for i, tc := range testCases {
		mockService := new(MockService)
		h := NewHandler(mockService)
		e := echo.New()

		if !tc.errorExpects {
			mockService.On("UpdatePost", mock.Anything).Return(tc.post, nil).Once()
		}

		req := httptest.NewRequest(http.MethodPut, "/posts/"+tc.id, bytes.NewBufferString(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/posts/:id")
		c.SetParamNames("id")
		c.SetParamValues(tc.id)

		err := h.UpdatePost(c)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("case %d", i))

		if tc.errorExpects {
			resp := ErrorResponse{}

			err = json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				assert.NoError(t, err, fmt.Sprintf("case %d", i))
			}

			assert.Equal(t, tc.errorMessage, resp.Error, fmt.Sprintf("case %d", i))
		} else {
			resp := models.Post{}

			err = json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				assert.NoError(t, err, fmt.Sprintf("case %d", i))
			}

			assert.Equal(t, tc.post.Title, resp.Title, fmt.Sprintf("case %d", i))
			assert.Equal(t, tc.post.Content, resp.Content, fmt.Sprintf("case %d", i))
		}
		mockService.AssertExpectations(t)
	}
}

func TestGetPost(t *testing.T) {
	testCases := []struct {
		id           string
		post         models.Post
		errorExpects bool
		errorMessage string
		status       int
	}{
		{
			id:           "1",
			post:         models.Post{ID: 1, Title: "Test Post", Content: "Test Content"},
			errorExpects: false,
			status:       http.StatusOK,
		},
		{
			id:           "0",
			post:         models.Post{},
			errorExpects: true,
			status:       http.StatusBadRequest,
			errorMessage: "invalid post id",
		},
		{
			id:           "a",
			post:         models.Post{},
			errorExpects: true,
			status:       http.StatusBadRequest,
			errorMessage: "invalid post id",
		},
	}

	for i, tc := range testCases {
		mockService := new(MockService)
		h := NewHandler(mockService)
		e := echo.New()

		if !tc.errorExpects {
			mockService.On("GetPost", 1).Return(tc.post, nil)
		}

		req := httptest.NewRequest(http.MethodGet, "/posts/"+tc.id, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/posts/:id")
		c.SetParamNames("id")
		c.SetParamValues(tc.id)

		err := h.GetPost(c)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("case %d", i))

		if tc.errorExpects {
			resp := ErrorResponse{}

			err = json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				assert.NoError(t, err, fmt.Sprintf("case %d", i))
			}

			assert.Equal(t, tc.errorMessage, resp.Error, fmt.Sprintf("case %d", i))
		} else {
			resp := models.Post{}

			err = json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.post.Title, resp.Title, fmt.Sprintf("case %d", i))
			assert.Equal(t, tc.post.Content, resp.Content, fmt.Sprintf("case %d", i))
		}
		mockService.AssertExpectations(t)
	}
}

func TestDeletePost(t *testing.T) {
	testCases := []struct {
		id           string
		post         models.Post
		errorExpects bool
		errorMessage string
		status       int
	}{
		{
			id:           "1",
			errorExpects: false,
			status:       http.StatusNoContent,
		},
		{
			id:           "0",
			errorExpects: true,
			status:       http.StatusBadRequest,
			errorMessage: "invalid post id",
		},
		{
			id:           "a",
			errorExpects: true,
			status:       http.StatusBadRequest,
			errorMessage: "invalid post id",
		},
	}

	for i, tc := range testCases {
		mockService := new(MockService)
		h := NewHandler(mockService)
		e := echo.New()

		if !tc.errorExpects {
			mockService.On("DeletePost", 1).Return(nil)
		}

		req := httptest.NewRequest(http.MethodDelete, "/posts/"+tc.id, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/posts/:id")
		c.SetParamNames("id")
		c.SetParamValues(tc.id)

		err := h.DeletePost(c)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("case %d", i))

		mockService.AssertExpectations(t)
	}

}
