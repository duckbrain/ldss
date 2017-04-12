package cmd

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download scripture content",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lang := lang()
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
}
