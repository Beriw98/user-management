.PHONY: run

run:
	docker-compose up -d --build

stop:
	docker-compose down

lint:
	gofmt -w .
	goimports-reviser -format ./...
	go mod tidy
	GOARCH=arm64 golangci-lint run

show-coverage:
	go test $(go list ./... | grep -v "./ent") -coverprofile=coverage.out
	go tool cover -html=coverage.out
