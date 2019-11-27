package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/duckbrain/ldss/lib"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type cmdcolors struct {
	title, subtitle, summary, verse, content, message *color.Color
}

var LookupOpts struct {
	ForceDownload bool
	Format        string
}

func colors(enabled bool) *cmdcolors {
	c := cmdcolors{}
	c.content = color.New()
	c.title = color.New(color.Bold).Add(color.Underline).Add(color.FgWhite).Add(color.BgHiMagenta)
	c.subtitle = color.New(color.Bold).Add(color.FgGreen)
	c.summary = color.New(color.Italic).Add(color.BgBlue).Add(color.FgBlack)
	c.verse = color.New(color.Bold).Add(color.FgRed)
	c.message = color.New(color.FgHiYellow).Add(color.Italic)
	color.NoColor = !enabled
	return &c
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Prints a scripture reference to the stdout",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		colors := colors(true)
		if len(refs) != 1 {
			panic("multiple references not implemented")
		}

		if LookupOpts.ForceDownload {
			err := library.Download(ctx, refs[0].Index)
			if err != nil {
				return errors.Wrap(err, "Download")
			}
		}
		item, err := library.LookupAndDownload(ctx, refs[0].Index)
		if err != nil {
			return errors.Wrap(err, "LookupAndDownload")
		}

		switch LookupOpts.Format {
		case "default":
		case "json":
			data, err := json.Marshal(item)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		default:
			return errors.New("uknown format")
		}

		if z := item.Content.Parse(); z != nil {
			for z.NextParagraph() {
				color := colors.content
				switch z.ParagraphStyle() {
				case lib.ParagraphStyleTitle:
					color = colors.title
				case lib.ParagraphStyleSummary:
					color = colors.summary
				case lib.ParagraphStyleChapter:
					color = colors.subtitle
				}
				if z.ParagraphVerse() > 0 {
					colors.verse.Printf("%v ", z.ParagraphVerse())
				}
				for z.NextText() {
					if z.TextStyle() == lib.TextStyleFootnote {
						continue
					}
					color.Print(z.Text())
				}
				color.Println("")
			}
		} else {
			fmt.Println(item.Name)
		}

		for _, child := range item.Children {
			childItem, err := library.Lookup(ctx, child)
			if err != nil {
				return err
			}
			fmt.Printf("- %v {%v}\n", childItem.Name, childItem.Path)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(lookupCmd)
	lookupCmd.Flags().BoolVarP(&LookupOpts.ForceDownload, "force-download", "d", false, "Force the download, even if it's already downloaded")
	lookupCmd.Flags().StringVarP(&LookupOpts.Format, "format", "f", "default", "Format to output in: default, json")
}
