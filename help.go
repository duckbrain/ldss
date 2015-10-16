package main

import "fmt"

func PrintInstructions() {
	fmt.Println("usage: ldss <command> [<arguments>]")
	fmt.Println()
	fmt.Println("The following are available as commands:")
	fmt.Println("   help [command name]           Print these instructions or get more detailed information on a command")
	fmt.Println("   download lang|<lang>|<book>   Download the catalog or book for offline use")
	fmt.Println("   languages                     List all available languages")
	fmt.Println("   read <ref>                    Open up a scripture reference to read")
	fmt.Println("   index <ref>                   List the child nodes of the speficied reference")
	fmt.Println("   print <ref>                   Print out the contents of the scripture reference given without the default reader or HTML parsing")
	fmt.Println("   lookup <ref>                  Perform a lookup on the scripture reference to determine if it is valid")
	fmt.Println("   <ref>                         Shorthand for \"ldss read <ref>\"")
	fmt.Println()
	fmt.Println("What to fill in the <blanks>:")
	fmt.Println("   <lang>                        Information on providing languages")
	fmt.Println("   <ref>                         catalog|<folder>|<book>|<node>")
	fmt.Println("   <folder>                      ID number of a folder or special name; eg: \"scriptures\"")
	fmt.Println()
	fmt.Println("Other topics to get help on:")
	fmt.Println("   config                        Information on the configuration file and its settings")
}

func PrintCommandInstructions(command string) {
	switch command {
	case "help":
		PrintInstructions()
	case "config":
		PrintConfigInstructions()
	default:
		fmt.Printf("Command \"%s\" is not recognized or no extra help is available.\n", command)
	}
}

func PrintConfigInstructions() {
	fmt.Println("Your config file is generally in ~/.ldss/config.json")
	fmt.Println()
	fmt.Println("The following are properties you can set:")
}
