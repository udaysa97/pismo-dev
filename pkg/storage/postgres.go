package storage

import (
	// Postgres package

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (s *Store) CreateGormPostgresClient(postgresURL string) {

	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	s.GormPsql = db
}
