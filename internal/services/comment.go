package services

import (
	"errors"
	"time"

	"ozon_test/internal/models"
	"ozon_test/internal/storage"
)

type CommentService struct {
	storage     storage.CommentStorage
	postStorage storage.PostStorage
}

func NewCommentService(storage storage.CommentStorage, postStorage storage.PostStorage) *CommentService {
	return &CommentService{storage: storage, postStorage: postStorage}
}

func (s *CommentService) CreateComment(postID int, parentCommentID *int, text, author string) (*models.Comment, error) {
	if len(text) > 2000 {
		return nil, errors.New("текст комментария превышает 2000 символов")
	}
	post, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	if !post.AllowComments {
		return nil, storage.ErrCommentsNotAllowed
	}
	comment := &models.Comment{
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Text:            text,
		Author:          author,
		CreatedAt:       time.Now(),
	}
	err = s.storage.CreateComment(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID int, limit, offset int) ([]*models.Comment, error) {
	return s.storage.GetCommentsByPostID(postID, limit, offset)
}
