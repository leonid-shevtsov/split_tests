release: release-linux release-osx

release-linux:
	$(call release-base,linux)

release-osx:
	$(call release-base,darwin)


release-base = \
	mkdir -p build; \
	GOOS=$(1) GOARCH=amd64 go build -ldflags="-s -w" -o build/split_tests; \
	gzip -S .$(1).gz build/split_tests

clean:
	rm -rf build
