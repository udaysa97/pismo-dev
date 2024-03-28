migrateup:
    migrate -path database/migration -database "postgresql://coindcx:coindcx@localhost:5432/okto_bff?sslmode=disable" up
migratedown:
    migrate -path database/migration -database "postgresql://coindcx:coindcx@localhost:5432/okto_bff?sslmode=disable" down
.PHONY: migrateup migratedown