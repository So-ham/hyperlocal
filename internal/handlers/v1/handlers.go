package v1

import (
	"hyperlocal/internal/services"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type handlerV1 struct {
	Service  services.Service
	Validate *validator.Validate
}

type HandlerV1 interface {
	// Auth handlers
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	
	// Post handlers
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetNearbyPosts(w http.ResponseWriter, r *http.Request)
	UpvotePost(w http.ResponseWriter, r *http.Request)
	DownvotePost(w http.ResponseWriter, r *http.Request)
	ReportPost(w http.ResponseWriter, r *http.Request)
	
	// Comment handlers
	CreateComment(w http.ResponseWriter, r *http.Request)
	GetComments(w http.ResponseWriter, r *http.Request)
	
	// Admin handlers
	GetFlaggedPosts(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	BanUser(w http.ResponseWriter, r *http.Request)
}

func New(s services.Service, v *validator.Validate) HandlerV1 {
	return &handlerV1{Service: s, Validate: v}
}
