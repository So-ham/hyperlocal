package services

import (
	"time"

	"github.com/google/uuid"
)

// CreateCommentRequest represents the request body for creating a comment
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

// CommentResponse represents the response for a comment
type CommentResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Username  *string   `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateComment creates a new comment on a post
func (s *service) CreateComment(req CreateCommentRequest, postID, userID uuid.UUID) (*CommentResponse, error) {
	// Create the comment
	comment, err := s.model.CreateComment(postID, userID, req.Content)
	if err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.model.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &CommentResponse{
		ID:        comment.ID.String(),
		Content:   comment.Content,
		Username:  user.Username,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// GetCommentsByPostID retrieves all comments for a post
func (s *service) GetCommentsByPostID(postID uuid.UUID) ([]CommentResponse, error) {
	// Get comments
	comments, err := s.model.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		response[i] = CommentResponse{
			ID:        comment.ID.String(),
			Content:   comment.Content,
			Username:  comment.User.Username,
			CreatedAt: comment.CreatedAt,
		}
	}

	return response, nil
}
