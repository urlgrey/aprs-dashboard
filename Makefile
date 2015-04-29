GO ?= godep go

all: build test

godep:
	go get github.com/tools/godep

godep-save:
	godep save ./...

all: build test

build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -v -o bin/aprs-dashboard

fmt:
	$(GO) fmt ./...

test:
	$(GO) test -v ./...

bench:
	$(GO) test ./... -bench .

run: build
	$(CURDIR)/aprs-dashboard

docker-build:
	docker info
	docker build -t urlgrey/aprs-dashboard:latest .

docker-deploy:
	docker login -e ${DOCKER_EMAIL} -u ${DOCKER_USER} -p ${DOCKER_PASS}
	docker push urlgrey/aprs-dashboard:latest
