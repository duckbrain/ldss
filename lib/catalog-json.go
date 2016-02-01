package lib

type jsonCatalogBase struct {
	Catalog         *jsonFolder `json:"catalog"`
	CoverArtBaseUrl string      `json:"cover_art_base_url"`
}

type jsonFolder struct {
	Name    string        `json:"name"`
	Folders []*jsonFolder `json:"folders"`
	Books   []*jsonBook   `json:"books"`
	ID      int           `json:"id"`
}

type jsonBook struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	GlURI string `json:"gl_uri"`
}
