package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
	"connection"
)

type Config struct {
	OnlineContent    *LDSContent
	OfflineContent   *LocalContent
	Connection       Connection
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

	return &ConfigurationOptions{
		"eng", 
		path.Join(currentUser.HomeDir, ".ldss"), 
		"https://tech.lds.org/glweb"
	}
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

	c.OnlineContent = connection.NewLDSContent(op.ServerURL)
	c.OfflineContent = connection.NewLocalContent(op.DataDirectory)
	
	cache := NewCacheConnection()
	cache.Open(c.OfflineContent.GetCachePath())

	c.Languages = cache

	c.Download = new(Downloader)
	c.Download.online = c.OnlineContent
	c.Download.offline = c.OfflineContent

	c.SelectedLanguage = c.Languages.GetByUnknown(op.Language)

	return c
}
