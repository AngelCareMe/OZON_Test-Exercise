package services

import (
	"testing"

	"ozon_test/internal/storage"
)

func TestCreatePost(t *testing.T) {
	postStorage := storage.NewInMemoryPostStorage()
	service := NewPostService(postStorage)
	post, err := service.CreatePost("Test", "Text", "Author")
	if err != nil {
		t.Fatal(err)
	}
	if post.Title != "Test" {
		t.Errorf("Ожидался заголовок 'Test', получено '%s'", post.Title)
	}
}

func TestDisableComments(t *testing.T) {
	postStorage := storage.NewInMemoryPostStorage()
	service := NewPostService(postStorage)
	post, _ := service.CreatePost("Test", "Text", "Author")
	err := service.DisableComments(post.ID)
	if err != nil {
		t.Fatal(err)
	}
	updatedPost, err := postStorage.GetPostByID(post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedPost.AllowComments {
		t.Error("Ожидалось, что комментарии будут отключены")
	}
}
