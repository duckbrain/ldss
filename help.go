package main

import "fmt"

func printAsset(path string) {
	data, err := Asset(path)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))
}

func PrintInstructions() {
	printAsset("data/help/help")
}

func PrintCommandInstructions(command string) {
	data, err := Asset("data/help/" + command)
	if err != nil {
		fmt.Printf("Command \"%s\" is not recognized or no extra help is available.\n", command)
		return
	}
	fmt.Print(string(data))
}
