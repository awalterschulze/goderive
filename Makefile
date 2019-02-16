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
	go version
	make test
	go vet ./derive/...
	go vet ./example/...
	go vet ./plugin/...
	go vet ./test/normal/...
	make diff

updatedeps:
	govendor fetch +vendor
	git checkout vendor/vendortest/vendortest.go

diff:
	git diff --exit-code .