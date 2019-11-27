package cmd

import (
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the storage",
	Run: func(cmd *cobra.Command, args []string) {
		err := library.Store.Clear(ctx)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(clearCmd)
}
