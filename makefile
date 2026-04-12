build:
	go build -o bin/bldd cmd/bldd/main.go

ci:
	GOWORK=off CGO_ENABLED=1 go test -vet=all -race ./...

.PHONY: build ci
