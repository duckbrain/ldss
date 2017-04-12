package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/duckbrain/ldss/lib"
	"github.com/spf13/cobra"
)

// referenceCmd represents the reference command
var referenceCmd = &cobra.Command{
	Use:   "reference",
	Short: "Generates a reference parsing file for a given language",
	Long:  `This command is likely only useful for development. When adding support for parsing a new language, this command will generate a file that provides a basic parsing file that can be used as a basis for the customized parsing file.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := &generateReference{lang()}
		app.run()
	},
}

func init() {
	generateCmd.AddCommand(referenceCmd)
}

type generateReference struct {
	lang *lib.Lang
}

func (app *generateReference) lookup(path string) lib.Item {
	ref := lib.ParsePath(app.lang, path)
	item, err := ref.Lookup()
	if err != nil {
		panic(err)
	}
	return item
}

func (app *generateReference) run() {
	err := lib.DownloadAll(app.lang, false)
	if err != nil {
		panic(err)
	}

	app.runScriptureVolume("/scriptures/ot")
	app.runScriptureVolume("/scriptures/nt")
	app.runScriptureVolume("/scriptures/bofm")
	app.runScriptureVolume("/scriptures/pgp")
	app.runDandC(app.lookup("/scriptures/dc-testament").(*lib.Book))
}

// Generates lookup names from user readable strings
func (app *generateReference) userNames(name string) []string {
	name = strings.ToLower(name)
	return []string{name}
}

// Generates lookup names from the last component of a path
func (app *generateReference) pathNames(name string) []string {
	name = name[strings.LastIndex(name, "/")+1:]
	name = strings.ToLower(name)
	name = strings.Replace(name, "-", " ", 100)
	return []string{name}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func (app *generateReference) genSimple(n lib.Item, hash string) {
	app.gen(append(app.userNames(n.Name()), app.pathNames(n.Path())...), n.Path()+hash)
}

func (app *generateReference) genParent(n lib.Item, hash string) {
	p := n.Parent()
	names := []string{}
	pnames := append(app.userNames(p.Name()), app.pathNames(p.Path())...)
	nnames := append(app.userNames(n.Name()), app.pathNames(n.Path())...)
	for _, pname := range pnames {
		for _, nname := range nnames {
			names = append(names, pname+" "+nname)
		}
	}
	app.gen(names, n.Path())
}

func (app *generateReference) gen(matches []string, path string) {
	cleaned := []string{}
	for _, value := range matches {
		if !stringInSlice(value, cleaned) {
			cleaned = append(cleaned, value)
		}
	}
	fmt.Printf("%v:%v\n", strings.Join(cleaned, ":"), path)
}

func (app *generateReference) comment(comment string) {
	fmt.Printf("#%v\n", comment)
}

func (app *generateReference) runScriptureVolume(path string) {
	b, err := lib.ParsePath(app.lang, path).Lookup()
	if err != nil {
		return
	}
	fmt.Println("")
	app.comment(b.Name())

	app.genSimple(b, "")
	nodes, err := b.Children()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, n := range nodes {
		app.runScriptureBook(n)
	}
}

func (app *generateReference) runScriptureBook(n lib.Item) {
	if _, err := lib.ParsePath(app.lang, n.Path()+"/1").Lookup(); err == nil {
		if _, err := lib.ParsePath(app.lang, n.Path()+"/2").Lookup(); err != nil {
			// Is a single chapter book
			names := app.userNames(n.Name())
			names = append(names, app.userNames(n.Name()+" 1")...)
			names = append(names, app.pathNames(n.Path())...)
			names = append(names, app.pathNames(n.Path()+"-1")...)
			app.gen(names, n.Path())
		} else {
			// Is a multiple chapter book
			app.genSimple(n, "#")
		}
	} else {
		app.genParent(n, "")
	}
}

func (app *generateReference) runDandC(n lib.Item) {
	fmt.Println("")
	app.comment(n.Name())
	app.genSimple(n, "#")
	//TODO Generate the number regex
	children, err := n.Children()
	if err != nil {
		panic(err)
	}
	for _, c := range children {
		path := c.Path()
		path = path[strings.LastIndex(path, "/")+1:]
		if _, err := strconv.Atoi(path); err != nil {
			// Print the ones that don't end in numbers
			app.genParent(c, "")
		}
	}
}
