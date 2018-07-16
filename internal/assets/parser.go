package assets

import "github.com/duckbrain/ldss/lib"

func init() {
	lib.SetReferenceParseReader(func(lang lib.Lang) ([]byte, error) {
		return Asset("data/reference/" + lang.Code())
	})
}
