package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/duckbrain/ldss/lib"
	_ "github.com/duckbrain/ldss/lib/sources/ldsorg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var langName string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ldss",
	Short: "A tool for viewing the scriptures",
	Long:  `LDS Scriptures is a set of tools for downloading, parsing, and reading the Gospel Library content from The Church of Jesus Christ of Latter-day Saints.`,
	// TODO: Figure out a way, so scripture references can be looked up
	// without specifying the lookup command.
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

func lang() lib.Lang {
	return lib.LookupLanguage(langName)
}

// Init lib for usage
func init() {
	err := lib.Open()
	if err != nil {
		panic(err)
	}
}
