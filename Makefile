release: release-linux release-osx

release-linux:
	$(call release-base,linux)

release-osx:
	$(call release-base,darwin)


release-base = \
	mkdir -p build/$(1); \
	GOOS=$(1) GOARCH=amd64 go build -ldflags="-s -w" -o build/$(1)/split_tests; \
	tar czf build/split_tests-$(1)-$(shell git rev-parse --short HEAD).tar.gz --directory=build/$(1) .

clean:
	rm -rf build
