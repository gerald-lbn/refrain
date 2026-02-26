.PHONY: run test coverage view-coverage build clean lint

run:
	go run cmd/refrain/main.go

test:
	go test ./...

coverage:
	go test -cover -coverprofile=coverage.out ./...

view-coverage: coverage
	go tool cover -html=coverage.out

build:
	go build -o bin/refrain cmd/refrain/main.go

clean:
	rm -rf bin/
	rm -f coverage.out

lint:
	golangci-lint run
