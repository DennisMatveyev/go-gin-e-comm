build:
	@go build -o bin/gin-e-comm

run: build
	@./bin/gin-e-comm

test:
	@go test -v ./...