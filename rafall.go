package rafall

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
)

func GenerateHtmlFromFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	return blackfriday.MarkdownCommon(content), err
}
