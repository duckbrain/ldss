DEBUG ?= 1
BINARY = ${GOPATH}/bin/ldss
BINDATA = ${GOPATH}/bin/go-bindata
DEPENDS = ldss/*.go lib/*.go .depends $(BINDATA) 

all: $(BINARY)

ifeq ($(DEBUG), 1)
$(BINARY): $(DEPENDS) ldss/bindata-debug.go
	go install --tags debug ./ldss
else
$(BINARY): $(DEPENDS) ldss/bindata.go
	go install ./ldss
endif

run: $(BINARY)
	$(BINARY)
run-lookup: $(BINARY)
	$(BINARY) lookup 1 Ne 3:17

ldss/bindata-debug.go:
	$(BINDATA) -nomemcopy -debug -tags "debug" -o "$@" data/...
ldss/bindata.go: $(shell find data -print)
	$(BINDATA) -nomemcopy -tags "!debug" -o "$@" data/...

$(BINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

.depends: 
	go get -d ./...
	@echo "Flags make that dependances are gotten" > .depends

clean:
	rm -f ${GOPATH}/bin/ldss
	go clean -r

clean-tree: clean
	@read -r -p "This will delete all GitHub and GoLang.org source. Cancel if this is not what you want." response
	rm -f .depends
	rm -f ${GOPATH}/bin/*
	rm -rf ${GOPATH}/pkg/*
	rm -rf ${GOPATH}/src/github.com
	rm -rf ${GOPATH}/src/golang.org
