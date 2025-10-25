package models

import (
	"gorm.io/gorm"
)

// Model represents the data access layer
type Model struct {
	db *gorm.DB
}

// New creates a new instance of Model
func New(gdb *gorm.DB) *Model {
	return &Model{
		db: gdb,
	}
}
