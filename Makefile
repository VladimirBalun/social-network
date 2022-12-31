CUR_DIR=$(shell pwd)
BIN_DIR=${CUR_DIR}/bin

.PHONY: build
build:
	go build -o ${BIN_DIR}/profiles_stat ${CUR_DIR}/cmd/profiles_stat

.PHONY: test
test:
	go test -race ./... -v -cover

.PHONY: lint
lint:
	golangci-lint run

.PHONY: .install-goose
.install-goose:
	go install github.com/pressly/goose/cmd/goose@latest

.PHONY: migrate-up
migrate-up: .install-goose
	goose -dir ./migrations/mysql mysql "${DB_USER}:${DB_PASSWORD}@/${DB_NAME}?parseTime=true" up

.PHONY: migrate-down
migrate-down: .install-goose
	goose -dir ./migrations/mysql "user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} sslmode=disable" down
