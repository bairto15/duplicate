package repository

import (
	"duplicate/logging"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	Store
}

func New(db *sqlx.DB, logger logging.Logger, mutex *sync.RWMutex) *Repo {	
	return &Repo{
		Store: NewPostgres(db, logger, mutex),
	}
}