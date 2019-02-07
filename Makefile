.PHONY: build clean dep reset run setup

# build the application binaries
build:
	go build -o ./build/gvent-api github.com/jmckind/gvent-api/cmd/gvent-api

# clean out any build artifacts
clean:
	rm -fr ./build

# ensure dependencies are installed
dep:
	dep ensure -v && dep status

# reset the local project to a pristine state
reset: clean
	rm -fr vendor

# run the application locally
run: build
	LOG_LEVEL=debug ./build/gvent-api

# ensure local project is ready for development
setup: dep
	mkdir -p ./build
