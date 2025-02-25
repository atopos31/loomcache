fmt:
	@go fmt ./...

test:
	@go clean -testcache
	@go test -v ./...

build:
	@go build -o ./bin/server .

tidy:
	@go mod tidy