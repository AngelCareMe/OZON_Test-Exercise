package services

import (
	"testing"

	"ozon_test/internal/storage"
)

func TestCreateComment(t *testing.T) {
	postStorage := storage.NewInMemoryPostStorage()
	commentStorage := storage.NewInMemoryCommentStorage()
	service := NewCommentService(commentStorage, postStorage)
	post, _ := NewPostService(postStorage).CreatePost("Test", "Text", "Author")
	comment, err := service.CreateComment(post.ID, nil, "Test comment", "User")
	if err != nil {
		t.Fatal(err)
	}
	if comment.Text != "Test comment" {
		t.Errorf("Ожидался текст 'Test comment', получено '%s'", comment.Text)
	}
}

func TestCreateCommentExceedsLimit(t *testing.T) {
	postStorage := storage.NewInMemoryPostStorage()
	commentStorage := storage.NewInMemoryCommentStorage()
	service := NewCommentService(commentStorage, postStorage)
	post, _ := NewPostService(postStorage).CreatePost("Test", "Text", "Author")
	longText := string(make([]byte, 2001))
	_, err := service.CreateComment(post.ID, nil, longText, "User")
	if err == nil {
		t.Error("Ожидалась ошибка для текста, превышающего 2000 символов")
	}
}

func TestCreateCommentWhenCommentsDisabled(t *testing.T) {
	postStorage := storage.NewInMemoryPostStorage()
	commentStorage := storage.NewInMemoryCommentStorage()
	postService := NewPostService(postStorage)
	commentService := NewCommentService(commentStorage, postStorage)
	post, _ := postService.CreatePost("Test", "Text", "Author")
	_ = postService.DisableComments(post.ID)
	_, err := commentService.CreateComment(post.ID, nil, "Test comment", "User")
	if err != storage.ErrCommentsNotAllowed {
		t.Error("Ожидалась ошибка ErrCommentsNotAllowed")
	}
}
