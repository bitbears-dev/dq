ver ?= 0.0.0
bin = dq
sources = $(wildcard *.go) $(wildcard ../../cli/*.go) $(wildcard ../../builtin/*.go)

.PHONY: all
all: test $(bin)

.PHONY: build-for-release
build-for-release: test $(bin) # the dependency $(bin) is required to build the binary to run on the local env for e2e-test
	GOOS=linux   GOARCH=amd64   go build -o dist/$(ver)/linux-amd64/dq
	GOOS=linux   GOARCH=arm64   go build -o dist/$(ver)/linux-arm64/dq
	GOOS=linux   GOARCH=arm     go build -o dist/$(ver)/linux-arm/dq
	GOOS=linux   GOARCH=riscv64 go build -o dist/$(ver)/linux-riscv64/dq
	GOOS=darwin  GOARCH=arm64   go build -o dist/$(ver)/darwin-arm64/dq
	GOOS=darwin  GOARCH=amd64   go build -o dist/$(ver)/darwin-amd64/dq
	GOOS=windows GOARCH=amd64   go build -o dist/$(ver)/windows-amd64/dq

.PHONY: package
package:
	cp ../../README.md ../../LICENSE dist/$(ver)/linux-amd64/
	cp ../../README.md ../../LICENSE dist/$(ver)/linux-arm64/
	cp ../../README.md ../../LICENSE dist/$(ver)/linux-arm/
	cp ../../README.md ../../LICENSE dist/$(ver)/linux-riscv64/
	cp ../../README.md ../../LICENSE dist/$(ver)/darwin-arm64/
	cp ../../README.md ../../LICENSE dist/$(ver)/darwin-amd64/
	cp ../../README.md ../../LICENSE dist/$(ver)/windows-amd64/

	mkdir -p dist/$(ver)/pkg

	tar -C dist/$(ver)/linux-amd64   -zcf dist/$(ver)/pkg/dq_$(ver)_linux_amd64.tar.gz   dq README.md LICENSE
	tar -C dist/$(ver)/linux-arm64   -zcf dist/$(ver)/pkg/dq_$(ver)_linux_arm64.tar.gz   dq README.md LICENSE
	tar -C dist/$(ver)/linux-arm     -zcf dist/$(ver)/pkg/dq_$(ver)_linux_arm.tar.gz     dq README.md LICENSE
	tar -C dist/$(ver)/linux-riscv64 -zcf dist/$(ver)/pkg/dq_$(ver)_linux_riscv64.tar.gz dq README.md LICENSE

	(cd dist/$(ver)/darwin-arm64 && zip -q -u ../pkg/dq_$(ver)_darwin_arm64.zip dq README.md LICENSE)
	(cd dist/$(ver)/darwin-amd64 && zip -q -u ../pkg/dq_$(ver)_darwin_amd64.zip dq README.md LICENSE)

	(cd dist/$(ver)/windows-amd64 && zip -q -u ../pkg/dq_$(ver)_windows_amd64.zip dq README.md LICENSE)

.PHONY: release
release:
	ghr -prerelease -replace -owner bitbears-dev -repository dq v$(ver) dist/$(ver)/pkg/

$(bin): $(sources)
	go build

.PHONY: test
test:
	go test

.PHONY: clean
clean:
	rm -f $(bin)
	rm -rf dist/
