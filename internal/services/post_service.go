package services

import (
	"time"

	"github.com/google/uuid"
)

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Content   string  `json:"content" validate:"required,min=1,max=500"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// PostResponse represents the response for a post
type PostResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Username  *string   `json:"username"`
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
	CreatedAt time.Time `json:"created_at"`
	IsFlagged bool      `json:"is_flagged,omitempty"`
}

// CreatePost creates a new post
func (s *service) CreatePost(req CreatePostRequest, userID uuid.UUID) (*PostResponse, error) {
	// Create the post
	post, err := s.model.CreatePost(userID, req.Content, req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.model.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &PostResponse{
		ID:        post.ID.String(),
		Content:   post.Content,
		Username:  user.Username,
		Upvotes:   post.Upvotes,
		Downvotes: post.Downvotes,
		CreatedAt: post.CreatedAt,
	}, nil
}

// GetNearbyPosts retrieves posts within a specified radius of a location
func (s *service) GetNearbyPosts(latitude, longitude float64) ([]PostResponse, error) {
	// Get posts within 5km radius
	posts, err := s.model.GetNearbyPosts(latitude, longitude, 5000)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]PostResponse, len(posts))
	for i, post := range posts {
		response[i] = PostResponse{
			ID:        post.ID.String(),
			Content:   post.Content,
			Username:  post.User.Username,
			Upvotes:   post.Upvotes,
			Downvotes: post.Downvotes,
			CreatedAt: post.CreatedAt,
		}
	}

	return response, nil
}

// GetPostByID retrieves a post by ID
func (s *service) GetPostByID(id uuid.UUID) (*PostResponse, error) {
	// Get the post
	post, err := s.model.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &PostResponse{
		ID:        post.ID.String(),
		Content:   post.Content,
		Username:  post.User.Username,
		Upvotes:   post.Upvotes,
		Downvotes: post.Downvotes,
		CreatedAt: post.CreatedAt,
		IsFlagged: post.IsFlagged,
	}, nil
}

// DeletePost deletes a post
func (s *service) DeletePost(id uuid.UUID) error {
	return s.model.DeletePost(id)
}

// UpvotePost upvotes a post
func (s *service) UpvotePost(postID, userID uuid.UUID) error {
	return s.model.VoteOnPost(userID, postID, "upvote")
}

// DownvotePost downvotes a post
func (s *service) DownvotePost(postID, userID uuid.UUID) error {
	return s.model.VoteOnPost(userID, postID, "downvote")
}

// GetFlaggedPosts retrieves all flagged posts
func (s *service) GetFlaggedPosts() ([]PostResponse, error) {
	// Get flagged posts
	posts, err := s.model.GetFlaggedPosts()
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]PostResponse, len(posts))
	for i, post := range posts {
		response[i] = PostResponse{
			ID:        post.ID.String(),
			Content:   post.Content,
			Username:  post.User.Username,
			Upvotes:   post.Upvotes,
			Downvotes: post.Downvotes,
			CreatedAt: post.CreatedAt,
			IsFlagged: post.IsFlagged,
		}
	}

	return response, nil
}
