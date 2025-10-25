package entities

import "errors"

var (
	ErrBookNotFound = errors.New("book not found")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCredentials = errors.New("invalid credentials")
)
