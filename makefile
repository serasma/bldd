build:
	go build -o bin/goflowlens cmd/goflowlens/main.go

ci:
	GOWORK=off CGO_ENABLED=1 go test -vet=all -race ./...

.PHONY: build ci
