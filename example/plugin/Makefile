.PHONY: example
example:
	make clean
	go install awalterschulze.org/go/goderive
	goderive ./...
	go test -v ./...
	go install awalterschulze.org/go/goderive/example/plugin/cmd/goderive-genreadme
	goderive-genreadme .

.PHONY: clean
clean:
	find . -type f -name "derived.gen.go" -exec rm -f {} \;