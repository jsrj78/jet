VERS = $(shell git describe)
DATE = $(shell date +%F)

GO_BUILDINFO = -ldflags "-X main.vers=$(VERS) -X main.date=$(DATE)"

# not used yet
test:
	go test

# build and run the hub locally
run: hub
	./hub

# build the hub for local use
hub:
	go build $(GO_BUILDINFO) -o $@

# build the hub and report its full git tag
version: hub
	jet

# build the hub for several platforms and compress for release
builds:
	@ rm -rf $@; mkdir $@
	@ (echo JET/Hub $(VERS) $(DATE); go version) >$@/README.txt
	@ echo; cat $@/README.txt; echo
	@ echo Re-building binaries:
	@ echo "  MacOSX 64-bit"
	@ GOOS=darwin GOARCH=amd64 go build $(GO_BUILDINFO) -o $@/hub-darwin-amd64
	@ echo "  ARMv6 32-bit"
	@ GOOS=linux GOARCH=arm GOARM=6 go build $(GO_BUILDINFO) -o $@/hub-linux-arm
	@ echo "  Linux 64-bit"
	@ GOOS=linux GOARCH=amd64 go build $(GO_BUILDINFO) -o $@/hub-linux-amd64
	@ gzip $@/hub-*
	@ echo; ls -l builds/*; echo

clean:
	rm -f hub builds

.PHONY: test run hub version builds clean

# vim: set noexpandtab :
