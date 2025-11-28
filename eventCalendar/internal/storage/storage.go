package storage

import (
	"context"
	"eventCalendar/internal/config"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	DB *pgx.Conn
}

func New(confDb *config.DBConfig) *Storage {
	conn, err := pgx.Connect(context.Background(), confDb.URL)
	if err != nil {
		panic(err)
	}
	return &Storage{DB: conn}
}
