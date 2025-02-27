fmt:
	@go fmt ./...

test:
	@go clean -testcache
	@go test -v ./...

build:
	@go build -o ./bin/server .

tidy:
	@go mod tidy

vet:
	@go vet ./...

proto:
	@protoc --go_out=. ./proto/cache.proto