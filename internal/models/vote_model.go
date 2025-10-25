package models

import (
	"errors"
	"hyperlocal/internal/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VoteModel handles vote-related database operations
type VoteModel struct{}

// VoteOnPost records a user's vote on a post
func (m *Model) VoteOnPost(userID, postID uuid.UUID, voteType string) error {
	// Check if user has already voted on this post
	var existingVote entities.UserPostVote
	err := m.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingVote).Error

	// If no vote exists, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		vote := &entities.UserPostVote{
			ID:        uuid.New(),
			UserID:    userID,
			PostID:    postID,
			VoteType:  voteType,
			CreatedAt: time.Now(),
		}

		// Begin transaction
		tx := m.db.Begin()

		// Create the vote
		if err := tx.Create(vote).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Update the post's vote count
		if voteType == "upvote" {
			if err := tx.Model(&entities.Post{}).Where("id = ?", postID).Update("upvotes", gorm.Expr("upvotes + 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if voteType == "downvote" {
			if err := tx.Model(&entities.Post{}).Where("id = ?", postID).Update("downvotes", gorm.Expr("downvotes + 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// Commit transaction
		return tx.Commit().Error
	}

	// If vote exists but is different type, update it
	if err == nil && existingVote.VoteType != voteType {
		// Begin transaction
		tx := m.db.Begin()

		// Update the vote type
		if err := tx.Model(&existingVote).Update("vote_type", voteType).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Update the post's vote counts
		if existingVote.VoteType == "upvote" && voteType == "downvote" {
			// Change from upvote to downvote
			if err := tx.Model(&entities.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
				"upvotes":   gorm.Expr("upvotes - 1"),
				"downvotes": gorm.Expr("downvotes + 1"),
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if existingVote.VoteType == "downvote" && voteType == "upvote" {
			// Change from downvote to upvote
			if err := tx.Model(&entities.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
				"downvotes": gorm.Expr("downvotes - 1"),
				"upvotes":   gorm.Expr("upvotes + 1"),
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// Commit transaction
		return tx.Commit().Error
	}

	// If vote exists and is the same type, return error
	return errors.New("user has already voted on this post")
}

// HasUserVotedOnPost checks if a user has already voted on a post
func (m *Model) HasUserVotedOnPost(userID, postID uuid.UUID) (bool, string, error) {
	var vote entities.UserPostVote
	err := m.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&vote).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, "", nil
	}

	if err != nil {
		return false, "", err
	}

	return true, vote.VoteType, nil
}
