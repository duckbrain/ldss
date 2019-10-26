package file

import (
	"github.com/duckbrain/ldss/lib"
)

var _ lib.Storage = *Store(nil)

type Store struct {
}
