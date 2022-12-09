DATASINK_DOCKER_REPO=krishnaiyer/datasink
DATASINK_VERSION=0.1.0
DATASINK_GIT_COMMIT=$(shell git rev-parse --short HEAD)
DATASINK_DATE=$(shell date)
DATASINK_PACKAGE="krishnaiyer.dev/golang/datasink"

.PHONY: init

init:
	@echo "Initialize development environment..."
	@mkdir -p .dev
	@go mod tidy

.PHONY: deps

deps:
	@echo "Install dependencies..."
	@go mod tidy

.PHONY: clean

clean:
	@echo "Clean development files..."
	@rm -rf .dev

.PHONY: build

build.local:
	go build \
	-ldflags="-X '${DATASINK_PACKAGE}/cmd.version=${DATASINK_VERSION}' \
	-X '${DATASINK_PACKAGE}/cmd.gitCommit=${DATASINK_GIT_COMMIT}' \
	-X '${DATASINK_PACKAGE}/cmd.buildDate=${DATASINK_DATE}'" \
	-o datasink main.go

build.docker.binary:
	GOOS=linux GOARCH=amd64 go build \
	-ldflags="-X '${DATASINK_PACKAGE}/cmd.version=${DATASINK_VERSION}' \
	-X '${DATASINK_PACKAGE}/cmd.gitCommit=${DATASINK_GIT_COMMIT}' \
	-X '${DATASINK_PACKAGE}/cmd.buildDate=${DATASINK_DATE}'" \
	-o datasink-docker main.go

docker.build: build.docker.binary
	docker build -t ${DATASINK_DOCKER_REPO}:${DATASINK_VERSION} .

docker.push:
	docker push ${DATASINK_DOCKER_REPO}:${DATASINK_VERSION}
	docker tag ${DATASINK_DOCKER_REPO}:${DATASINK_VERSION} ${DATASINK_DOCKER_REPO}:latest
	docker push ${DATASINK_DOCKER_REPO}:latest

build.dist:
	goreleaser --snapshot --skip-publish --rm-dist

clean.dist:
	rm -rf dist
