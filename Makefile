.PHONY: test
test:
	go install .
	goderive ./...
	make gofmt
	go test -v ./...

.PHONY: gofmt
gofmt:
	gofmt -l -s -w .