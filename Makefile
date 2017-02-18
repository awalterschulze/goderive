.PHONY: test
test:
	go install .
	make -C test test
	make -C example example

.PHONY: gofmt
gofmt:
	gofmt -l -s -w .

.PHONY: travis
travis:
	make test
