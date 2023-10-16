package handler

import (
	"duplicate/internal/repository"
	"duplicate/internal/service"
	"duplicate/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func Test_Duplicate(t *testing.T) {
	// Создание mock-объекта для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock-объекта: %s", err)
	}
	defer db.Close()

	//Инициализируем обьекты
	logger := logging.GetLogger()	
	mutex := &sync.RWMutex{}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.New(sqlxDB, logger, mutex)
	service := service.New(repo, logger, mutex)
	h := New(service, logger)

	// Определение ожидаемого запроса и его результатов
	mock.ExpectQuery("SELECT user_id, ip_addr FROM conn_log").
		WillReturnRows(sqlmock.
			NewRows([]string{"user_id", "ip_addr"}).
			AddRow(1, "127.0.0.1").
			AddRow(1, "127.0.0.20").
			AddRow(2, "127.0.0.2").
			AddRow(2, "127.0.0.2").
			AddRow(3, "127.0.0.2").
			AddRow(4, "127.0.0.4").
			AddRow(1, "127.0.0.1").
			AddRow(5, "127.0.0.5").
			AddRow(4564645, "127.0.120.55").
			AddRow(55466, "127.0.40.225").
			AddRow(79878456, "127.0.40.225").
			AddRow(798784561, "127.0.40.225").
			AddRow(455, "127.0.0.5").
			AddRow(5, "127.0.0.5").
			AddRow(7, "127.0.0.3").
			AddRow(8, "127.0.0.3").
			AddRow(50, "127.0.0.50").
			AddRow(100, "127.0.0.100"))

	err = repo.SaveLogInCash()
	if err != nil {
		logger.Error(err)
	}

	// Проверка результатов
	cases := []testCase{
		{
			message:  "Одинаковое user_id",
			user_id1: 1,
			user_id2: 1,
			expect:   true,
		},
		{
			message:  "Одинаковое ip адреса 1",
			user_id1: 2,
			user_id2: 3,
			expect:   true,
		},
		{
			message:  "Одинаковое ip адреса 2",
			user_id1: 55466,
			user_id2: 79878456,
			expect:   true,
		},
		{
			message:  "Одинаковое ip адреса 3",
			user_id1: 798784561,
			user_id2: 79878456,
			expect:   true,
		},
		{
			message:  "Разные ip адреса 1",
			user_id1: 1,
			user_id2: 3,
			expect:   false,
		},
		{
			message:  "Разные ip адреса 2",
			user_id1: 100,
			user_id2: 50,
			expect:   false,
		},
		{
			message:  "Разные ip адреса 3",
			user_id1: 5,
			user_id2: 7,
			expect:   false,
		},
		{
			message:  "Не существющий 1 user_id ",
			user_id1: 30,
			user_id2: 50,
			expect:   false,
		},
		{
			message:  "Не существющие 2 user_id ",
			user_id1: 30,
			user_id2: 80,
			expect:   false,
		},
	}

	router := httprouter.New()
	router.GET("/:user_id1/:user_id2", h.ServeHTTP)

	for _, v := range cases {
		url := fmt.Sprintf("/%d/%d", v.user_id1, v.user_id2)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		out, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
		}
		
		var res result
		json.Unmarshal(out, &res)
		
		assert.Equal(t, v.expect, res.Dupes, v.message)
		logger.Info(res)
	}
}

func BenchmarkDuplicate(b *testing.B) {
	// Создание mock-объекта для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("Ошибка при создании mock-объекта: %s", err)
	}
	defer db.Close()

	//Инициализируем обьекты
	logger := logging.GetLogger()
	mutex := &sync.RWMutex{}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.New(sqlxDB, logger, mutex)
	service := service.New(repo, logger, mutex)
	h := New(service, logger)

	// Определение ожидаемого запроса и его результатов
	mock.ExpectQuery("SELECT user_id, ip_addr FROM conn_log").
		WillReturnRows(sqlmock.
			NewRows([]string{"user_id", "ip_addr"}).
			AddRow(1, "127.0.0.1").
			AddRow(1, "127.0.0.20").
			AddRow(2, "127.0.0.2").
			AddRow(2, "127.0.0.2").
			AddRow(3, "127.0.0.2").
			AddRow(4, "127.0.0.4").
			AddRow(1, "127.0.0.1").
			AddRow(5, "127.0.0.5").
			AddRow(4564645, "127.0.120.55").
			AddRow(55466, "127.0.40.225").
			AddRow(79878456, "127.0.40.225").
			AddRow(798784561, "127.0.40.225").
			AddRow(455, "127.0.0.5").
			AddRow(5, "127.0.0.5").
			AddRow(7, "127.0.0.3").
			AddRow(8, "127.0.0.3").
			AddRow(50, "127.0.0.50").
			AddRow(100, "127.0.0.100"))

	err = repo.SaveLogInCash()
	if err != nil {
		logger.Error(err)
	}

	// Проверка результатов
	cases := []testCase{
		{
			message:  "Одинаковое user_id",
			user_id1: 1,
			user_id2: 1,
			expect:   true,
		},
		{
			message:  "Одинаковое ip адреса 1",
			user_id1: 2,
			user_id2: 3,
			expect:   true,
		},
		{
			message:  "Одинаковое ip адреса 2",
			user_id1: 55466,
			user_id2: 79878456,
			expect:   true,
		},
		{
			message:  "Одинаковое ip адреса 3",
			user_id1: 798784561,
			user_id2: 79878456,
			expect:   true,
		},
		{
			message:  "Разные ip адреса 1",
			user_id1: 1,
			user_id2: 3,
			expect:   false,
		},
		{
			message:  "Разные ip адреса 2",
			user_id1: 100,
			user_id2: 50,
			expect:   false,
		},
		{
			message:  "Разные ip адреса 3",
			user_id1: 5,
			user_id2: 7,
			expect:   false,
		},
		{
			message:  "Не существющий 1 user_id ",
			user_id1: 30,
			user_id2: 50,
			expect:   false,
		},
		{
			message:  "Не существющие 2 user_id ",
			user_id1: 30,
			user_id2: 80,
			expect:   false,
		},
	}
	
	router := httprouter.New()
	router.GET("/:user_id1/:user_id2", h.ServeHTTP)
	
	b.ResetTimer()

	for _, v := range cases {
		url := fmt.Sprintf("/%d/%d", v.user_id1, v.user_id2)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		out, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
		}
		
		var res result
		json.Unmarshal(out, &res)
		
		logger.Info(res)
	}
}
