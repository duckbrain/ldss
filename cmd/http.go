package cmd

import (
	"log"
	"net/http"

	ldsshttp "github.com/duckbrain/ldss/lib/http"
	"github.com/spf13/cobra"
)

var webOpts struct {
	Addr string
}
var webCmd = &cobra.Command{
	Use:     "http",
	Aliases: []string{"web"},
	Short:   "Launch the web server",
	Long:    `Launch the web server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := webOpts
		server := ldsshttp.Server{
			Lang: lang,
			Lib:  library,
		}

		http.Handle("/", server.Handler())
		log.Printf("Listening on: %v\n", opts.Addr)
		return http.ListenAndServe(opts.Addr, nil)
	},
}

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVarP(&webOpts.Addr, "addr", "a", ":1830", "The listen address to use for http")
}
