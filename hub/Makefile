VERS = $(shell git describe)
DATE = $(shell date +%F)

GO_BUILDINFO = -ldflags "-X main.vers=$(VERS) -X main.date=$(DATE)"

test:
	go test

run: hub
	./hub

hub:
	go build $(GO_BUILDINFO) -o $@

builds:
	@ rm -rf $@; mkdir $@
	@ (echo JET/Hub $(VERS) $(DATE); go version) >$@/README.txt
	@ echo; cat $@/README.txt; echo
	@ echo Re-building binaries:
	@ echo "  MacOSX 64-bit"
	@ GOOS=darwin GOARCH=amd64 go build $(GO_BUILDINFO) -o $@/hub-darwin-amd64
	@ echo "  ARMv5 32-bit"
	@ GOOS=linux GOARCH=arm GOARM=5 go build $(GO_BUILDINFO) -o $@/hub-linux-arm
	@ echo "  Linux 64-bit"
	@ GOOS=linux GOARCH=amd64 go build $(GO_BUILDINFO) -o $@/hub-linux-amd64
	@ gzip $@/hub-*
	@ echo; ls -l builds/*; echo

clean:
	rm -f hub builds

.PHONY: test run hub builds clean

# vim: set noexpandtab :
