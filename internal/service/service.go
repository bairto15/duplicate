package service

import (
	"duplicate/logging"
	"sync"

	"duplicate/internal/repository"
)

type Service struct {
	logger logging.Logger
	repo   repository.Repo
	mutex  *sync.RWMutex
}

func New(repo *repository.Repo, logger logging.Logger, mutex *sync.RWMutex) *Service {
	return &Service{logger: logger, repo: *repo, mutex: mutex}
}
