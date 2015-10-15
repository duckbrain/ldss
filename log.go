
package main

import "github.com/wsxiaoys/terminal/color"


func Print (err error) {
	color.Println("@rerror@{|}: " + err.Error())
}
