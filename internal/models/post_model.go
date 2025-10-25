package models

import (
	"hyperlocal/internal/entities"
	"time"

	"github.com/google/uuid"
)

// CreatePost creates a new post
func (m *Model) CreatePost(userID uuid.UUID, content string, latitude, longitude float64) (*entities.Post, error) {
	post := &entities.Post{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
	}

	if err := m.db.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

// GetPostByID retrieves a post by ID
func (m *Model) GetPostByID(id uuid.UUID) (*entities.Post, error) {
	var post entities.Post
	if err := m.db.Preload("User").First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetNearbyPosts retrieves posts within a specified radius of a location
func (m *Model) GetNearbyPosts(latitude, longitude float64, radiusMeters float64) ([]entities.Post, error) {
	var posts []entities.Post

	// Using PostGIS ST_DWithin to find posts within radius
	// ST_MakePoint creates a point from longitude and latitude (note the order)
	// ST_SetSRID sets the spatial reference system identifier (SRID) to 4326 (WGS84)
	// ST_DWithin checks if the distance between two geometries is within a given value
	query := `
		SELECT * FROM posts 
		WHERE ST_DWithin(
			ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		ORDER BY created_at DESC
	`

	if err := m.db.Raw(query, longitude, latitude, radiusMeters).Scan(&posts).Error; err != nil {
		return nil, err
	}

	// Load user data for each post
	for i := range posts {
		if err := m.db.Model(&posts[i]).Association("User").Find(&posts[i].User); err != nil {
			return nil, err
		}
	}

	return posts, nil
}

// UpdatePost updates a post
func (m *Model) UpdatePost(post *entities.Post) error {
	return m.db.Save(post).Error
}

// DeletePost deletes a post
func (m *Model) DeletePost(id uuid.UUID) error {
	return m.db.Delete(&entities.Post{}, "id = ?", id).Error
}

// FlagPost marks a post as flagged
func (m *Model) FlagPost(id uuid.UUID) error {
	return m.db.Model(&entities.Post{}).Where("id = ?", id).Update("is_flagged", true).Error
}

// GetFlaggedPosts retrieves all flagged posts
func (m *Model) GetFlaggedPosts() ([]entities.Post, error) {
	var posts []entities.Post
	if err := m.db.Where("is_flagged = ?", true).Preload("User").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// UpvotePost increments the upvote count for a post
func (m *Model) UpvotePost(postID uuid.UUID) error {
	return m.db.Model(&entities.Post{}).Where("id = ?", postID).Update("upvotes", m.db.Raw("upvotes + 1")).Error
}

// DownvotePost increments the downvote count for a post
func (m *Model) DownvotePost(postID uuid.UUID) error {
	return m.db.Model(&entities.Post{}).Where("id = ?", postID).Update("downvotes", m.db.Raw("downvotes + 1")).Error
}
