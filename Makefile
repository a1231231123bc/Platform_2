GO := /usr/local/go/bin/go
APP_NAME := platform-server
BUILD_DIR := ./bin

.PHONY: run build test test-cover lint migrate-up migrate-down clean tidy docker-verify

run:
	$(GO) run ./cmd/server/

build:
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/

test:
	$(GO) test -v ./...

test-cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

tidy:
	$(GO) mod tidy

docker-verify:
	@set +e; \
	docker compose down -v --remove-orphans >/dev/null 2>&1; \
	docker compose up --build --abort-on-container-exit --exit-code-from verifier verifier; \
	status=$$?; \
	docker compose down -v --remove-orphans; \
	exit $$status
