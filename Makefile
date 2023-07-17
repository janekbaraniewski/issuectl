.PHONY: test
test:
	go test ./...

issuectl:
	go build -o issuectl cmd/main.go
