// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rafall

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
)

func GenerateHtmlFromFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	return blackfriday.MarkdownCommon(content), err
}
