DEBUG ?= 1

all: ldss

ifeq (DEBUG, 1)
ldss: *.go ldslib/*.go bindata-debug.go
	go build 
else
ldss: *.go ldslib/*.go bindata.go
	go build 
endif

run: ldss
	./ldss

bindata-debug.go:
	rm -f bindata.go
	${GOPATH}/bin/go-bindata -debug data/...
bindata.go: $(shell find data -print)
	rm -f bindata-debug.go
	${GOPATH}/bin/go-bindata data/...

run-lookup: ldss
	./ldss lookup 1 Ne 3:17

depends:
	go get -u github.com/jteeuwen/go-bindata/...
	go get ./...
