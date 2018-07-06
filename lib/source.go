package lib

type Source interface {
	Downloader
	Langs() ([]*Lang, error)
	Root(lang *Lang) (Item, error)
}

var srcs map[string]Source

func init() {
	srcs = make(map[string]Source)
}

func Register(name string, src Source) {
	if _, ok := srcs[name]; ok {
		panic(fmt.Errorf("Cannot have two sources with name %v", name))
	}
	srcs[name] = src
}
