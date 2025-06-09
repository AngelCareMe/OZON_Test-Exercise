package services

import (
	"time"

	"ozon_test/internal/models"
	"ozon_test/internal/storage"
)

type PostService struct {
	storage storage.PostStorage
}

func NewPostService(storage storage.PostStorage) *PostService {
	return &PostService{storage: storage}
}

func (s *PostService) CreatePost(title, text, author string) (*models.Post, error) {
	post := &models.Post{
		Title:         title,
		Text:          text,
		AllowComments: true,
		Author:        author,
		CreatedAt:     time.Now(),
	}
	err := s.storage.CreatePost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetAllPosts() ([]*models.Post, error) {
	return s.storage.GetAllPosts()
}

func (s *PostService) DisableComments(postID int) error {
	post, err := s.storage.GetPostByID(postID)
	if err != nil {
		return err
	}
	post.AllowComments = false
	return s.storage.UpdatePost(post)
}
