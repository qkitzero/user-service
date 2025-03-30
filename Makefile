MOCKGEN_CMD=go run go.uber.org/mock/mockgen@v0.5.0

mockgen:
	$(MOCKGEN_CMD) -source=internal/domain/user/user.go -destination=mocks/domain/user/mock_user.go -package=mocks
	$(MOCKGEN_CMD) -source=internal/domain/user/repository.go -destination=mocks/domain/user/mock_repository.go -package=mocks
	$(MOCKGEN_CMD) -source=internal/application/user/service.go -destination=mocks/application/user/mock_service.go -package=mocks

test:
	mkdir -p tmp
	go test -cover ./internal/domain/... -coverprofile=./tmp/cover.out
	go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	open ./tmp/cover.html