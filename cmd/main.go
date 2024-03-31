package main

import (
	"pismo-dev/api"
	"pismo-dev/initializer"
	"pismo-dev/internal/appconfig"
	"pismo-dev/pkg/logger"
	"pismo-dev/pkg/storage"
)

func init() {
	//if env := os.Getenv("ENV"); env == "DEVELOPMENT" {
	initializer.LoadEnvVariables()
	//	}

	appconfig.SetEnvVariables()
	logger.SetAppName("pismo-dev")
}

func main() {

	databaseUrl := appconfig.DATABASE_URL

	if len(databaseUrl) == 0 {
		panic("env variable DATABASE_URL is missing")
	}

	db := &storage.Store{}
	db.InitPostgresClient(databaseUrl)

	repositories := initializer.InitRepositories(db)
	services := initializer.InitServices(repositories)

	api.InitServer(services, db)
}
