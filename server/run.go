package server

import (
	"duplicate/internal/handler"
	"duplicate/internal/repository"
	"duplicate/internal/service"
	"duplicate/logging"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/julienschmidt/httprouter"
)

func Run() error {
	logger := logging.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic: %v,\n%s", r, debug.Stack())
		}
	}()

	mutex := &sync.RWMutex{}

	postgres := repository.NewPostgresDB("postgres.sql", logger)
	repo := repository.New(postgres, logger, mutex)
	service := service.New(repo, logger, mutex)
	handler := handler.New(service, logger)

	router := httprouter.New()
	router.GET("/:user_id1/:user_id2", handler.ServeHTTP)
	http.ListenAndServe(":12345", router)

	return nil
}
