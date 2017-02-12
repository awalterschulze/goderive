.PHONY: test
test:
	go install .
	make -C test test

.PHONY: gofmt
gofmt:
	gofmt -l -s -w .