DEBUG ?= 1
GOPATH ?= ${HOME}/go
BINARY = ${GOPATH}/bin/ldss
BINDATA = ${GOPATH}/bin/go-bindata
DEPENDS = .depends *.go lib/*.go

all: $(BINARY) assets/bindata_debug.go assets/bindata_release.go

ifeq ($(DEBUG), 1)
$(BINARY): $(DEPENDS) assets/bindata_debug.go
	go install --tags debug .
else
$(BINARY): $(DEPENDS) assets/bindata_release.go
	go install .
endif

assets/bindata_debug.go: $(BINDATA) $(shell find data -print)
	$(BINDATA) -prefix ${PWD} -pkg assets -nomemcopy -debug -tags "debug" -o "$@" data/...
	gofmt -w "$@"
assets/bindata_release.go: $(BINDATA) $(shell find data -print)
	$(BINDATA) -prefix ${PWD} -pkg assets -nomemcopy -tags "!debug" -o "$@" data/...
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


# Android Builds

android:
	cd ldssa/build; ./gradlew build

android-install: android
	adb install -r ldssa/build/build/outputs/apk/build-debug.apk
