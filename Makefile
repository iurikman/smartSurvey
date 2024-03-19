hello:
	echo "Hello my friend!"

run:
	go run cmd/service/main.go

lint:
	gofumpt -w .
	go mod tidy
	golangci-lint run -c .golangci.yml ./...
