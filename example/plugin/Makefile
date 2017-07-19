.PHONY: example
example:
	make clean
	go install github.com/awalterschulze/goderive
	goderive ./...
	go test -v ./...
	go install github.com/awalterschulze/goderive/example/plugin/cmd/goderive-genreadme
	goderive-genreadme .

.PHONY: clean
clean:
	find . -type f -name "derived.gen.go" -exec rm -f {} \;