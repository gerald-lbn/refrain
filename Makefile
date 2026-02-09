.PHONY: run test build clean lint

run:
	go run cmd/refrain/main.go

test:
	ginkgo run ./...

coverage:
	ginkgo run -cover -coverprofile=coverage.out ./...

view-coverage: coverage
	go tool cover -html=coverage.out

build:
	go build -o bin/refrain cmd/refrain/main.go

clean:
	rm -rf bin/
	rm -f coverage.out

lint:
	golangci-lint run
