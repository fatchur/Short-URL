.PHONY: tidy lint migrate seed up drop-table clear-table mocks integration-test build-monolith build-user build-short-url up-monolith down-monolith up-user down-user up-short-url down-short-url up-db down-db

tidy:
	go mod tidy
	cd cmd && go mod tidy
	cd pkg/short-url && go mod tidy
	cd pkg/user && go mod tidy
	cd pkg && go mod tidy

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	elif [ -f "$(shell go env GOPATH)/bin/golangci-lint" ]; then \
		$(shell go env GOPATH)/bin/golangci-lint run ./...; \
	else \
		echo "golangci-lint not found; running 'go vet ./...'"; \
		go vet ./...; \
		echo "Tip: install golangci-lint: https://golangci-lint.run"; \
	fi

migrate:
	cd cmd && go run . -d migrate

seed:
	cd cmd && go run . -d=seed

up:
	docker-compose up -d

drop-table:
	cd cmd && go run . -d=drop-table

clear-table:
	cd cmd && go run . -d=clear-table

mocks:
	$(shell go env GOPATH)/bin/mockery

integration-test:
	@echo "Running integration tests with coverage..."
	@echo "Testing User Service Controller..."
	cd pkg/user && go test -v -cover ./api/controller -run TestUserControllerIntegrationTestSuite
	@echo "Testing Short URL Service Controller..."
	cd pkg/short-url && go test -v -cover ./api/controller -run TestShortUrlControllerIntegrationTestSuite
	@echo "All integration tests completed!"

build-monolith:
	docker build -t short-url-monolith -f pkg/Dockerfile .

build-user:
	docker build -t user-service -f pkg/user/Dockerfile .

build-short-url:
	docker build -t short-url-service -f pkg/short-url/Dockerfile .

up-monolith:
	docker-compose -f docker-compose.monolith.yml up -d

down-monolith:
	docker-compose -f docker-compose.monolith.yml down

up-user:
	docker-compose -f docker-compose.user.yml up -d

down-user:
	docker-compose -f docker-compose.user.yml down

up-short-url:
	docker-compose -f docker-compose.short-url.yml up -d

down-short-url:
	docker-compose -f docker-compose.short-url.yml down

up-db:
	docker-compose -f docker-compose.db.yml up -d

down-db:
	docker-compose -f docker-compose.db.yml down

