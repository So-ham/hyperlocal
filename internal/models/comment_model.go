package models

import (
	"hyperlocal/internal/entities"
	"time"

	"github.com/google/uuid"
)

// CreateComment creates a new comment on a post
func (m *Model) CreateComment(postID, userID uuid.UUID, content string) (*entities.Comment, error) {
	comment := &entities.Comment{
		ID:        uuid.New(),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := m.db.Create(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByPostID retrieves all comments for a post
func (m *Model) GetCommentsByPostID(postID uuid.UUID) ([]entities.Comment, error) {
	var comments []entities.Comment
	if err := m.db.Where("post_id = ?", postID).Preload("User").Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// DeleteComment deletes a comment
func (m *Model) DeleteComment(id uuid.UUID) error {
	return m.db.Delete(&entities.Comment{}, "id = ?", id).Error
}
