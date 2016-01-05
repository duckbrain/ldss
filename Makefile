DEBUG ?= 1
BINARY = ${GOPATH}/bin/ldss
BINDATA = ${GOPATH}/bin/go-bindata

all: $(BINARY)

ifeq ($(DEBUG), 1)
$(BINARY): ldss/*.go lib/*.go ldss/bindata-debug.go
	go install --tags debug ./ldss
else
$(BINARY): ldss/*.go lib/*.go ldss/bindata.go
	go install ./ldss
endif

run: $(BINARY)
	$(BINARY)
run-lookup: $(BINARY)
	$(BINARY) lookup 1 Ne 3:17

ldss/bindata-debug.go:
	$(BINDATA) -debug -tags "debug" -o "$@" data/...
ldss/bindata.go: $(shell find data -print)
	$(BINDATA) -tags "!debug" -o "$@" data/...

$(BINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

depends: $(BINDATA)
	go get ./...
