# LDS Scriptures

LDS Scriptures port to native command line in Go.

LDS Scriptures is a set of tools for downloading, parsing, and reading the Gospel Library content from [The Church of Jesus Christ of Latter-day Saints](http://lds.org). 

The following interfaces are planned to be implemented as part of this project or as a separate project. Currently, only the web-based interface is intended for use.

- [x] Web server
- [ ] Android App
- [ ] Command line interface
- [ ] Native OS graphical interface using [github.com/andlabs/ui](https://github.com/andlabs/ui)
- [ ] Curses interface

## Installation instructions

There are currently no binary releases for the application, but you can compile it with go.

1. [Install Go](https://golang.org/doc/install)
2. Run `go get github.com/duckbrain/ldss`. This should download and compile ldss and all its dependencies.
3. Run `ldss web` to start the web server. It will default to port 1830. If you would like to use a different port, [that's currently broken](https://github.com/duckbrain/ldss/issues/1).

## For Developers

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
