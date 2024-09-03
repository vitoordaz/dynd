ALL: clean lint test build

build:
	go build -mod=vendor -o build/dynd cmd/dynd/*.go

lint: gofmt goimports
	docker run --rm -e LOG_LEVEL=error -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

test:
	go test ./...

gofmt:
	gofmt -w .

goimports:
	goimports -local "github.com/vitoordaz/dynd" -w .

clean:
	rm -rf build

vendor:
	go mod vendor
