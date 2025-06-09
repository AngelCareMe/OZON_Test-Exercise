package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ozon_test/internal/services"
	"ozon_test/internal/storage"
)

func TestGetAllPosts(t *testing.T) {
	postStorage := storage.NewInMemoryPostStorage()
	postService := services.NewPostService(postStorage)
	handler := NewPostHandler(postService)

	_, _ = postService.CreatePost("Test1", "Text1", "Author1")
	_, _ = postService.CreatePost("Test2", "Text2", "Author2")

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllPosts(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался код 200, получено %v", status)
	}

	var posts []struct {
		ID            int    `json:"id"`
		Title         string `json:"title"`
		Text          string `json:"text"`
		AllowComments bool   `json:"allow_comments"`
		Author        string `json:"author"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Errorf("Ожидалось 2 поста, получено %d", len(posts))
	}
}

func TestCreateCommentInvalidJSON(t *testing.T) {
	commentStorage := storage.NewInMemoryCommentStorage()
	postStorage := storage.NewInMemoryPostStorage()
	commentService := services.NewCommentService(commentStorage, postStorage)
	handler := NewCommentHandler(commentService)

	req, err := http.NewRequest("POST", "/comments/create", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateComment(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Ожидался код 400, получено %v", status)
	}
}
