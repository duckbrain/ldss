package ldslib

import "fmt"

type glLanguageDescription struct {
	Languages []Language `json:"languages"`
	Success   bool       `json:"success"`
}

type Language struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	EnglishName string `json:"eng_name"`
	Code        string `json:"code"`
	GlCode      string `json:"code_three"`
}

func (l *Language) String() string {
	var id, name, code string

	id = fmt.Sprintf("%v: ", l.ID)
	if l.Name == l.EnglishName {
		name = l.Name
	} else {
		name = fmt.Sprintf("%v (%v)", l.Name, l.EnglishName)
	}
	if l.Code == l.GlCode {
		code = fmt.Sprintf(" [%v]", l.Code)
	} else {
		code = fmt.Sprintf(" [%v/%v]", l.Code, l.GlCode)
	}

	return id + name + code
}
