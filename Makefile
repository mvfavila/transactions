# variables

CONTAINER_NAME=transactions
PORT = 8080

# vendoring

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

# Testing

.PHONY: test
test:
	go test -v ./...

.PHONY: test-cover
test-cover:
	go test -cover ./...

.PHONY: test-cover-visual
test-cover-visual:
	go test ./... -coverprofile=temp/coverage.out
	go tool cover -html=temp/coverage.out

# Linter

.PHONY: lint
lint:
	go vet
	golangci-lint run --disable errcheck

# Exec

.PHONY: run
run:
	APP_ENV=dev go run .

# Exec - docker

.PHONY: run-docker
run-docker:
	docker run --rm -it \
	-e APP_ENV='dev' \
	-p $(PORT):$(PORT) \
	$(CONTAINER_NAME)

# Build

.PHONY: build	
build:
	docker build --no-cache -t $(CONTAINER_NAME) .
