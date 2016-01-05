DEBUG ?= 1

all: ldss

ifeq ($(DEBUG), 1)
ldss: *.go lib/*.go bindata-debug.go
	go install --tags debug
else
ldss: *.go lib/*.go bindata.go
	go install 
endif

run: ldss
	./ldss

bindata-debug.go:
	${GOPATH}/bin/go-bindata -debug -tags "debug" -o "$@" data/...
bindata.go: $(shell find data -print)
	${GOPATH}/bin/go-bindata -tags "!debug" -o "$@" data/...

run-lookup: ldss
	./ldss lookup 1 Ne 3:17

depends:
	go get -u github.com/jteeuwen/go-bindata/...
	go get ./...
