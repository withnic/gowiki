package pages

import (
	"errors"
	"html/template"
	"io/ioutil"
)

type Page struct {
	Title    string
	Body     []byte
	HtmlBody template.HTML
}

func (p *Page) write(path string) error {
	return ioutil.WriteFile(path, p.Body, 0600)
}

func (p *Page) Save() error {
	if !validateTitle(p.Title) {
		return errors.New("Invalid title. Title is only alphabet.")
	}

	path := getDataPath(p.Title)
	return p.write(path)
}
