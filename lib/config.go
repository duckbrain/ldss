package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"unicode/utf8"
)

var c *Configuration

func Config() *Configuration {
	return c
}

type Configuration struct {
	args        []string
	values      map[string]interface{}
	shortParams map[rune]configParam
	longParams  map[string]configParam
}

type configParam interface {
	needValue() bool
	handleValue(string, *Configuration) error
}

type AppOption struct {
	Name     string
	Default  interface{}
	ShortArg rune
	LongArg  string
	Parse    func(string) (interface{}, error)
}

type AppFlag struct {
	ShortArg rune
	LongArg  string
	Action   func(*Configuration) error
}

func (o AppOption) needValue() bool {
	return true
}

func (o AppOption) handleValue(s string, c *Configuration) error {
	val, err := o.Parse(s)
	if err == nil {
		c.Set(o.Name, val)
	}
	return err
}

func (o AppFlag) needValue() bool {
	return false
}

func (o AppFlag) handleValue(s string, c *Configuration) error {
	return o.Action(c)
}

func init() {
	c = &Configuration{
		values:      make(map[string]interface{}),
		shortParams: make(map[rune]configParam),
		longParams:  make(map[string]configParam),
	}
	c.loadDefaults()
	c.loadFile()
	c.loadParams()
}

func (c *Configuration) RegisterFlag(o AppFlag) {
	c.shortParams[o.ShortArg] = o
	c.longParams[o.LongArg] = o
}

func (c *Configuration) RegisterOption(o AppOption) {
	c.shortParams[o.ShortArg] = o
	c.longParams[o.LongArg] = o
	c.Set(o.Name, o.Default)
}

func (c *Configuration) Set(name string, value interface{}) {
	c.values[name] = value
}

func (c *Configuration) Get(name string) interface{} {
	return c.values[name]
}

func (c *Configuration) String() string {
	var buffer bytes.Buffer
	nameLen := 0

	for key, _ := range c.values {
		if len(key) > nameLen {
			nameLen = len(key)
		}
	}

	for key, value := range c.values {
		spaces := strings.Repeat(" ", nameLen-len(key)+1)
		buffer.WriteString(fmt.Sprintf("%v:%v%v\n", key, spaces, value))
	}

	return buffer.String()
}

func (c *Configuration) loadDefaults() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	c.Set("Language", "eng")
	c.Set("DataDirectory", path.Join(currentUser.HomeDir, ".ldss"))
	c.Set("ServerURL", "https://tech.lds.org/glweb")
}

func (c *Configuration) loadFile() {
	file, err := os.Open(path.Join(c.Get("DataDirectory").(string), "config.json"))
	if err != nil {
		return
	}
	if err = json.NewDecoder(file).Decode(c.values); err != nil {
		panic(err)
	}
}

func (c *Configuration) loadParams() {
	args := os.Args[1:]
	for i := 0; i < len(args); {
		arg := args[i]
		if arg[0] == '-' {
			var op configParam
			var ok bool
			if arg[1] == '-' {
				if op, ok = c.longParams[arg[2:]]; !ok {
					panic(fmt.Errorf("Argument \"%v\" invalid", arg))
				}
			} else {
				for j := 1; j < len(arg); j++ {
					r, _ := utf8.DecodeRuneInString(arg[j:])
					op, ok = c.shortParams[r]
					if !ok {
						panic(fmt.Errorf("Argument \"-%v\" invalid", arg[j]))
					}
					if op.needValue() && j != len(arg)-1 {
						panic(fmt.Errorf("Argument \"-%v\" needs a value", arg[j]))
					}
					op.handleValue("", c)
				}
			}
			if op.needValue() {
				if i == len(args)-1 {
					panic(fmt.Errorf("Argument \"%v\" needs a value", arg))
				}
				op.handleValue(args[i+1], c)
				args = args[:i+copy(args[i:], args[i+2:])]
			} else {
				op.handleValue("", c)
				args = args[:i+copy(args[i:], args[i+1:])]
			}

		} else {
			i++
		}
	}
}
