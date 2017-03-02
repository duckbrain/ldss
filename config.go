package main

import (
	"os/user"
	"path"
)

import (
	"github.com/spf13/viper"
)

func init() {
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
}
