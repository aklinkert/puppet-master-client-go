.PHONY: vendors
vendors:
	go mod download
	go mod tidy

.PHONY: test
test:
	go test -cover ./...
	golangci-lint run
	golint -set_exit_status ./...
