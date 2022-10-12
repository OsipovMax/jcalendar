MIGRATE_DRV=postgres
POSTGRES_HOST?=0.0.0.0
MIGRATE_DSN?=${PG_DSN}
MIGRATE_DSN?=host=${POSTGRES_HOST} user=postgres password=sbermarket_paas dbname=paas_db port=5432 sslmode=disable
MIGRATE_DIR=./migrations/migrate
GOOSE_BASE_CMD=${GOBIN}/goose -dir ${MIGRATE_DIR} ${MIGRATE_DRV} "${MIGRATE_DSN}"

migration-up:		## Migrate the DB to the most recent version available
	${GOOSE_BASE_CMD} up