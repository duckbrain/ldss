all: ldss

ldss: *.go
	go build --tags "libsqlite3 linux"

run: ldss
	./ldss

run-lookup: ldss
	./ldss lookup 1 Ne 3:17
