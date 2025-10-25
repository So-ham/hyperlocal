package models

import (
	"errors"
	"hyperlocal/internal/entities"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserModel handles user-related database operations
type UserModel struct{}

// CreateUser creates a new user in the database
func (m *Model) CreateUser(username *string, password string) (*entities.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := m.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (m *Model) GetUserByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := m.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (m *Model) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := m.db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (m *Model) UpdateUser(user *entities.User) error {
	return m.db.Save(user).Error
}

// BanUser sets the user's banned status
func (m *Model) BanUser(userID uuid.UUID, banned bool) error {
	return m.db.Model(&entities.User{}).Where("id = ?", userID).Update("is_banned", banned).Error
}

// StoreRefreshToken stores a refresh token for a user
func (m *Model) StoreRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	refreshToken := &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return m.db.Create(refreshToken).Error
}

// GetRefreshToken retrieves a refresh token
func (m *Model) GetRefreshToken(token string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	if err := m.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken deletes a refresh token
func (m *Model) DeleteRefreshToken(token string) error {
	return m.db.Where("token = ?", token).Delete(&entities.RefreshToken{}).Error
}

// DeleteExpiredRefreshTokens deletes all expired refresh tokens
func (m *Model) DeleteExpiredRefreshTokens() error {
	return m.db.Where("expires_at < ?", time.Now()).Delete(&entities.RefreshToken{}).Error
}

// VerifyPassword checks if the provided password matches the stored hash
func (m *Model) VerifyPassword(user *entities.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
