package main

import (
	"fmt"

	"github.com/duckbrain/ldss/assets"
)

func printAsset(path string) {
	data, err := assets.Asset(path)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))
}

func PrintInstructions() {
	printAsset("data/help/help")
}

func PrintCommandInstructions(command string) {
	data, err := assets.Asset("data/help/" + command)
	if err != nil {
		fmt.Printf("Command \"%s\" is not recognized or no extra help is available.\n", command)
		return
	}
	fmt.Print(string(data))
}
