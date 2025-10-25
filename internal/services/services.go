package services

import (
	"hyperlocal/internal/models"

	"github.com/google/uuid"
)

// Service represents the service layer having
// all the services from all service packages
type service struct {
	model models.Model
}

// New creates a new instance of Service
func New(model *models.Model) Service {
	return &service{
		model: *model,
	}
}

// Service defines the interface for the service layer
type Service interface {
	// Auth services
	Register(req RegisterRequest) (*TokenResponse, error)
	Login(req LoginRequest) (*TokenResponse, error)
	RefreshToken(req RefreshTokenRequest) (*TokenResponse, error)
	ValidateToken(tokenString string) (*JWTClaims, error)

	// Post services
	CreatePost(req CreatePostRequest, userID uuid.UUID) (*PostResponse, error)
	GetNearbyPosts(latitude, longitude float64) ([]PostResponse, error)
	GetPostByID(id uuid.UUID) (*PostResponse, error)
	DeletePost(id uuid.UUID) error

	// Vote services
	UpvotePost(postID, userID uuid.UUID) error
	DownvotePost(postID, userID uuid.UUID) error

	// Comment services
	CreateComment(req CreateCommentRequest, postID, userID uuid.UUID) (*CommentResponse, error)
	GetCommentsByPostID(postID uuid.UUID) ([]CommentResponse, error)

	// Report services
	ReportPost(req ReportPostRequest, postID, userID uuid.UUID) error

	// Admin services
	GetFlaggedPosts() ([]PostResponse, error)
	BanUser(userID uuid.UUID) error
}
