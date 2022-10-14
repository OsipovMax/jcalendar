MIGRATE_DRV=postgres
POSTGRES_HOST=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=service_db
POSTGRES_PORT=5432
MIGRATE_DSN=host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} port=${POSTGRES_PORT} sslmode=disable TimeZone=Europe/Moscow
MIGRATE_DIR=./migrations/migrate
GOOSE_BASE_CMD=${GOBIN}/goose -dir ${MIGRATE_DIR} ${MIGRATE_DRV} "${MIGRATE_DSN}"

install:
	go install github.com/pressly/goose/v3/cmd/goose@latest


migration-up: install
	${GOOSE_BASE_CMD} up