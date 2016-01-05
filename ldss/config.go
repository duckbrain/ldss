package main

import (
	"encoding/json"
	"fmt"
	"ldss/lib"
	"os"
	"os/user"
	"path"
	"strconv"
)

type Config struct {
	op             ConfigurationOptions
	OnlineContent  lib.Source
	OfflineContent lib.Source
	Library        *lib.Library
	Reference      lib.RefParser
	Download       *lib.Downloader
}

type ConfigurationOptions struct {
	Language        string
	DataDirectory   string
	ServerURL       string
	WebPort         int
	WebTemplatePath string
}

func (op *ConfigurationOptions) String() string {
	return fmt.Sprintf("Language:      %v\n", op.Language) +
		fmt.Sprintf("ServerURL:       %v\n", op.ServerURL) +
		fmt.Sprintf("DataDirectory:   %v\n", op.DataDirectory) +
		fmt.Sprintf("WebPort:         %v\n", op.WebPort) +
		fmt.Sprintf("WebTemplatePath: %v\n", op.WebTemplatePath)
}

func loadDefaultOptions() ConfigurationOptions {
	currentUser, err := user.Current()

	if err != nil {
		panic(err)
	}

	op := ConfigurationOptions{}
	op.Language = "eng"
	op.DataDirectory = path.Join(currentUser.HomeDir, ".ldss")
	op.ServerURL = "https://tech.lds.org/glweb"
	op.WebPort = 1830

	return op
}

func loadParameterOptions(op ConfigurationOptions) (ConfigurationOptions, []string) {
	args := os.Args[1:]
	var err error
	for i := 0; i < len(args); {
		switch args[i] {
		case "-p":
			op.WebPort, err = strconv.Atoi(args[i+1])
			if err != nil {
				panic(fmt.Errorf("Could not convert port \"%v\" to an integer", args[i+1]))
			}
			args = args[:i+copy(args[i:], args[i+2:])]
		case "-l":
			op.Language = args[i+1]
			args = args[:i+copy(args[i:], args[i+2:])]
		default:
			i++
		}
	}
	return op, args
}

func loadFileOptions(op ConfigurationOptions) ConfigurationOptions {
	file, err := os.Open(path.Join(op.DataDirectory, "config.json"))
	if err != nil {
		// File does not exits, continue
		return op
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(op)
	if err != nil {
		panic(err)
	}
	return op
}

func LoadConfiguration(op ConfigurationOptions) Config {
	c := Config{op: op}
	c.OnlineContent = lib.NewOnlineSource(op.ServerURL)
	c.OfflineContent = lib.NewOfflineSource(op.DataDirectory)
	c.Library = lib.NewLibrary(c.OfflineContent)
	c.Download = lib.NewDownloader(c.OnlineContent, c.OfflineContent)
	return c
}

func (c *Config) Languages() []lib.Language {
	langs, err := c.Library.Languages()
	if err != nil {
		c.Download.Languages()
		langs, err = c.Library.Languages()
		if err != nil {
			panic(err)
		}
	}
	return langs
}

func (c *Config) Language(s string) *lib.Language {
	lang, err := c.Library.Language(s)
	if err != nil {
		//TODO: Output stderr
		c.Download.Languages()
		lang, err = c.Library.Language(s)
		if err != nil {
			panic(err)
		}
	}
	return lang
}

func (c *Config) SelectedLanguage() *lib.Language {
	return c.Language(c.op.Language)
}

func (c *Config) Catalog(lang *lib.Language) *lib.Catalog {
	catalog, err := c.Library.Catalog(lang)
	if err != nil {
		c.Download.Catalog(lang)
		catalog, err = c.Library.Catalog(lang)
		if err != nil {
			panic(err)
		}
	}
	return catalog
}

func (c *Config) SelectedCatalog() *lib.Catalog {
	return c.Catalog(c.SelectedLanguage())
}
