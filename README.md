# LDS Scriptures

LDS Scriptures port to native command line in Go.

LDS Scriptures is a set of tools for downloading, parsing, and reading the Gospel Library content from [The Church of Jesus Christ of Latter-day Saints](http://lds.org). It contains interfaces the following interfaces, either in part, in plan, or completely finished. 

- [ ] Command line
- [ ] Web server
- [ ] GTK GUI
- [ ] Interactive shell
- [ ] Curses interface

## Build Instructions

LDS Scriptures can compile for debug or release. The major difference between the two is that the resources in the data/ directory will be statically linked  in release mode. There may also be more error and logging output in debug mode as well.

### Makefile build

1. Check out the repository into $GOPATH/src/ldss and cd into the directory
2. Run `make DEBUG=1`

The `DEBUG=1` can be ommitted because that is the default value. If you set `DEBUG=0`, then it will compile in release mode.

### Manual Build

1. Check out the repository into `$GOPATH/src/ldss` and cd into the directory
2. Run `go get -u github.com/jteeuwen/go-bindata/...` to download go-bindata
3. Run `$GOPATH/bin/go-bindata data/... -o ldss/bindata.go` to generate bindata.go. You can pass `--debug` to generate a debug version as well.
4. Run `go get -d ./...` to download all dependencies
5. Run `go install ./ldss` to compile and you will find the binary in `$GOPATH/bin/ldss`.
