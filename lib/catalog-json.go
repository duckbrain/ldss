package lib

type jsonCatalogBase struct {
	Catalog         *jsonCatalog `json:"catalog"`
	CoverArtBaseUrl string       `json:"cover_art_base_url"`
}

type jsonFolderBase struct {
	Name    string    `json:"name"`
	Folders []*Folder `json:"folders"`
	Books   []*Book   `json:"books"`
}

type jsonCatalog jsonFolderBase

type jsonFolder struct {
	jsonFolderBase
	ID int `json:"id"`
}

type jsonBook struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	GlURI string `json:"gl_uri"`
}
