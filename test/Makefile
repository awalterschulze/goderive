.PHONY: derived.gen.go
derived.gen.go:
	goderive ./...

.PHONY: test
test: derived.gen.go gofmt
	go test -v ./...

gofmt:
	gofmt -l -s -w .
