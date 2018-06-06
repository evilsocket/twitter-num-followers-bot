SHELL := bash

all: tnfbot

godep:
	@go get -u github.com/golang/dep/...

deps: godep
	@dep ensure

tnfbot: deps
	@go build -o tnfbot .

clean:
	@rm -rf tnfbot

install:
	@cp tnfbot /usr/local/bin/
