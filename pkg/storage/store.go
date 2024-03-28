package storage

import (
	"database/sql"
	"gorm.io/gorm"
)

type Store struct {
	Postgres *sql.DB
	GormPsql *gorm.DB
}

func (s *Store) InitPostgresClient(databaseUrl string) {
	if len(databaseUrl) == 0 {
		panic("env variable DATABASE_URL is missing")
	}
	s.CreatePostgresClient(databaseUrl, false)
	s.CreateGormPostgresClient(databaseUrl)
}
