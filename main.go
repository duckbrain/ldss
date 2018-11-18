// ldss is a program to download and read the scriptures from the Church of Jesus Christ
// of Latter-day Saints in a variety of formats. It contains a web server, GUI, and
// command line interface. The content parsing and lookup is implemented in the ldss/lib
// subpackage under this package.
package main

import (
	"github.com/duckbrain/ldss/cmd"
	_ "github.com/duckbrain/ldss/lib/sources/ldsorg"
)

func main() {
	cmd.Execute()
}
