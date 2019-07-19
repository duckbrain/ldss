/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/bmaupin/go-epub"
	"github.com/duckbrain/ldss/lib"
	"github.com/spf13/cobra"
)

var epubFilename string

// epubCmd represents the epub command
var epubCmd = &cobra.Command{
	Use:   "epub",
	Short: "Generate an epub file from a section of scripture",
	Long:  `Generates an epub file in the speficied directory using the provided scripture `,
	Run: func(cmd *cobra.Command, args []string) {
		lang := lang()
		lookupString := strings.Join(args, " ")
		refs := lib.Parse(lang, lookupString)
		if len(refs) != 1 {
			panic(fmt.Errorf("Multiple references not allowed for generating epubs"))
		}
		item, err := refs[0].Lookup()
		if err != nil {
			panic(err)
		}

		if len(epubFilename) == 0 {
			// TODO Convert filename to a slug
			epubFilename = fmt.Sprintf("%v.epub", item.Name())
		}

		e := epub.NewEpub(item.Name())
		e.SetAuthor("The Church of Jesus Christ of Latter-day Saints")
		e.SetLang(lang.Code())
		e.SetIdentifier(fmt.Sprintf("ldss:%v", item.Path()))
		// TODO e.SetCover to a cover image if the node has one

		err = e.Write(epubFilename)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(epubCmd)

	epubCmd.Flags().StringVarP(&epubFilename, "filename", "f", "", "Filename of the epub file to create. If none is specified, one will be generated in the current directory.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// epubCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// epubCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
