// +build !noshell

package main

import (
	"bufio"
	"fmt"
	_ "io"
	"os"
)

type shell struct {
	appinfo
}

func init() {
	apps["shell"] = &shell{}
}

func (app shell) run() {
	fmt.Printf("Welcome to the LDS Scriptures interactive shell.\n")
	cin := bufio.NewReader(os.Stdin)

	for {
		app.handleLine(cin)
	}
}

func (app shell) handleLine(cin *bufio.Reader) {
	fmt.Printf("> ")
	line, isPrefix, err := cin.ReadLine()
	_ = line
	_ = isPrefix
	if err != nil {
		panic(err)
	}
}
