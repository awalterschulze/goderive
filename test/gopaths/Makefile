.PHONY: test
test:
	GO111MODULE=off GOPATH=$(PWD)/gopath1:$(PWD)/gopath2 goderive ./...
	GO111MODULE=off GOPATH=$(PWD)/gopath1:$(PWD)/gopath2 go test -v ./...