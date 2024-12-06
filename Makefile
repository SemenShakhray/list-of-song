include .env
LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0

migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v

migration-add:
	goose -dir ${MIGRATION_DIR} create ${MIGRATION_NAME} sql

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v
