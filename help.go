package main

import "fmt"

func PrintInstructions() {
	fmt.Println("usage: ldss <command> [<arguments>]")
	fmt.Println()
	fmt.Println("The following are available as commands:")
	fmt.Println("   help [command name]           Print these instructions or get more detailed information on a command")
	fmt.Println("   download lang|<lang>|<book>   Download the catalog or book specified so it can be read or referenced")
	fmt.Println("   read <ref>                    Opens up a scripture reference to read; See \"ldss help ref\" to see how to write references.")
	fmt.Println("   index <ref>                   List the child nodes of the speficied reference")
	fmt.Println("   print <ref>                   Print out the contents of the scripture reference given without the default reader or HTML parsing")
	fmt.Println("   lookup <ref>                  Perform a lookup on the scripture reference to determine if it is valid")
	fmt.Println("   <ref>                         Shorthand for \"ldss read <ref>\"")
	fmt.Println()
	fmt.Println("Other topics to get help on:")
	fmt.Println("   config                        Information on the configuration file and its settings")
	fmt.Println("   ref                           Information on scripture references and valid formats to provide them in")
	fmt.Println("   lang                          Information on providing languages")
}

func PrintCommandInstructions(command string) {
	switch (command) {
	case "config":
		PrintConfigInstructions()
	default:
		fmt.Printf("Command \"%s\" is not recognized.\n", command);
	}
}

func PrintConfigInstructions() {
}
