.PHONY: tidy lint migrate seed up drop-table

tidy:
	go mod tidy
	cd cmd && go mod tidy

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


