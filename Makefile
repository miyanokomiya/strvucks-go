GOOSE_CONNECTION := "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${DB_HOST} sslmode=disable"
GOOSE_DRIVER := "postgres"
GOOSE_DIR := "db/goose"

rsa:
	openssl genrsa -out jwt.rsa 2048
	openssl rsa -in jwt.rsa -pubout > jwt.rsa.pub
	base64 jwt.rsa
	base64 jwt.rsa.pub

goose:
	goose -dir ${GOOSE_DIR} ${GOOSE_DRIVER} ${GOOSE_CONNECTION} ${ARG}

goose/status:
	goose -dir ${GOOSE_DIR} ${GOOSE_DRIVER} ${GOOSE_CONNECTION} status

goose/up:
	goose -dir ${GOOSE_DIR} ${GOOSE_DRIVER} ${GOOSE_CONNECTION} up

goose/create:
	goose -dir ${GOOSE_DIR} ${GOOSE_DRIVER} ${GOOSE_CONNECTION} create ${ARG} sql

test:
	POSTGRES_DB=${POSTGRES_DB}_test make goose/up
	POSTGRES_DB=${POSTGRES_DB}_test go test ./...
