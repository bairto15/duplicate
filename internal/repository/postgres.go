package repository

import (
	"duplicate/logging"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	logger logging.Logger
	db     *sqlx.DB
	Cash   Cash
	mutex  *sync.RWMutex
}


func NewPostgres(db *sqlx.DB, logger logging.Logger, mutex *sync.RWMutex) *Postgres {
	cash := Cash{UsersWithIps: make(UsersWithIps), IpsWithUsers: make(IpsWithUsers)}

	postgres := &Postgres{db: db, logger: logger, Cash: cash, mutex: mutex}

	logger.Info("Сохранение в кэш данные. Подождите ...")

	postgres.SaveLogInCash()

	logger.Info("Сохранение в кэш данные завершено!")

	return postgres
}

func NewPostgresDB(path string, logger logging.Logger) *sqlx.DB {
	cr := "user=postgres password=qwerty host=localhost port=5432 dbname=postgres sslmode=disable"

	db, err := sqlx.Connect("postgres", cr)
	if err != nil {
		logger.Fatal(err)
	}

	sqlFile, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Генерация рандомных данных. Подождите ...")

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Генерация данных завершено!")

	return db
}

type Cash struct {
	UsersWithIps UsersWithIps
	IpsWithUsers IpsWithUsers
}

type SetIps map[string]struct{}
type SetUserId map[int32]struct{}

type UsersWithIps map[int32]SetIps

type IpsWithUsers map[string]SetUserId

type log struct {
	UserId int32  `db:"user_id"`
	IpAddr string `db:"ip_addr"`
	Ts     string `db:"ts"`
}
