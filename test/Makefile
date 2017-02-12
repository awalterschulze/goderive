derived.gen.go: structs.go
	goderive ./...

.PHONY: test
test: derived.gen.go gofmt
	go test -v ./...

gofmt:
	gofmt -l -s -w .
