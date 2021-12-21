.PHONY: example
example:
	goderive --dedup=true --prefix="" ./...
	go run *.go

clean:
	rm derived.gen.go || true