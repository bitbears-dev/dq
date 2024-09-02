subdirs = builtin cli cmd/dq

.PHONY: all
all: lint test $(subdirs)

.PHONY: build-for-release
build-for-release: lint test $(subdirs) e2e-test

.PHONY: lint
lint:
	docker run -it --rm devdrops/staticcheck:latest staticcheck

.PHONY: clean test package release
clean test package release: $(subdirs)

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
	TZ=America/Los_Angeles test/test.sh
