all: cqlmm

cqlmm: $(shell find . -type f -iname "*.go")
	go build -o $@ github.com/nesv/cqlmm/cmd/cqlmm

clean:
	rm -f cqlmm

install: all
	go install github.com/nesv/cqlmm/cmd/cqlmm

.PHONY: clean install
