bin = dq

.PHONY: all
all: test $(bin)

$(bin):
	go build

.PHONY: test
test:
	go test

.PHONY: clean
clean:
	rm -f $(bin)
