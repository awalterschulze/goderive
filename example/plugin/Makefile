.PHONY: example
example:
	make clean
	go install github.com/awalterschulze/goderive
	goderive ./...
	go test -v ./...
	go install github.com/awalterschulze/goderive/example/plugin/cmd/goderive-genreadme
	goderive-genreadme ./equal ./compare ./copyto ./keys ./sort ./unique ./set ./min ./max

.PHONY: clean
clean:
	rm ./equal/derived.gen.go || true
	rm ./compare/derived.gen.go || true
	rm ./copyto/derived.gen.go || true
	rm ./keys/derived.gen.go || true
	rm ./sort/derived.gen.go || true
	rm ./unique/derived.gen.go || true
	rm ./set/derived.gen.go || true
	rm ./min/derived.gen.go || true
	rm ./max/derived.gen.go || true