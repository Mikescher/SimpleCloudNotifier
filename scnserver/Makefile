DOCKER_REPO=registry.blackforestbytes.com
DOCKER_NAME=mikescher/simplecloudnotifier
PORT=9090

NAMESPACE=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git rev-parse HEAD)

.PHONY: test swagger pygmentize docker migrate dgi pygmentize lint docker

SWAGGO_VERSION=v1.8.12
SWAGGO=github.com/swaggo/swag/cmd/swag@$(SWAGGO_VERSION)

build: ids enums swagger pygmentize fmt
	mkdir -p _build
	rm -f ./_build/scn_backend
	CGO_ENABLED=1 go build -v -o _build/scn_backend -tags "timetzdata sqlite_fts5 sqlite_foreign_keys" ./cmd/scnserver

enums:
	go generate models/enums.go

ids:
	go generate models/ids.go

run: build
	mkdir -p .run-data
	_build/scn_backend

gow:
	which gow || go install github.com/mitranim/gow@latest
	gow -e "go,mod,html,css,json,yaml,js" run -tags "timetzdata sqlite_fts5 sqlite_foreign_keys" blackforestbytes.com/simplecloudnotifier/cmd/scnserver

dgi:
	[ ! -f "DOCKER_GIT_INFO" ] || rm DOCKER_GIT_INFO
	echo -n "VCSTYPE="     >> DOCKER_GIT_INFO ; echo "git"                         >> DOCKER_GIT_INFO
	echo -n "BRANCH="      >> DOCKER_GIT_INFO ; git rev-parse --abbrev-ref HEAD    >> DOCKER_GIT_INFO
	echo -n "HASH="        >> DOCKER_GIT_INFO ; git rev-parse              HEAD    >> DOCKER_GIT_INFO
	echo -n "COMMITTIME="  >> DOCKER_GIT_INFO ; git log -1 --format=%cd --date=iso >> DOCKER_GIT_INFO
	echo -n "REMOTE="      >> DOCKER_GIT_INFO ; git config --get remote.origin.url >> DOCKER_GIT_INFO

docker: dgi
	cp ../scn_send.sh .
	docker build \
            -t "$(DOCKER_NAME):$(HASH)" \
            -t "$(DOCKER_NAME):$(NAMESPACE)-latest" \
            -t "$(DOCKER_NAME):latest" \
            -t "$(DOCKER_REPO)/$(DOCKER_NAME):$(HASH)" \
            -t "$(DOCKER_REPO)/$(DOCKER_NAME):$(NAMESPACE)-latest" \
            -t "$(DOCKER_REPO)/$(DOCKER_NAME):latest" \
            .
	[ -f "scn_send.sh" ] && rm scn_send.sh

swagger-setup:
	mkdir -p ".swaggobin"
	[ -f ".swaggobin/swag_$(SWAGGO_VERSION)" ] || { GOBIN=/tmp/_swaggo go install $(SWAGGO); cp "/tmp/_swaggo/swag" ".swaggobin/swag_$(SWAGGO_VERSION)"; rm -rf "/tmp/_swaggo"; }

swagger: swagger-setup
	".swaggobin/swag_$(SWAGGO_VERSION)" init -generalInfo ./api/router.go --propertyStrategy camelcase --output ./swagger/ --outputTypes "json,yaml"

pygmentize: website/scn_send.html

website/scn_send.html: ../scn_send.sh
	_pygments/pygmentizew -l bash    -f html "$(shell pwd)/../scn_send.sh" > "$(shell pwd)/website/scn_send.html"
	_pygments/pygmentizew -S monokai -f html                               > "$(shell pwd)/website/css/pygmnetize-dark.css"
	_pygments/pygmentizew -S borland -f html                               > "$(shell pwd)/website/css/pygmnetize-light.css"

run-docker-local: docker
	mkdir -p .run-data
	docker run --rm \
	           --init \
	           --env "CONF_NS=local-docker" \
			   --volume "$(shell pwd)/.run-data/docker-local:/data" \
			   --publish "8080:80" \
			   $(DOCKER_NAME):latest

inspect-docker: docker
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
	rm -rf _pygments/env
	git clean -fdx
	! which go 2>&1 >> /dev/null || go clean
	! which go 2>&1 >> /dev/null || go clean -testcache

fmt: swagger-setup
	go fmt ./...
	".swaggobin/swag_$(SWAGGO_VERSION)" fmt

test:
	which gotestsum || go install gotest.tools/gotestsum@latest
	gotestsum --format "testname" -- -tags="timetzdata sqlite_fts5 sqlite_foreign_keys" "./test"

migrate:
	CGO_ENABLED=1 go build -v -o _build/scn_migrate -tags "timetzdata sqlite_fts5 sqlite_foreign_keys" ./cmd/migrate
	./_build/scn_migrate

lint:
	# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.2
	golangci-lint run ./...


