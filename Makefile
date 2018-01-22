CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/aaronland/picturebook
	cp -r flags src/github.com/aaronland/picturebook/
	cp -r functions src/github.com/aaronland/picturebook/
	cp -r *.go src/github.com/aaronland/picturebook/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-image-tools"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/microcosm-cc/exifutil"
	@GOPATH=$(GOPATH) go get -u "github.com/rwcarlsen/goexif/exif"
	@GOPATH=$(GOPATH) go get -u "github.com/jung-kurt/gofpdf"
	@GOPATH=$(GOPATH) go get -u "github.com/rainycape/unidecode"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/picturebook cmd/picturebook.go
