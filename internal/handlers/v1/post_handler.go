package v1

import (
	"encoding/json"
	"hyperlocal/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreatePost handles post creation
// @Summary Create a new post
// @Description Create a new post with content and location
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreatePostRequest true "Post details"
// @Success 201 {object} services.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts [post]
func (h *handlerV1) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req services.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	post, err := h.Service.CreatePost(req, userID.(uuid.UUID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetNearbyPosts handles retrieving posts near a location
// @Summary Get nearby posts
// @Description Get posts within 5km of the specified location
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Success 200 {array} services.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts [get]
func (h *handlerV1) GetNearbyPosts(w http.ResponseWriter, r *http.Request) {
	// Get latitude and longitude from query parameters
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		http.Error(w, "Latitude and longitude are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	posts, err := h.Service.GetNearbyPosts(lat, lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// UpvotePost handles upvoting a post
// @Summary Upvote a post
// @Description Upvote a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/upvote [post]
func (h *handlerV1) UpvotePost(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Service.UpvotePost(postID, userID.(uuid.UUID)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post upvoted successfully"})
}

// DownvotePost handles downvoting a post
// @Summary Downvote a post
// @Description Downvote a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/downvote [post]
func (h *handlerV1) DownvotePost(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Service.DownvotePost(postID, userID.(uuid.UUID)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post downvoted successfully"})
}

// ReportPost handles reporting a post
// @Summary Report a post
// @Description Report a post by ID with a reason
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body services.ReportPostRequest true "Report details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/report [post]
func (h *handlerV1) ReportPost(w http.ResponseWriter, r *http.Request) {
	var req services.ReportPostRequest
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

	if err := h.Service.ReportPost(req, postID, userID.(uuid.UUID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post reported successfully"})
}

// DeletePost handles deleting a post (admin only)
// @Summary Delete a post
// @Description Delete a post by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/posts/{id} [delete]
func (h *handlerV1) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Get post ID from URL
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeletePost(postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}