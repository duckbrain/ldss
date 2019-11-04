package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download scripture content",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		for _, ref := range refs {
			err := library.Download(context.Background(), ref.Index)
			if err != nil {
				panic(err)
			}
		}
	},
}
