default:
	mkdir -p ../build
	go build -o ../build/jethub .

dist:
	goxc -+tasks=bump,interpolate-source
	git commit -am "new minor release"
