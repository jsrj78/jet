test:
	go test .

run:
	go run *.go

hub:
	go build -ldflags "-X main.version=`git describe`"

hub-arm:
	GOARCH=arm GOOS=linux \
	go build -ldflags "-X main.version=`git describe`" -o $@

clean:
	rm -f hub hub-arm

.PHONY: test run hub hub-arm clean

# vim: set noexpandtab :
