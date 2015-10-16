package main

import (
	"os"
	"os/user"
	"path"
)

type Config struct {
	OnlineContent    *LDSContent
	OfflineContent   *LocalContent
	Languages        *LanguageLoader
	Download         *Downloader
	SelectedLanguage *Language
}

type ConfigurationOptions struct {
	Language      string
	DataDirectory string
	ServerURL     string
}

func loadDefaultOptions() *ConfigurationOptions {
	currentUser, err := user.Current()

	if err != nil {
		panic(err)
	}

	op := new(ConfigurationOptions)
	op.Language = "eng"
	op.DataDirectory = path.Join(currentUser.HomeDir, ".ldss")
	op.ServerURL = "https://tech.lds.org/glweb"

	return op
}

func loadParameterOptions(op *ConfigurationOptions) []string {
	args := os.Args[1:]
	return args
}

func loadFileOptions(op *ConfigurationOptions) {

}

func LoadConfiguration(op *ConfigurationOptions) Config {
	c := Config{}

	c.OnlineContent = new(LDSContent)
	c.OnlineContent.BasePath = op.ServerURL
	c.OfflineContent = new(LocalContent)
	c.OfflineContent.BasePath = op.DataDirectory

	c.Languages = new(LanguageLoader)
	c.Languages.content = c.OfflineContent

	c.Download = new(Downloader)
	c.Download.online = c.OnlineContent
	c.Download.offline = c.OfflineContent

	c.SelectedLanguage = c.Languages.GetByUnknown(op.Language)

	return c
}
