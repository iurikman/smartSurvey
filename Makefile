hello:
	echo "Hello my friend!"

run:
	go run cmd/service/main.go

lint:
	gofumpt -w .
	go mod tidy
	golangci-lint run --fix -c .golangci.yml ./...

test:
	make run &
	go test -v ./...
	fg
	PID=$!
	kill $PID