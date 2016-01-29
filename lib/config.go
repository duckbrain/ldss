package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"bytes"
	"strings"
	"errors"
)

// Package configuration object loaded on package load
var Config ConfigurationOptions

type DefaultCommands int

const (
	DefaultCommandLookup DefaultCommands = iota
	DefaultCommandCurses
)

/*type Config struct {
	op             ConfigurationOptions
	OnlineContent  Source
	OfflineContent Source
	Library        *Library
	Reference      RefParser
	Download       *Downloader
	DefaultCommand DefaultCommands
}*/

type ConfigurationOptions struct {
	Language        string
	DataDirectory   string
	ServerURL       string
	//WebPort         int TODO: Move to ldss/web.go 
	//WebTemplatePath string
	AppParams	map[string]interface{}
}

type AppOption interface {
	Name() string
	Default() interface{}
	Params() (rune, string) // Return the rune for -abc params and the string for --full-name params, omit the starting hyphens
	ParseParam(string) (interface{}, error)
}

func (op ConfigurationOptions) String() string {
	var buffer bytes.Buffer
	nameLen := 13

	for key, _ := range op.AppParams {
		if len(key) > nameLen {
			nameLen = len(key)
		}
	}

	printLn := func(key string, value interface{}) {
		spaces := strings.Repeat(" ", nameLen - len(key) + 1)
		buffer.WriteString(fmt.Sprintf("%v:%v%v\n", key, spaces, value));
	}

	printLn("Language", op.Language)
	printLn("ServerURL", op.ServerURL)
	printLn("DataDirectory", op.DataDirectory)

	for key, value := range op.AppParams {
		printLn(key, value);
	}

	return buffer.String()
}

func (c *ConfigurationOptions) MarshalJSON() ([]byte, error) {
	cmap := make(map[string]interface{})
	cmap["Language"] = c.Language
	cmap["DataDirectory"] = c.DataDirectory
	cmap["ServerURL"] = c.ServerURL

	for k, v := range c.AppParams {
		cmap[k] = v
	}

	return json.Marshal(cmap)
}

func (c *ConfigurationOptions) UnMarshalJSON (data []byte) (error) {
	var cmap map[string]interface{}

	if c == nil {
	    return errors.New("RawString: UnmarshalJSON on nil pointer")
	}

	if err := json.Unmarshal(data, &cmap); err != nil {
	    return err
	}

	for key, val := range cmap {
	    switch key {
		case "Language":
		    c.Language = val.(string)
		case "DataDirectory":
		    c.DataDirectory = val.(string)
		case "ServerURL":
		    c.ServerURL = val.(string)
		default:
		    c.AppParams[key] = val
	    }
	}
	return nil
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

	return op
}

func loadParameterOptions(op ConfigurationOptions) (ConfigurationOptions, []string) {
	args := os.Args[1:]
	for i := 0; i < len(args); {
		switch args[i] {
		/*case "-p":
			op.WebPort, err = strconv.Atoi(args[i+1])
			if err != nil {
				panic(fmt.Errorf("Could not convert port \"%v\" to an integer", args[i+1]))
			}
			args = args[:i+copy(args[i:], args[i+2:])]*/
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

/*func LoadConfiguration(op ConfigurationOptions) Config {
	c := Config{op: op}
	c.OnlineContent = NewOnlineSource(op.ServerURL)
	c.OfflineContent = NewOfflineSource(op.DataDirectory)
	c.Library = NewLibrary(c.OfflineContent)
	c.Download = NewDownloader(c.OnlineContent, c.OfflineContent)
	return c
}

func (c *Config) Languages() []Language {
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

func (c *Config) Language(s string) *Language {
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

func (c *Config) SelectedLanguage() *Language {
	return c.Language(c.op.Language)
}

func (c *Config) Catalog(lang *Language) *Catalog {
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

func (c *Config) SelectedCatalog() *Catalog {
	return c.Catalog(c.SelectedLanguage())
}*/
