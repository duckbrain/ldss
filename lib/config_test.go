package lib

import (
	"testing"
)

func TestConfig(t *testing.T) {
	c := newConfiguration()
	c.Set("myvalue", 123)
	if c.Get("myvalue").(int) != 123 {
		t.Fail()
	}
	c.Set("myvalue", "hello")
	if c.Get("myvalue").(string) != "hello" {
		t.Fail()
	}
	t.Logf("TestConfig: %v", c)
}

func TestConfigParams(t *testing.T) {
	c := newConfiguration()
	c.RegisterFlag(AppFlag{
		ShortArg: 'f',
		LongArg:  "--four",
		Action: func(c *Configuration) error {
			c.Set("Four", true)
			return nil
		},
	})
	if err := c.loadParams([]string{"one", "two", "three", "-four"}); err == nil {
		t.Fail()
	}
	if err := c.loadParams([]string{"one", "two", "three", "-f"}); err != nil {
		t.Fail()
	}
	if len(c.args) != 3 {
		t.Error(c.args)
	}

	t.Log(c, c.args)
}

func TestConfigInit(t *testing.T) {
	c := newConfiguration()
	if err := c.Init(); err != nil {
		t.Error(err, c)
	}
}
