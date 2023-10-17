package service

import (
	"test/internal/service/archive"
	"test/internal/service/email"
	"test/config"
)

type Service struct {
	Archive *archive.Archive
	Email *email.Email
	Config *config.Config
}

func NewService(a *archive.Archive, e *email.Email, cfg *config.Config) *Service {
	return &Service{
		Archive: a,
		Email: e,
		Config: cfg,
	}
}
