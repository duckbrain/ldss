// +build exclude

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type shell struct {
	appinfo
	cmd *cmd
}

func init() {
	addApp("shell", &shell{})
}

func (app shell) register(*Configuration) {}

func (app *shell) run() {
	fmt.Printf("Welcome to the LDS Scriptures interactive shell.\n")
	cin := bufio.NewReader(os.Stdin)
	app.cmd = new(cmd)
	app.cmd.appinfo = app.appinfo

	for {
		app.handleLine(cin)
	}
}

func (app shell) handleLine(cin *bufio.Reader) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	fmt.Printf("> ")
	line, isPrefix, err := cin.ReadLine()
	if err != nil {
		panic(err)
	}
	if isPrefix {
		panic(fmt.Errorf("Line too long"))
	}
	args := strings.Fields(string(line))
	switch args[0] {
	case "exit":
		os.Exit(0)
	default:
		app.cmd.args = args
		app.cmd.run()
	}
}
