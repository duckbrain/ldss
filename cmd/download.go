// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download scripture content",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lang := lang()
		// TODO: Work your own magic here
		fmt.Println("download called")
		if len(args) == 1 {
			lib.DownloadAll(lang, false)
			return
		}
		switch args[1] {
		case "languages", "lang":
			fmt.Println("Downloading language list")
			if err := lib.DownloadLanguages(); err != nil {
				panic(err)
			}
		case "all":
			lib.DownloadAll(lang, true)
		case "missing":
			lib.DownloadAll(lang, false)
		case "cat", "catalog":
			fmt.Println("Downloading \"" + lang.Name + "\" language catalog")
			lib.DownloadCatalog(lang)
		default:
			item, err := lib.AutoDownload(func() (lib.Item, error) {
				refs := lib.Parse(lang, args[1])
				if len(refs) != 1 {
					return nil, fmt.Errorf("Cannot yet handle multiple references")
				}
				return refs[0].Lookup()
			})
			if err != nil {
				panic(err)
			}
			if book, ok := item.(*lib.Book); ok {
				fmt.Printf("Downloading book \"%v\"\n", book.Name())
				lib.DownloadBook(book)
			} else if folder, ok := item.(*lib.Folder); ok {
				fmt.Printf("Downloading folder \"%v\"\n", folder.Name())
				lib.DownloadChildren(folder, false)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
