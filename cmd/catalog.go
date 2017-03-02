package cmd

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
	"github.com/spf13/cobra"
)

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Prints the root level catalog",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			langs, err := lib.Languages()
			if err != nil {
				panic(err)
			}
			for _, l := range langs {
				fmt.Println(l.String())
			}
		} else {
			catalog, err := lang().Catalog()
			if err != nil {
				panic(err)
			}
			fmt.Println(catalog.String())
		}
	},
}

func init() {
	RootCmd.AddCommand(catalogCmd)
}
