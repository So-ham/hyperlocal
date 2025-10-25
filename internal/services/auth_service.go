package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse represents the response for successful authentication
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // in seconds
}

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register creates a new user account
func (s *service) Register(req RegisterRequest) (*TokenResponse, error) {
	// Check if username is already taken
	_, err := s.model.GetUserByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already taken")
	}

	// Create the user
	var username *string = &req.Username
	user, err := s.model.CreateUser(username, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(user.ID)
}

// Login authenticates a user and returns tokens
func (s *service) Login(req LoginRequest) (*TokenResponse, error) {
	// Find the user
	user, err := s.model.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, errors.New("account is banned")
	}

	// Verify password
	if !s.model.VerifyPassword(user, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	return s.generateTokens(user.ID)
}

// RefreshToken refreshes the access token using a refresh token
func (s *service) RefreshToken(req RefreshTokenRequest) (*TokenResponse, error) {
	// Get the refresh token from the database
	refreshToken, err := s.model.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if the token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		// Delete the expired token
		s.model.DeleteRefreshToken(req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	// Get the user
	user, err := s.model.GetUserByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, errors.New("account is banned")
	}

	// Delete the old refresh token
	if err := s.model.DeleteRefreshToken(req.RefreshToken); err != nil {
		return nil, err
	}

	// Generate new tokens
	return s.generateTokens(user.ID)
}

// generateTokens generates access and refresh tokens for a user
func (s *service) generateTokens(userID uuid.UUID) (*TokenResponse, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	// Set token expiration times
	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour) // 30 days

	// Create the claims for the access token
	claims := JWTClaims{
		UserID: userID.String(),
		Role:   "user", // Default role
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hyperlocal",
			Subject:   userID.String(),
		},
	}

	// Create the access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	// Generate a refresh token (a random UUID)
	refreshToken := uuid.New().String()

	// Store the refresh token in the database
	if err := s.model.StoreRefreshToken(userID, refreshToken, refreshTokenExpiry); err != nil {
		return nil, err
	}

	// Return the tokens
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTokenExpiry.Sub(time.Now()).Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *service) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
