.PHONY: build clean dep run setup

# build the application binaries
build: setup
	go build -o ./build/gvent-api github.com/jmckind/gvent-api/cmd

# clean out any build artifacts
clean:
	rm -fr ./build

# ensure dependencies are installed
dep:
	dep ensure && dep status

# reset the local project to a pristine state
reset: clean
	rm -fr vendor

# run the application locally
run: build
	LOG_LEVEL=debug ./build/gvent-api

# ensure build directories and dependencies are present
setup: dep
	mkdir -p ./build
