GO ?= go	
export GOPATH = $(CURDIR)/_vendor

all: build test
build:
	$(GO) build -o aprs-dashboard

fmt:
	$(GO) fmt

test:
	$(GO) test

bench:
	$(GO) test -bench . -parallel 2

benchmem:
	$(GO) test -bench . -benchmem -parallel 2

run: build
	$(CURDIR)/aprs-dashboard
