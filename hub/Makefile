VERS = $(shell git describe)
DATE = $(shell date +%F)

GO_BUILDINFO = -ldflags "-X main.vers=$(VERS) -X main.date=$(DATE)"

test:
	go test

run: hub
	./hub

hub:
	go build $(GO_BUILDINFO) -o $@

hub-arm:
	GOARCH=arm GOOS=linux go build $(GO_BUILDINFO) -o $@

clean:
	rm -f hub hub-arm

.PHONY: test run hub hub-arm clean

# vim: set noexpandtab :
