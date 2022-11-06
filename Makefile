subdirs = builtin cli cmd/dq

.PHONY: all
all: lint test $(subdirs)

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: clean test
clean test: $(subdirs)

$(subdirs): force
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: force
force:

.PHONY: update-deps
update-deps:
	go get -u ./...

.PHONY: e2e-test
e2e-test:
	TZ=Asia/Tokyo test/test.sh
	TZ=US/Pacific test/test.sh
