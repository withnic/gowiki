package pages

import (
	"errors"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
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
