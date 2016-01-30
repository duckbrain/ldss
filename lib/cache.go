package lib

import (
	"sync"
)

type cache struct {
	lock      sync.Mutex
	val       interface{}
	construct func() (interface{}, error)
}

func (c *cache) get() (interface{}, error) {
	if c.val != nil {
		return c.val, nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.val != nil {
		return c.val, nil
	}
	val, err := c.construct()
	if err != nil {
		return nil, err
	}
	c.val = val
	return val, nil
}
