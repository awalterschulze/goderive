.PHONY: example
example:
	make clean
	go install go.awalterschulze.org/goderive
	goderive ./...
	go test -v ./...
	go install go.awalterschulze.org/goderive/example/plugin/cmd/goderive-genreadme
	goderive-genreadme .

.PHONY: clean
clean:
	find . -type f -name "derived.gen.go" -exec rm -f {} \;