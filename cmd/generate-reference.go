// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		lang := lang()
		cat, err := lang.Catalog()
		if err != nil {
			panic(err)
		}
		app := &generateReference{lang, cat, args}
		app.run()
	},
}

func init() {
	generateCmd.AddCommand(referenceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// referenceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// referenceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

type generateReference struct {
	lang *lib.Lang
	cat  *lib.Catalog
	args []string
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
	if len(app.args) != 2 {
		panic("Invalid number of arguments, expects language id")
	}
	langID := app.args[1]
	fmt.Println(langID)
	var err error
	if app.lang, err = lib.LookupLanguage(langID); err != nil {
		panic(err)
	}
	if app.cat, err = app.lang.Catalog(); err != nil {
		panic(err)
	}
	fmt.Println(app.lang.String())
	err = lib.DownloadAll(app.lang, false)
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
