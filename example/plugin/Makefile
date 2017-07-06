.PHONY: example
example:
	go install github.com/awalterschulze/goderive
	goderive ./...
	go test -v ./...
	go install github.com/awalterschulze/goderive/example/plugin/cmd/goderive-genreadme
	goderive-genreadme ./equal ./compare