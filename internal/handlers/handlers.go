package handlers

import (
	v1 "hyperlocal/internal/handlers/v1"
	"hyperlocal/internal/services"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	V1 v1.HandlerV1
}

func New(service services.Service, validate *validator.Validate) *Handler {
	return &Handler{
		V1: v1.New(service, validate),
	}

}
