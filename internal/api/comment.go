package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ozon_test/internal/services"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostID          int    `json:"post_id"`
		ParentCommentID *int   `json:"parent_comment_id"`
		Text            string `json:"text"`
		Author          string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	comment, err := h.service.CreateComment(req.PostID, req.ParentCommentID, req.Text, req.Author)
	if err != nil {
		http.Error(w, "Не удалось создать комментарий: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}
	comments, err := h.service.GetCommentsByPostID(postID, limit, offset)
	if err != nil {
		http.Error(w, "Не удалось получить комментарии", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
