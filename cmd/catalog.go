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
			cItem, err := lib.AutoDownload(func() (lib.Item, error) {
				return lang().Catalog()
			})
			if err != nil {
				panic(err)
			}
			catalog := cItem.(*lib.Catalog)
			fmt.Println(catalog.String())
		} else {
			langs, err := lib.Languages()
			if err != nil {
				panic(err)
			}
			for _, l := range langs {
				fmt.Println(l.String())
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(catalogCmd)
}
