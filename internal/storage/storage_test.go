package storage

import (
	"testing"
	"time"

	"ozon_test/internal/models"
)

func TestInMemoryPostStorage(t *testing.T) {
	store := NewInMemoryPostStorage()

	// Тест создания поста
	post := &models.Post{
		Title:         "Test",
		Text:          "Text",
		AllowComments: true,
		Author:        "Author",
		CreatedAt:     time.Now(),
	}
	if err := store.CreatePost(post); err != nil {
		t.Fatal(err)
	}

	// Тест получения поста по ID
	retrieved, err := store.GetPostByID(post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if retrieved.Title != "Test" {
		t.Errorf("Ожидался заголовок 'Test', получено '%s'", retrieved.Title)
	}

	// Тест получения всех постов
	posts, err := store.GetAllPosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 {
		t.Errorf("Ожидался 1 пост, получено %d", len(posts))
	}
}

func TestInMemoryCommentStorage(t *testing.T) {
	store := NewInMemoryCommentStorage()

	// Тест создания комментария
	comment := &models.Comment{
		PostID:    1,
		Text:      "Test comment",
		Author:    "User",
		CreatedAt: time.Now(),
	}
	if err := store.CreateComment(comment); err != nil {
		t.Fatal(err)
	}

	// Тест получения комментариев с пагинацией
	comments, err := store.GetCommentsByPostID(1, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 {
		t.Errorf("Ожидался 1 комментарий, получено %d", len(comments))
	}
	if comments[0].Text != "Test comment" {
		t.Errorf("Ожидался текст 'Test comment', получено '%s'", comments[0].Text)
	}
}
