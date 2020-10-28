.PHONY: all docker docker-clean publish-docker

REPO=github.com/xtao/xtor
#VERSION?=$(shell git describe HEAD | sed s/^v//)
VERSION?=1.0
DATE?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-X ${REPO}/version.Version=${VERSION} -X ${REPO}/version.BuildDate=${DATE}
TOOLS=${GOPATH}/bin/go-bindata \
      ${GOPATH}/bin/go-bindata-assetfs \
      ${GOPATH}/bin/goconvey
SRC=$(shell find . -name '*.go')
TESTFLAGS="-v"

PACKAGES=$(shell go list ./... | grep -v /vendor/)


vet:
	go vet ${PACKAGES}

lint:
	go list ./... | grep -v /vendor/ | grep -v assets | xargs -L1 golint -set_exit_status


rpm/xtorsvr/xtorsvr: ${SRC}
	CGO_ENABLED=1 GOOS=linux go build -ldflags "${LDFLAGS}" -a -installsuffix cgo -o $@ github.com/xtao/xtor/cmd/xtorsvr

rpm/xtorcli/xtorcli: ${SRC}
	CGO_ENABLED=1 GOOS=linux go build -ldflags "${LDFLAGS}" -a -installsuffix cgo -o $@ github.com/xtao/xtor/cmd/xtorcli

rpm: rpm/xtorsvr/xtorsvr rpm/xtorcli/xtorcli
	cd rpm && ./autogen.sh
	cd rpm && ./configure
	cd rpm && make clean && make rpm

clean:
	rm -rf cmd/xtorsvr/xtorsvr
	rm -rf cmd/xtorcli/xtorcli
	rm -rf rpm/xtorsvr/xtorsvr
	rm -rf rpm/xtorcli/xtorcli
	rm -rf rpm/rpmbuild
	rm -rf rpm/config.*
	rm -rf rpm/configure
	rm -rf rpm/install.sh
	rm -rf rpm/Makefile
	rm -rf rpm/Makefile.in
	rm -rf rpm/*.tar.gz
	rm -rf rpm/*.rpm
