build:
	@go build .

test:
	@go test -v ./...

coverage:
	@go tool cover -html=coverage.out -o coverage.html