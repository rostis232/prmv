package service

import (
	"github.com/rostis232/prmv/models"
)

type Service struct {
	Repo Repository
}

type Repository interface {
	AddPost(post models.Post) (int, error)
	GetAllPosts() ([]models.Post, error)
	UpdatePost(post models.Post) (int, error)
	GetPost(id int) (models.Post, error)
	DeletePost(id int) error
}

func NewService(repo Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) AddPost(newPost models.Post) (models.Post, error) {
	id, err := s.Repo.AddPost(newPost)
	if err != nil {
		return models.Post{}, err
	}

	post, err := s.Repo.GetPost(id)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *Service) GetAllPosts() ([]models.Post, error) {
	posts, err := s.Repo.GetAllPosts()
	if err != nil {
		return []models.Post{}, err
	}

	return posts, nil
}

func (s *Service) UpdatePost(updatedPost models.Post) (models.Post, error) {
	post, err := s.Repo.GetPost(updatedPost.ID)
	if err != nil {
		return models.Post{}, err
	}
	if updatedPost.Title != "" && post.Title != updatedPost.Title {
		post.Title = updatedPost.Title
	}
	if updatedPost.Content != "" && post.Content != updatedPost.Content {
		post.Content = updatedPost.Content
	}
	id, err := s.Repo.UpdatePost(post)
	if err != nil {
		return models.Post{}, err
	}

	post, err = s.Repo.GetPost(id)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (s *Service) GetPost(id int) (models.Post, error) {
	post, err := s.Repo.GetPost(id)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (s *Service) DeletePost(id int) error {
	err := s.Repo.DeletePost(id)
	if err != nil {
		return err
	}
	return nil
}
