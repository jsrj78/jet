VERS = $(shell git describe)
DATE = $(shell date +%F)

test:
	go test

run: hub
	./hub

hub:
	go build -ldflags "-X main.vers=$(VERS) -X main.date=$(DATE)" -o $@

hub-arm:
	GOARCH=arm GOOS=linux \
	go build -ldflags "-X main.vers=$(VERS) -X main.date=$(DATE)" -o $@

clean:
	rm -f hub hub-arm

.PHONY: test run hub hub-arm clean

# vim: set noexpandtab :
