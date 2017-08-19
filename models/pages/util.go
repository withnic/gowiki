package pages

import (
	"errors"
	"io/ioutil"
	"regexp"

	"fmt"
)

type UtilFuncs interface {
	Load(title string) (*Page, error)
}

var Util UtilFuncs = utilFuncs{}

type utilFuncs struct{}

var validateTitleRegexp = regexp.MustCompile(`^([a-zA-Z0-9]+)$`)

var dataDirPattern = "data/%s.txt"

func getDataPath(str string) string {
	return fmt.Sprintf(dataDirPattern, str)
}

func validateTitle(title string) bool {
	return validateTitleRegexp.MatchString(title)
}

func readBody(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Load returns Page struct and error
func (u utilFuncs) Load(title string) (*Page, error) {
	if !validateTitle(title) {
		return nil, errors.New("Invalid title. Title is only alphabet.")
	}

	filename := getDataPath(title)
	body, err := readBody(filename)

	return &Page{
		Title: title,
		Body:  body,
	}, err
}
