package services

import (
	"github.com/google/uuid"
)

// BanUser bans a user
func (s *service) BanUser(userID uuid.UUID) error {
	return s.model.BanUser(userID, true)
}
