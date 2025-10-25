package v1

import (
	"encoding/json"
	"hyperlocal/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateComment handles comment creation
// @Summary Create a new comment
// @Description Create a new comment on a post
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body services.CreateCommentRequest true "Comment details"
// @Success 201 {object} services.CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/comments [post]
func (h *handlerV1) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req services.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get post ID from URL
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	comment, err := h.Service.CreateComment(req, postID, userID.(uuid.UUID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetComments handles retrieving comments for a post
// @Summary Get comments for a post
// @Description Get all comments for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {array} services.CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/comments [get]
func (h *handlerV1) GetComments(w http.ResponseWriter, r *http.Request) {
	// Get post ID from URL
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := h.Service.GetCommentsByPostID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}