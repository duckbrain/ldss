package cmd

import (
	"fmt"
	"strings"

	"github.com/duckbrain/ldss/lib"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type cmdcolors struct {
	title, subtitle, summary, verse, content, message *color.Color
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
	Run: func(cmd *cobra.Command, args []string) {
		colors := colors(true)
		lang := lang()
		lookupString := strings.Join(args, " ")
		refs := lib.Parse(lang, lookupString)
		if len(refs) != 1 {
			panic(fmt.Errorf("Multiple references not implemented"))
		}

		item, err := refs[0].Lookup()
		if err != nil {
			panic(err)
		}

		if node, ok := item.(lib.Contenter); ok {
			if content, err := node.Content(); err == nil {
				z := content.Parse()
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
						colors.verse.Print(z.ParagraphVerse())
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
				fmt.Println(err)
			}
		} else {
			fmt.Println(item.Name())
		}

		for _, child := range item.Children() {
			fmt.Printf("- %v {%v}\n", child.Name(), child.Path())
		}
	},
}

func init() {
	RootCmd.AddCommand(lookupCmd)
}
