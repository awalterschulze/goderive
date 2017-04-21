.PHONY: test
test:
	rm derived.gen.go || true
	rm dedup_test.go || true
	cp dedup_test.gold dedup_test.go
	goderive -dedup .
	go test -v ./...
	rm derived.gen.go
	rm dedup_test.go
