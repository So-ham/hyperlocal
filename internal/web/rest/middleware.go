package rest

import (
	"context"
	"hyperlocal/internal/services"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AuthMiddleware authenticates requests using JWT
func AuthMiddleware(service services.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Check if the header has the Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := service.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Parse the user ID
			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			// Set the user ID in the context
			ctx := context.WithValue(r.Context(), "userID", userID)
			ctx = context.WithValue(ctx, "role", claims.Role)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthMiddlewareFunc returns a middleware function compatible with r.With()
func AuthMiddlewareFunc(service services.Service) func(http.Handler) http.Handler {
	return AuthMiddleware(service)
}

// AdminMiddleware ensures the user has admin role
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the role from the context
		role := r.Context().Value("role")
		if role == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		if role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimiterMiddleware implements a simple in-memory rate limiter
func RateLimiterMiddleware(next http.Handler) http.Handler {
	// Simple in-memory store for rate limiting
	// In production, use Redis or a similar distributed store
	type rateLimitEntry struct {
		count    int
		lastSeen int64
	}
	
	// Store user IDs and their request counts
	// In a real app, this would be a Redis cache or similar
	rateLimits := make(map[string]rateLimitEntry)
	var mu sync.Mutex
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the context
		userIDValue := r.Context().Value("userID")
		if userIDValue == nil {
			next.ServeHTTP(w, r)
			return
		}
		
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		
		userIDStr := userID.String()
		
		mu.Lock()
		defer mu.Unlock()
		
		now := time.Now().Unix()
		entry, exists := rateLimits[userIDStr]
		
		// Reset count if it's been more than a minute
		if exists && now - entry.lastSeen > 60 {
			entry.count = 0
		}
		
		// Update the entry
		entry.count++
		entry.lastSeen = now
		rateLimits[userIDStr] = entry
		
		// Check if the user has exceeded the rate limit
		if entry.count > 10 { // 10 requests per minute
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
		// Get the user ID from the context
		userIDValue := r.Context().Value("userID")
		if !exists {
			c.Next()
			return
		}
		
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			c.Next()
			return
		}
		
		userIDStr := userID.String()
		
		// Check if this is a POST request to /posts
		if c.Request.Method == "POST" && strings.HasPrefix(c.Request.URL.Path, "/posts") && !strings.Contains(c.Request.URL.Path, "/comments") {
			// Get the current timestamp
			now := c.Request.Context().Value("requestTime").(int64)
			
			// Get the user's rate limit entry
			entry, exists := rateLimits[userIDStr]
			
			// If the entry doesn't exist or is older than an hour, create a new one
			if !exists || now-entry.lastSeen > 3600 {
				rateLimits[userIDStr] = rateLimitEntry{
					count:    1,
					lastSeen: now,
				}
			} else {
				// If the user has already made 3 requests in the last hour, reject the request
				if entry.count >= 3 {
					c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded (3 posts per hour)"})
					c.Abort()
					return
				}
				
				// Otherwise, increment the count
				entry.count++
				entry.lastSeen = now
				rateLimits[userIDStr] = entry
			}
		}
		
		c.Next()
	}
}