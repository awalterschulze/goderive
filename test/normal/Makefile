.PHONY: derived.gen.go
derived.gen.go:
	goderive ./...

.PHONY: clean
clean:
	rm gostring_gen_test.go || true

.PHONY: test
test: clean derived.gen.go gofmt
	go test -v ./...

.PHONY: gofmt
gofmt:
	gofmt -d .
	! gofmt -d . 2>&1 | read
