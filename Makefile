GO ?= godep go
COVERAGEDIR = ./coverage
all: build test cover

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
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	$(GO) test -v ./handlers -race -cover -coverprofile=$(COVERAGEDIR)/handlers.coverprofile
	$(GO) test -v ./parser -race -cover -coverprofile=$(COVERAGEDIR)/parser.coverprofile

cover:
	$(GO) tool cover -html=$(COVERAGEDIR)/handlers.coverprofile -o $(COVERAGEDIR)/handlers.html
	$(GO) tool cover -html=$(COVERAGEDIR)/parser.coverprofile -o $(COVERAGEDIR)/parser.html

bench:
	$(GO) test ./... -cpu 2 -bench .

run: build
	$(CURDIR)/aprs-dashboard

docker-build:
	docker info
	docker build -t urlgrey/aprs-dashboard:latest .

docker-deploy:
	docker login -e ${DOCKER_EMAIL} -u ${DOCKER_USER} -p ${DOCKER_PASS}
	docker push urlgrey/aprs-dashboard:latest
