package pages

import (
	"errors"
	"html/template"
	"io/ioutil"
	"regexp"

	"fmt"

	"github.com/russross/blackfriday"
)

type UtilFuncs interface {
	Load(title string) (*Page, error)
	Parse(raw []byte) []byte
}

var Util UtilFuncs = utilFuncs{}

type utilFuncs struct{}

var validateTitleRegexp = regexp.MustCompile(`^([a-zA-Z0-9]+)$`)

var dataDirPattern = "data/%s.md"

func getDataPath(str string) string {
	return fmt.Sprintf(dataDirPattern, str)
}

func validateTitle(title string) bool {
	return validateTitleRegexp.MatchString(title)
}

func readBody(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (u utilFuncs) Parse(raw []byte) []byte {
	return blackfriday.MarkdownCommon(raw)
}

// Load returns Page struct and error
func (u utilFuncs) Load(title string) (*Page, error) {
	if !validateTitle(title) {
		return nil, errors.New("Invalid title. Title is only alphabet.")
	}

	filename := getDataPath(title)
	body, err := readBody(filename)

	output := u.Parse(body)
	return &Page{
		Title:    title,
		Body:     body,
		HtmlBody: template.HTML(output),
	}, err
}
