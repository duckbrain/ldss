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
		if len(args) == 0 {
			for _, l := range lib.Languages() {
				fmt.Println(l.Name())
			}
		} else {
			lang := lib.LookupLanguage(args[0])
			fmt.Println(lang.Name())
		}
	},
}

func init() {
	RootCmd.AddCommand(languagesCmd)

	//TODO: Formatting options?
}
