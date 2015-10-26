/*
 * 
 */

package main

import (
	
)

type JSONConnection struct {
	c *Content
}

func NewJSONConnection(content *Content) *JSONConnection {
	j := new(JSONConnection)
	j.c = content
	return j
}