package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/sources/churchofjesuschrist"
	"github.com/duckbrain/ldss/lib/storages/filestore"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var library *lib.Library
var refs []lib.Reference
var lang lib.Lang
var ctx context.Context = context.TODO()

var langName string
var logLevel string
var logLevels map[string]logrus.Level

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ldss",
	Short: "A tool for viewing the scriptures",
	Long:  `LDS Scriptures is a set of tools for downloading, parsing, and reading the Gospel Library content from The Church of Jesus Christ of Latter-day Saints.`,
	// TODO: Figure out a way, so scripture references can be looked up
	// without specifying the lookup command.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := logrus.New()
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		logger.SetLevel(level)
		store, err := filestore.New(".ldss")
		if err != nil {
			logger.Error(errors.Wrap(err, "store init"))
			return err
		}
		lang = lib.Lang(langName)
		library = lib.Default
		library.Store = store
		library.Index = store
		library.Logger = logger
		library.Register(churchofjesuschrist.Default)

		if len(args) > 0 {
			refs, err = library.Parse(ctx, lang, strings.Join(args, " "))
			library.Logger.Debugf("parsing refs for lang: %v, args: %v, refs: %v", lang, args, refs)
			if err != nil {
				library.Logger.Fatalf("parse reference: %v", err)
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&langName, "lang", "en", "language for scripture content")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "warning", "Logging level: info, debug, warning, error")
}
