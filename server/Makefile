default:
	go run *.go

dist:
	goxc -+tasks=bump,interpolate-source
	git commit -am "new minor release"
