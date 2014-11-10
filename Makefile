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
	$(GO) test -bench .

run: build
	$(CURDIR)/aprs-dashboard
