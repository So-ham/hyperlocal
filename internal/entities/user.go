package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Username     *string   `gorm:"unique"`
	PasswordHash string
	IsBanned     bool `gorm:"default:false"`
	CreatedAt    time.Time
}

// Post represents a post in the system
type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Content   string
	Latitude  float64
	Longitude float64
	Upvotes   int  `gorm:"default:0"`
	Downvotes int  `gorm:"default:0"`
	IsFlagged bool `gorm:"default:false"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	PostID    uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Content   string
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

// Report represents a report of a post
type Report struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	PostID    uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Reason    string
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

// UserPostVote tracks user votes on posts to prevent multiple votes
type UserPostVote struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	PostID    uuid.UUID `gorm:"type:uuid"`
	VoteType  string    // "upvote" or "downvote"
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

// RefreshToken stores refresh tokens for users
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Token     string    `gorm:"unique"`
	ExpiresAt time.Time
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}
