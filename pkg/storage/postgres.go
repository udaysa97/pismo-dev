package storage

import (
	"github.com/lib/pq" // Postgres package

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (s *Store) CreatePostgresClient(postgresURL string, sslMode bool) {
	if !sslMode {
		postgresURL += "?sslmode=disable"
	}

	sqltrace.Register("pismo-dev-postgres", &pq.Driver{})
	db, err := sqltrace.Open("postgres", postgresURL)
	if err != nil {
		panic(err)
	}

	// maxOpenConnections, err := strconv.Atoi(appconfig.POSTGRES_MAX_CONNECTIONS)
	// if err != nil {
	// 	panic("MaxOpenConnections not found in environment")
	// }
	// maxIdleConnections, err := strconv.Atoi(appconfig.POSTGRES_MAX_IDLE_CONNECTIONS)
	// if err != nil {
	// 	panic("MaxIdleConnections not found in environment")
	// }
	// maxIdleTime, err := strconv.Atoi(appconfig.POSTGRES_MAX_IDLE_CONNECTION_TIME)
	// if err != nil {
	// 	panic("MaxIdleConnectionTime not found in environment")
	// }

	// db.SetMaxOpenConns(maxOpenConnections)
	// db.SetMaxIdleConns(maxIdleConnections)
	// db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	s.Postgres = db
}

func (s *Store) CreateGormPostgresClient(postgresURL string) {
	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	s.GormPsql = db
}
