-include .env

.PHONY: test lint lint-fix \
	proto-format proto-lint proto-gen \
	mock-gen \
	migrate-up migrate-up-one migrate-down migrate-reset migrate-create migrate-status

test:
	mkdir -p tmp
	go test -cover ./internal/... -coverprofile=./tmp/cover.out
	go tool cover -func=./tmp/cover.out | tail -n 1
	go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	open ./tmp/cover.html

lint:
	go tool golangci-lint run ./...

lint-fix:
	go tool golangci-lint run --fix ./...

proto-format:
	buf format -w

proto-lint:
	buf lint

proto-gen:
	buf generate

mock-gen:
	go tool mockgen -source=internal/domain/user/user.go -destination=mocks/domain/user/mock_user.go -package=mocks
	go tool mockgen -source=internal/domain/user/repository.go -destination=mocks/domain/user/mock_repository.go -package=mocks
	go tool mockgen -source=internal/domain/identity/identity.go -destination=mocks/domain/identity/mock_identity.go -package=mocks
	go tool mockgen -source=internal/application/user/usecase.go -destination=mocks/application/user/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/application/auth/service.go -destination=mocks/application/auth/mock_service.go -package=mocks
	go tool mockgen -destination=mocks/external/auth/v1/mock_client.go -package=mocks github.com/qkitzero/auth-service/gen/go/auth/v1 AuthServiceClient

MIGRATIONS_DIR=internal/infrastructure/db/migrations
MIGRATE=migrate -source file://$(MIGRATIONS_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_HOST_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

migrate-up:
	$(MIGRATE) up

migrate-up-one:
	$(MIGRATE) up 1

migrate-down:
	$(MIGRATE) down 1

migrate-reset:
	$(MIGRATE) drop -f

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -format 20060102150405 $(name)

migrate-status:
	$(MIGRATE) version
