package services

import (
	"github.com/google/uuid"
)

// ReportPostRequest represents the request body for reporting a post
type ReportPostRequest struct {
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

// ReportPost reports a post
func (s *service) ReportPost(req ReportPostRequest, postID, userID uuid.UUID) error {
	// Create the report
	_, err := s.model.CreateReport(postID, userID, req.Reason)
	return err
}
