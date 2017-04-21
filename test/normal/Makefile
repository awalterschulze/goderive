.PHONY: derived.gen.go
derived.gen.go:
	goderive ./...

.PHONY: test
test: derived.gen.go gofmt
	go test -v ./...

.PHONY: gofmt
gofmt:
	gofmt -d .
	! gofmt -d . 2>&1 | read
