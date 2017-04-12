package cmd

import (
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:    "generate",
	Short:  "A collection of generators",
	Long:   ``,
	Hidden: true,
}

func init() {
	RootCmd.AddCommand(generateCmd)

}
