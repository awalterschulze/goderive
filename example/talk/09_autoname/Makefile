.PHONY: example
example:
	goderive --autoname=true --prefix="" ./...
	go run *.go

clean:
	rm derived.gen.go || true