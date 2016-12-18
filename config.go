package main

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

// Represents the configuration options, flags, and values provided to direct
// operation.
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

// Defines the long and short version of argument options hat can be passed
// to the program as well as how to parse the following string as a value.
// Also defines a name used in config files and a default value.
// Any value can be ommitted, except fo the name.
type ConfigOption struct {
	Name     string
	Default  interface{}
	ShortArg rune
	LongArg  string
	Parse    func(string) (interface{}, error)
}

// Defines the long and short version of argument flags that can be passed
// to the program as well as the action to perform if the flag is encountered.
// Either of the ShortArg or LongArg may be permitted.
type ConfigFlag struct {
	ShortArg rune
	LongArg  string
	Action   func(*Configuration) error
}

func (o ConfigOption) needValue() bool {
	return true
}

func (o ConfigOption) handleValue(s string, c *Configuration) (err error) {
	var val interface{}
	if o.Parse != nil {
		val, err = o.Parse(s)
	} else {
		val = s
	}
	if err == nil {
		c.Set(o.Name, val)
	}
	return
}

func (o ConfigFlag) needValue() bool {
	return false
}

func (o ConfigFlag) handleValue(s string, c *Configuration) error {
	return o.Action(c)
}

func newConfiguration() *Configuration {
	return &Configuration{
		values:      make(map[string]interface{}),
		shortParams: make(map[rune]configParam),
		longParams:  make(map[string]configParam),
	}
}

// Gets the options for starting the program and sets them appropriately.
// First loads the default values, then overrites and adds values from the
// config file, finnaly overrites and adds values from command line paramters.
func (c *Configuration) Init() error {
	//TODO: Use errors instead of panics
	if err := c.loadDefaults(); err != nil {
		return err
	}
	if err := c.loadFile(); err != nil {
		return err
	}
	if err := c.loadParams(os.Args[1:]); err != nil {
		return err
	}
	return nil
}

// Defines a new ConfigFlag to parse from the config file and command line
// paramters before calling Init()
func (c *Configuration) RegisterFlag(o ConfigFlag) {
	c.shortParams[o.ShortArg] = o
	c.longParams[o.LongArg] = o
}

// Defines a new ConfigOption to parse from the config file and command line
// paramters before calling Init()
func (c *Configuration) RegisterOption(o ConfigOption) {
	c.shortParams[o.ShortArg] = o
	c.longParams[o.LongArg] = o
	c.Set(o.Name, o.Default)
}

// An ordered list of arguments passed to the program that were not flags or options
func (c *Configuration) Args() []string {
	return c.args
}

// Sets a value by it's string name
func (c *Configuration) Set(name string, value interface{}) {
	c.values[name] = value
}

// Gets a value by it's string name
func (c *Configuration) Get(name string) interface{} {
	return c.values[name]
}

// Gives a string representation of the configuration options, formatted for
// command line output
func (c *Configuration) String() string {
	buffer := new(bytes.Buffer)
	nameLen := 0

	for key := range c.values {
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

func (c *Configuration) loadDefaults() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	c.Set("Language", "eng")
	c.Set("DataDirectory", path.Join(currentUser.HomeDir, ".ldss"))
	c.Set("ServerURL", "https://tech.lds.org/glweb")
	return nil
}

func (c *Configuration) loadFile() error {
	file, err := os.Open(path.Join(c.Get("DataDirectory").(string), "config.json"))
	if err != nil {
		return nil
	}
	return json.NewDecoder(file).Decode(c.values)
}

func (c *Configuration) loadParams(args []string) error {
	for i := 0; i < len(args); {
		arg := args[i]
		if arg[0] == '-' {
			var op configParam
			var ok bool
			if arg[1] == '-' {
				if op, ok = c.longParams[arg[2:]]; !ok {
					return fmt.Errorf("Argument \"%v\" invalid", arg)
				}
			} else {
				for j := 1; j < len(arg); j++ {
					r, _ := utf8.DecodeRuneInString(arg[j:])
					op, ok = c.shortParams[r]
					if !ok {
						return fmt.Errorf("Argument \"-%v\" invalid", arg[j])
					}
					if op.needValue() && j != len(arg)-1 {
						return fmt.Errorf("Argument \"-%v\" needs a value", arg[j])
					}
					if err := op.handleValue("", c); err != nil {
						return err
					}
				}
			}
			if op.needValue() {
				if i == len(args)-1 {
					return fmt.Errorf("Argument \"%v\" needs a value", arg)
				}
				if err := op.handleValue(args[i+1], c); err != nil {
					return err
				}
				args = args[:i+copy(args[i:], args[i+2:])]
			} else {
				op.handleValue("", c)
				args = args[:i+copy(args[i:], args[i+1:])]
			}

		} else {
			i++
		}
	}
	c.args = args
	return nil
}
