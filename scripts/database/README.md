Steps to run migration on local
1. `brew install golang-migrate`
2. `migrate create -ext sql -dir database/migration -seq name`
3. use `migrateup` or `migratedown` commands from makefile to run migration scripts
4. `migrate -help` for more options