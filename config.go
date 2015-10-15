package main

import (
	"os/user"
	"path"
)

type Config struct {
	OnlineContent *LDSContent
	OfflineContent *LocalContent
	Languages *LanguageLoader
	Download *Downloader
}

func LoadConfiguration() Config {
	c := Config{}
	u, err := user.Current()
	
	if err != nil {
		panic(err)
	}
	
	c.OnlineContent = new(LDSContent)
	c.OnlineContent.BasePath = "https://tech.lds.org/glweb"
	c.OfflineContent = new(LocalContent)
	c.OfflineContent.BasePath = path.Join(u.HomeDir, ".ldss")
	
	c.Languages = new(LanguageLoader)
	c.Languages.c = c.OfflineContent
	c.Download = new(Downloader)
	c.Download.online = c.OnlineContent
	c.Download.offline = c.OfflineContent
	
	return c
}
