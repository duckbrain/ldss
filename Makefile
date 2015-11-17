all: ldss

ldss: *.go ldslib/*.go bindata.go
	go build 

run: ldss
	./ldss

bindata.go: $(shell find data -print)
	${GOPATH}/bin/go-bindata -debug data/...

run-lookup: ldss
	./ldss lookup 1 Ne 3:17
