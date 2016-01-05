DEBUG ?= 1
BINARY = ${GOPATH}/bin/ldss

all: $(BINARY)

ifeq ($(DEBUG), 1)
$(BINARY): ldss/*.go lib/*.go ldss/bindata-debug.go
	go install --tags debug ./ldss
else
$(BINARY): ldss/*.go lib/*.go ldss/bindata.go
	go install ./ldss
endif

run: ldss
	$(BINARY)

ldss/bindata-debug.go:
	${GOPATH}/bin/go-bindata -debug -tags "debug" -o "$@" data/...
ldss/bindata.go: $(shell find data -print)
	${GOPATH}/bin/go-bindata -tags "!debug" -o "$@" data/...

run-lookup: ldss
	$(BINARY) lookup 1 Ne 3:17

depends:
	go get -u github.com/jteeuwen/go-bindata/...
	go get ./...
