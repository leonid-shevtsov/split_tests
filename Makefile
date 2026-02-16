release: release-linux release-osx release-linux-arm

release-linux:
	$(call release-base,linux,amd64)

release-osx:
	$(call release-base,darwin,amd64)

release-linux-arm:
	$(call release-base,linux,arm64)


release-base = \
	mkdir -p build; \
	GOOS=$(1) GOARCH=$(2) go build -ldflags="-s -w" -o build/split_tests; \
	gzip -S .$(1)$(if $(filter arm64,$(2)),.$(2)).gz build/split_tests

clean:
	rm -rf build
