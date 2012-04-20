package rafall

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

var testfilenames = []string{"hello_world", "two_paragraphs"}

func getPathAndExpectedOutput(filename string) (string, []byte) {
	path, _ := filepath.Abs("testdata/" + filename + ".mkd")
	output, _ := ioutil.ReadFile("testdata/" + filename + ".html")
	return path, output
}

func TestGenerateHtmlFromFile(t *testing.T) {
	for _, filename := range testfilenames {
		path, expected := getPathAndExpectedOutput(filename)
		got, _ := GenerateHtmlFromFile(path)
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected:\n%s\nfor the file %s, but got:\n%s", string(expected), filename, string(got))
		}
	}
}

func TestGenerateHtmlFromFileReturnErrorsWhenTheFileDoesNotExist(t *testing.T) {
	content, err := GenerateHtmlFromFile("/some/path/that/should/not/exist")
	if content != nil {
		t.Errorf("Should return nil when the file does not exist, returned: %q", content)
	}
	if err == nil {
		t.Errorf("Should return error when the file does not exist, returned: %q", err)
	}
}
