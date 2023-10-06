DB_NAME ?= banco_dos_amigos
DB_USER ?= postgres
DB_PASS ?= mysecretpassword
DB_HOST ?= localhost
DB_PORT ?= 5432
# postgres:mysecretpassword@database/indaband_local
DB_URL ?= postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
export

GOTEST ?= $(if $(shell which gotestsum),gotestsum --format testname --,go test)
integration:
	$(GOTEST) -v -tags integration -count 1 -run $(or '$(TESTS_RX)',.) $(or $(TESTS_PKG),./...)

test:
	$(GOTEST) -v -count 1 -run $(or '$(TESTS_RX)',.) $(or $(TESTS_PKG),./...)

migration:
	goose -allow-missing -dir ./sql/migrations postgres "user=$(DB_USER) password=$(DB_PASS) host=$(DB_HOST) port=$(DB_PORT) dbname=$(DB_NAME) sslmode=disable" up

run:
	go run cmd/server/main.go

