.PHONY: example
example:
	make -C 01_min example
	make -C 02_min example
	make -C 03_sort example
	make -C 04_range example
	make -C 05_prefix example
	make -C 06_mem example
	make -C 07_fib example
	make -C 08_hash example
	make -C 09_autoname example
	make -C 10_dup example

.PHONY: install
install:
	go install github.com/awalterschulze/goderive

.PHONY: clean
clean:
	rm 01_min/derived.gen.go || true
	rm 02_min/derived.gen.go || true
	rm 03_sort/derived.gen.go || true
	rm 04_range/derived.gen.go || true
	rm 05_prefix/derived.gen.go || true
	rm 06_mem/derived.gen.go || true
	rm 07_fib/derived.gen.go || true
	rm 08_hash/derived.gen.go || true
	rm 09_autoname/derived.gen.go || true
	rm 10_dup/derived.gen.go || true
