// +build nobuild

package cmd

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
	"github.com/spf13/cobra"
)

var languagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "Lists the languages available",
	Long:  `Lists the available languages. If passed an extra parameter, it will return the language that matches. The parameter can be the name of the language in English or the native language, the 2 or 3 letter codes of the language, or the numeric ID of the language.`,
	Run: func(cmd *cobra.Command, args []string) {
		var langs []lib.Lang
		if len(args) == 0 {
			langs = lib.Languages()
		} else {
			q := args[0]
			lang := lib.LookupLanguage(q)
			if lang == nil {
				panic(fmt.Errorf("Language %v not found", q))
			}
			langs = []lib.Lang{lang}
		}
		for _, l := range langs {
			fmt.Printf("%5v %v\n", l.Code(), l.Name())
		}
	},
}

func init() {
	RootCmd.AddCommand(languagesCmd)

	//TODO: Formatting options?
}
