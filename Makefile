DEBUG ?= 1
BINARY = ${GOPATH}/bin/ldss
BINDATA = ${GOPATH}/bin/go-bindata
DEPENDS = *.go lib/*.go .depends $(BINDATA) 

all: $(BINARY) bindata_debug.go bindata_release.go

ifeq ($(DEBUG), 1)
$(BINARY): $(DEPENDS) bindata_debug.go
	go install .
else
$(BINARY): $(DEPENDS) bindata_release.go
	go install --tags release .
endif

run: $(BINARY)
	$(BINARY)
run-lookup: $(BINARY)
	$(BINARY) lookup 1 Ne 3:17

bindata_debug.go:
	$(BINDATA) -nomemcopy -debug -tags "!release" -o "$@" data/...
	gofmt -w "$@"
bindata_release.go: $(shell find data -print)
	$(BINDATA) -nomemcopy -tags "release" -o "$@" data/...
	gofmt -w "$@"

$(BINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

.depends: 
	go get -d ./...
	@echo "Flags make that dependances are gotten" > .depends
	
format:
	go fmt *
	go fmt lib/*

clean:
	rm -f $(BINARY)
	rm -f bindata.go bindata_debug.go bindata_release.go
	go clean -r

clean-tree: clean
	@read -r -p "This will delete all GitHub and GoLang.org source. Cancel if this is not what you want." response
	rm -f .depends
	rm -f ${GOPATH}/bin/*
	rm -rf ${GOPATH}/pkg/*
	rm -rf ${GOPATH}/src/github.com
	rm -rf ${GOPATH}/src/golang.org
