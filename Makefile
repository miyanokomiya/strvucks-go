GOOSE_CONNECTION := "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${DB_HOST} sslmode=disable"
GOOSE_DRIVER := "postgres"

goose:
	goose ${GOOSE_DRIVER} ${GOOSE_CONNECTION} ${ARG}

goose/status:
	goose ${GOOSE_DRIVER} ${GOOSE_CONNECTION} status

goose/up:
	goose ${GOOSE_DRIVER} ${GOOSE_CONNECTION} up

goose/create:
	goose ${GOOSE_DRIVER} ${GOOSE_CONNECTION} create ${ARG} sql
