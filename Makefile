all: run

build:
	go generate
	go build

run: build
	./gorum -testdata
