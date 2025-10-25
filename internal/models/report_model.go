package models

import (
	"hyperlocal/internal/entities"
	"time"

	"github.com/google/uuid"
)

// CreateReport creates a new report for a post
func (m *Model) CreateReport(postID, userID uuid.UUID, reason string) (*entities.Report, error) {
	report := &entities.Report{
		ID:        uuid.New(),
		PostID:    postID,
		UserID:    userID,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	if err := m.db.Create(report).Error; err != nil {
		return nil, err
	}

	// Check if post has reached the report threshold (3 reports)
	var reportCount int64
	if err := m.db.Model(&entities.Report{}).Where("post_id = ?", postID).Count(&reportCount).Error; err != nil {
		return nil, err
	}

	// If 3 or more reports, flag the post
	if reportCount >= 3 {
		if err := m.FlagPost(postID); err != nil {
			return nil, err
		}
	}

	return report, nil
}

// GetReportsByPostID retrieves all reports for a post
func (m *Model) GetReportsByPostID(postID uuid.UUID) ([]entities.Report, error) {
	var reports []entities.Report
	if err := m.db.Where("post_id = ?", postID).Preload("User").Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}
