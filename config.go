package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type Config struct {
	OnlineContent    *LDSContent
	OfflineContent   *LocalContent
	Languages        LanguageLoader
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
	file, err := os.Open(path.Join(op.DataDirectory, "config.json"))
	if err != nil {
		// File does not exits, continue
		return
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(op)
	if err != nil {
		panic(err)
	}
}

func LoadConfiguration(op *ConfigurationOptions) Config {
	c := Config{}

	c.OnlineContent = NewLDSContent(op.ServerURL, 17)
	c.OfflineContent = NewLocalContent(op.DataDirectory)
	
	cache := NewCacheConnection()
	cache.Open(c.OfflineContent.GetCachePath())

	c.Languages = cache

	c.Download = new(Downloader)
	c.Download.online = c.OnlineContent
	c.Download.offline = c.OfflineContent

	c.SelectedLanguage = c.Languages.GetByUnknown(op.Language)

	return c
}
