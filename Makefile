.PHONY: test
test:
	go install .
	make -C test test
	make -C example example

.PHONY: gofmt
gofmt:
	go fmt $(go list ./... | grep -v /vendor/)

.PHONY: action
action: gofmt
	go version
	make test
	go vet ./derive/...
	go vet ./example/...
	go vet ./plugin/...
	go vet ./test/normal/...
	make diff

updatedeps:
	go mod tidy
	go mod vendor
	git checkout vendor/vendortest/vendortest.go

diff:
	git diff --exit-code .
