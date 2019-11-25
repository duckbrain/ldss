package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/sources/churchofjesuschrist"
	"github.com/duckbrain/ldss/lib/storages/filestore"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var library *lib.Library
var refs []lib.Reference
var lang lib.Lang
var ctx context.Context = context.TODO()

var langName string

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
		logger.SetLevel(logrus.TraceLevel)
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
			refs, err = library.Parser.Parse(lang, strings.Join(args, " "))
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
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&langName, "lang", "en", "language for scripture content")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	currentUser, err := user.Current()
	if err == nil {
		viper.SetDefault("DataDirectory", path.Join(currentUser.HomeDir, ".ldss"))
	} else {
		viper.SetDefault("DataDirectory", ".ldss")
	}
	viper.SetDefault("Language", "eng")
	viper.SetDefault("ServerURL", "https://tech.lds.org/glweb")

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/ldss/")
	viper.AddConfigPath("$HOME/.ldss")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
