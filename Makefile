.PHONY: check fmt vet lint test coverage proto

check:
	@echo "Running all checks..."
	@bash scripts/check.sh

fmt:
	@echo "Checking Go file formatting..."
	@bash scripts/checks/fmt.sh

vet:
	@echo "Checking for common Go mistakes..."
	@bash scripts/checks/vet.sh

lint:
	@echo "Checking for style errors..."
	@bash scripts/checks/lint.sh

test:
	@echo "Running unit tests..."
	@go test -race ./...

coverage:
	@echo "Generating coverage profile..."
	@go test -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./...)

proto:
	@echo "Generating proto files..."
	@protoc --go_out=. *.proto
