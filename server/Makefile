DOCKER_REPO=registry.blackforestbytes.com
DOCKER_NAME=scn_server
PORT=9090

NAMESPACE=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git rev-parse HEAD)

build:
	rm -f ./_build/scn_backend
	go build -v -o _build/scn_backend ./cmd/scnserver

run: build
	mkdir -p .run-data
	_build/scn_backend

build-docker:
	[ ! -f "DOCKER_GIT_INFO" ] || rm DOCKER_GIT_INFO
	git rev-parse --abbrev-ref HEAD    >> DOCKER_GIT_INFO
	git rev-parse              HEAD    >> DOCKER_GIT_INFO
	git log -1 --format=%cd --date=iso >> DOCKER_GIT_INFO
	git config --get remote.origin.url >> DOCKER_GIT_INFO
	docker build \
            -t "$(DOCKER_NAME):$(HASH)" \
            -t "$(DOCKER_NAME):$(NAMESPACE)-latest" \
            -t "$(DOCKER_NAME):latest" \
            -t "$(DOCKER_REPO)/$(DOCKER_NAME):$(HASH)" \
            -t "$(DOCKER_REPO)/$(DOCKER_NAME):$(NAMESPACE)-latest" \
            -t "$(DOCKER_REPO)/$(DOCKER_NAME):latest" \
            .

.PHONY: swagger
swagger:
	which swag || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -generalInfo api/router.go --output ./swagger/ --outputTypes "json,yaml"

run-docker-local:
	mkdir -p .run-data
	docker run --rm \
	           --init \
	           --env "CONF_NS=local" \
			   --volume "$(shell pwd)/.run-data/docker-local:/data" \
			   --network host \
			   $(DOCKER_NAME):latest

inspect-docker:
	mkdir -p .run-data
	docker run -ti \
	           --rm \
	           --volume "$(shell pwd)/.run-data/docker-inspect:/data" \
	           $(DOCKER_NAME):latest \
	           bash

push-docker:
	docker image push "$(DOCKER_REPO)/$(DOCKER_NAME):$(HASH)"
	docker image push "$(DOCKER_REPO)/$(DOCKER_NAME):$(NAMESPACE)-latest"
	docker image push "$(DOCKER_REPO)/$(DOCKER_NAME):latest"

clean:
	rm -rf _build/*
	rm -rf .run-data/*
	git clean -fdx
	go clean

fmt:
	go fmt ./...
	swag fmt
