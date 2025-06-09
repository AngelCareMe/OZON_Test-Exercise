package api

import (
	"encoding/json"
	"net/http"

	"ozon_test/internal/services"
)

type PostHandler struct {
	service *services.PostService
}

func NewPostHandler(service *services.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string `json:"title"`
		Text   string `json:"text"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	post, err := h.service.CreatePost(req.Title, req.Text, req.Author)
	if err != nil {
		http.Error(w, "Не удалось создать пост", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		http.Error(w, "Не удалось получить посты", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) DisableComments(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostID int `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if err := h.service.DisableComments(req.PostID); err != nil {
		http.Error(w, "Не удалось отключить комментарии", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
