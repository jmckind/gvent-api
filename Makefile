.PHONY: build clean dep image run

build:
	mkdir -p ./build
	GOARCH=amd64 GOOS=linux go build -o ./build/gvent-api github.com/jmckind/gvent-api/cmd

clean:
	rm -fr ./build

dep:
	dep ensure && dep status

image:
	docker build -t jmckind/gvent-api:latest .

run: build
	LOG_LEVEL=debug ./build/gvent-api
