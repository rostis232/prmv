package service

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/rostis232/prmv/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddPost(post models.Post) (int, error) {
	args := m.Called(post)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetAllPosts() ([]models.Post, error) {
	args := m.Called()
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockRepository) UpdatePost(post models.Post) (int, error) {
	args := m.Called(post)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetPost(id int) (models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *MockRepository) DeletePost(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAddPost(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	post := models.Post{Title: "Test Title", Content: "Test Content"}
	mockRepo.On("AddPost", post).Return(1, nil)
	mockRepo.On("GetPost", 1).Return(post, nil)

	result, err := service.AddPost(post)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, result.Title)
	assert.Equal(t, post.Content, result.Content)

	mockRepo.AssertExpectations(t)
}

func TestGetAllPosts(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	posts := []models.Post{
		{Title: "Test Title 1", Content: "Test Content 1"},
		{Title: "Test Title 2", Content: "Test Content 2"},
	}

	mockRepo.On("GetAllPosts").Return(posts, nil)

	result, err := service.GetAllPosts()
	assert.NoError(t, err)
	assert.Equal(t, posts, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdatePost(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	originalPost := models.Post{ID: 1, Title: "Original Title", Content: "Original Content"}
	updatedPost := models.Post{ID: 1, Title: "Updated Title", Content: "Updated Content"}
	expectedPost := models.Post{ID: 1, Title: "Updated Title", Content: "Updated Content"}

	// Set up the expected calls
	mockRepo.On("GetPost", 1).Return(originalPost, nil).Once()
	mockRepo.On("UpdatePost", mock.AnythingOfType("models.Post")).Return(1, nil).Once()
	mockRepo.On("GetPost", 1).Return(expectedPost, nil).Once()

	result, err := service.UpdatePost(updatedPost)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdatePostValidation(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	invalidPost := models.Post{ID: 1, Title: "", Content: ""}

	_, err := service.UpdatePost(invalidPost)

	assert.Error(t, err)
}

func TestGetPost(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	post := models.Post{ID: 1, Title: "Test Title", Content: "Test Content"}

	mockRepo.On("GetPost", 1).Return(post, nil)

	result, err := service.GetPost(1)
	assert.NoError(t, err)
	assert.Equal(t, post, result)

	mockRepo.AssertExpectations(t)
}

func TestDeletePost(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("DeletePost", 1).Return(nil)

	err := service.DeletePost(1)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAddPostValidation(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	invalidPost := models.Post{Title: "", Content: "Test Content"}

	_, err := service.AddPost(invalidPost)
	assert.Error(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
}
