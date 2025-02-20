fmt:
	@go fmt ./...

test:
	@go test -v ./...

tidy:
	@go mod tidy