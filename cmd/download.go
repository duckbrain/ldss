package cmd

import (
	"fmt"
	"strings"

	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/dl"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download scripture content",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lang := lang()

		//TODO: Make these flags
		force := false
		recursive := false

		query := strings.Join(args[1:], " ")
		for _, ref := range lib.Parse(lang, query) {

			item, err := ref.Lookup()
			if err != nil {
				panic(err)
			}

			if x, ok := item.(dl.Downloader); ok {
				if force || !x.Downloaded() {
					dl.Enqueue(x, nil)
				}
			}

			if recursive {
				//TODO
				_ = item
			}

		}

		//TODO: Register download watcher and output status
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)

	dlUpdates := make(chan dl.DownloadStatus)
	go func() {
		var lastStarted dl.DownloadStatus
		for s := range dlUpdates {
			switch s.Stage {
			case dl.Started:
				fmt.Printf("Starting download for %v.\n", s.Name())
				lastStarted = s
			case dl.Complete:
				if lastStarted.Downloader != nil && lastStarted.Hash() == s.Hash() {
					fmt.Println("Done")
				} else {
					fmt.Printf("Download %v finished!\n", s.Name())
				}
			}
		}
	}()
	dl.Listen(dlUpdates)
}
