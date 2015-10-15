all: ldss

ldss: *.go
	go build

run: ldss
	./ldss

run-lookup: ldss
	./ldss lookup 1 Ne 3:17
