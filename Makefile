current_dir := $(abspath $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))

.PHONY: build
build:
	go build ./...

.PHONY: run
run:
	@docker build . -t chatserver
	@docker run -p 5050:5050 -p 8080:8080 --name=chatserver --rm --mount type=bind,source=$(current_dir)/logs,target=/logs  chatserver
