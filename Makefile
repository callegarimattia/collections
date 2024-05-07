lint:
	golangci-lint run

fmt: |
	goimports -w .
	gofumpt -l -w .
	golines -w -m 80 .

test:
	go test -v -cover ./...