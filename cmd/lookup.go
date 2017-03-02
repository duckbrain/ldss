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

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		colors := colors(true)
		lang := lang()
		lookupString := strings.Join(args[1:], " ")
		refs := lib.Parse(lang, lookupString)
		if len(refs) != 1 {
			panic(fmt.Errorf("Multiple references not implemented"))
		}
		item, err := refs[0].Lookup()
		if err != nil {
			panic(err)
		}

		if node, ok := item.(*lib.Node); ok {
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
				return
			}
		}

		children, err := item.Children()
		if err != nil {
			panic(err)
		}
		fmt.Println(item)
		for _, child := range children {
			fmt.Printf("- %v\n", child)
		}
	},
}

func init() {
	RootCmd.AddCommand(lookupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lookupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lookupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
