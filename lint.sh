find . -name "*.go" -type f -exec goimports -w {} \;
find . -name "*.go" -type f -exec go fmt {} \;
golangci-lint run
go test ./...
