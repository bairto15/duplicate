package handler

import (
	"duplicate/internal/service"
	"duplicate/logging"
)

type Handler struct {
	service service.Service
	logger  logging.Logger
}

func New(service *service.Service, logger logging.Logger) *Handler {
	logger.Info("Сервер запущен!")
	return &Handler{service: *service, logger: logger}
}

type result struct {
	Dupes bool `json:"dupes"`
}

type testCase struct {
	message  string
	user_id1 int32
	user_id2 int32
	expect   bool
}
