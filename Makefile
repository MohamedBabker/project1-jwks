.PHONY: build run test cover lint fmt vet tidy

# Build the binary for the main package only
build:
	go build -o project1-jwks .

# Run the main package only
run:
	go run .

# Run all tests (verbose)
test:
	go test ./... -v

# Run tests with coverage across code packages (exclude main from denominator)
cover:
	go test ./... -coverpkg=./internal/...,./jwks -coverprofile=coverage.out -covermode=atomic -v
	go tool cover -func=coverage.out

# Lint helpers
lint: vet fmt

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
