.PHONY: example
example:
	go install github.com/awalterschulze/goderive
	rm ./equal/derived.gen.go || true
	rm ./compare/derived.gen.go || true
	rm ./copyto/derived.gen.go || true
	rm ./keys/derived.gen.go || true
	goderive ./...
	go test -v ./...
	go install github.com/awalterschulze/goderive/example/plugin/cmd/goderive-genreadme
	goderive-genreadme ./equal
	goderive-genreadme ./compare
	goderive-genreadme ./copyto
	goderive-genreadme ./keys