.PHONY: test
test:
	rm derived.gen.go || true
	rm autoname_test.go || true
	cp autoname_test.gold autoname_test.go
	goderive -autoname .
	go test -v ./...
	rm derived.gen.go
	rm autoname_test.go
