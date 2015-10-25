package main

import (
	"encoding/json"
	"fmt"
)

type JSONLanguageLoader struct {
	content   Content
	languages []Language
}

type LanguageLoader interface {
	GetByUnknown(id string) *Language
	GetAll() []Language
}

func NewJSONLanguageLoader(c Content) *JSONLanguageLoader {
	l := new(JSONLanguageLoader)
	l.content = c
	return l
}

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

func (l *JSONLanguageLoader) populateIfNeeded() {
	if l.languages != nil {
		return
	}

	var description glLanguageDescription
	file := l.content.OpenRead(l.content.GetLanguagesPath())
	dec := json.NewDecoder(file)
	err := dec.Decode(&description)
	if err != nil {
		panic(err)
	}

	l.languages = description.Languages
}

func (l *JSONLanguageLoader) GetByUnknown(id string) *Language {
	l.populateIfNeeded()
	for _, lang := range l.languages {
		if lang.Name == id || fmt.Sprintf("%v", lang.ID) == id || lang.EnglishName == id || lang.Code == id || lang.GlCode == id {
			return &lang
		}
	}
	return nil
}

func (l *JSONLanguageLoader) GetAll() []Language {
	l.populateIfNeeded()
	return l.languages
}
