package handler

import (
	"test/internal/service"
)

type Handler struct {
	service		*service.Service
}

func NewHandler(s *service.Service) *Handler {
	handler := &Handler{
		service:     s,
	}

	return handler
}
