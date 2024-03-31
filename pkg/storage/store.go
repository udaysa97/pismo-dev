package storage

import (
	"gorm.io/gorm"
)

type Store struct {
	GormPsql *gorm.DB
}

func (s *Store) InitPostgresClient(databaseUrl string) {
	if len(databaseUrl) == 0 {
		panic("env variable DATABASE_URL is missing")
	}
	s.CreateGormPostgresClient(databaseUrl)
}
