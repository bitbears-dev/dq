bin = dq
sources = $(wildcard *.go) $(wildcard ../../cli/*.go) $(wildcard ../../builtin/*.go)

.PHONY: all
all: test $(bin)

$(bin): $(sources)
	go build

.PHONY: test
test:
	go test

.PHONY: clean
clean:
	rm -f $(bin)
