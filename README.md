# LDS Scriptures

![](data/web/static/favicon.ico)

ldss is a set of tools for downloading, parsing, and reading the Gospel Library content from [The Church of Jesus Christ of Latter-day Saints](http://lds.org). 

My goal would be to eventually have many ways to access the richest scripture content in as many places as possible, but for now, it mainly provides a web server.

The content available through this software is owned by the church, and is distributed by the church. This is simply an alternate client to view the content in.

### Web Interface (your best bet)

The web interface is what I use for my daily usage. It works great in institute classes and for preparing lessons. Most features will probably work best and show up here first.

### Native Graphical User Interface (easiest)

I tried to do th is in Go with [github.com/andlabs/ui](https://github.com/andlabs/ui), but I was basically re-implementing the web browser. I still think it's a really good project though. [github.com/zserge/webview](https://github.com/zserge/webview) make it easy to port the web interface to a desktop application, so it's almost as good as the web interface. It'll need some extra functionality to make up for not having browser navigation (back, forward, tabs), but it's the way the project will go for now.

### Android App (not well maintained)

I did manage to generate an Android app with [gowebview](https://github.com/microo8/gowebview), but it really needs some work before I would call it good enough for daily usage. If you are willing and able to help, pull requests would be much appreciated. I'm also perfectly happy with you forking the project, but let me know so I can link back.

### Command Line Interface (a little broken)

There was a half decent command line interface, but since I moved to Viper, I haven't fully redone it yet. I'll get around to it, but the web interface is so much easier to use (and to work on for me).

## Installation instructions

There are currently no binary releases for the application. If you know how to cross-compile [go-sqlite3](https://github.com/mattn/go-sqlite3) on Linux, let me know. I'm a little lacking on the build machines. You can install the package by doing the following.

1. [Install Go](https://golang.org/doc/install)
2. Run `go get github.com/duckbrain/ldss`. This should download and compile ldss and all its dependencies.
3. Run `ldss web` to start the web server. It will default to port 1830. If you would like to use a different port, specify it with the `--port` parameter.

## For Developers

To quickly download and build the project, install Go, and **set your `GOPATH`**. 

```bash
go get github.com/duckbrain/ldss
```

This will check the project out to `$GOPATH/src/github.com/duckbrain/ldss` and compile it to `$GOPATH/bin/ldss`.

This project uses [go-bindata](http://github.com/jteeuwen/go-bindata) to generate `assets/bindata_release.go`. If anything is changed in the `data\` directory, you must regenerate the file. The easiest way to do this is with `make` (if you have it installed). Otherwise, you need to do the following.

- `go get -u github.com/jteeuwen/go-bindata/...` to download and compile go-bindata.
- `$GOPATH/bin/go-bindata -pkg assets -nomemcopy -tags "!debug" -o assets/bindata_release.go data/...` to generate the new release build. Make sure all pull-requests use these parameters, but you can modify the generation for your own builds.

There is also a debug mode that will generate a binary that reads the files from the project instead of building them into the executable. You can read the `Makefile` for details on that.

If you are new to Go, I would recommend reading [How to Write Go Code](https://golang.org/doc/code.html) in the Go documentation for various other commands you can use in development.
