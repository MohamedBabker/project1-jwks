
.PHONY: build run test cover lint fmt vet tidy

build:
	go build -o project1-jwks ./...

run:
	go run ./...

test:
	go test ./... -v

cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic -v
	go tool cover -func=coverage.out

lint: vet fmt

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
