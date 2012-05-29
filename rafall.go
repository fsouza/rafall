// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

var (
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
)

// Time is used to parse dates in RFC822 format from JSON
type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"` + time.RFC822Z + `"`)), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var tim time.Time
	if tim, err = time.Parse(`"`+time.RFC822Z+`"`, string(data)); err == nil {
		*t = Time{tim}
	}
	return
}

type Metadata struct {
	Title string
	Date  Time
	Tags  []string
}

func extractMetadata(content []byte) (tail []byte, meta Metadata, err error) {
	tail = content
	if !bytes.HasPrefix(content, jsonStart) {
		return
	}
	end := bytes.Index(content, jsonEnd)
	content = content[len(jsonStart)-1 : end+1]
	err = json.Unmarshal(content, &meta)
	return tail[end+len(jsonEnd):], meta, err
}

func readConfig(content []byte) (config map[string]string, err error) {
	err = json.Unmarshal(content, &config)
	return
}

type Generator struct {
	config    map[string]string
	metaFiles map[string]string
}

func NewGenerator(configFile string) (g Generator, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	g.config, err = readConfig(b)
	g.metaFiles = map[string]string{
		"archive": "archive.html",
		"layout":  "layout.html",
		"post":    "post.html",
	}
	return
}

func (g *Generator) isMetaFile(filename string) bool {
	for _, v := range g.metaFiles {
		if filename == v {
			return true
		}
	}
	return false
}

// isValid indicates whether the give filename is valid for a post
// or not. To be valid, a file must end with .html and be not a metafile.
func (g *Generator) isValid(filename string) bool {
	if !strings.HasSuffix(filename, ".html") {
		return false
	}
	return !g.isMetaFile(filename)
}

func (g *Generator) collectFiles() (fl FileList, err error) {
	dir, err := os.Open("src")
	if err != nil {
		return
	}
	defer dir.Close()
	fis, err := dir.Readdir(-1)
	if err != nil {
		return
	}
	for _, fi := range fis {
		name := fi.Name()
		if !g.isValid(name) {
			continue
		}
		p := path.Join(".", "src", name)
		f, err := os.Open(p)
		if err != nil {
			fmt.Printf("Skipping %s: %s\n", name, err.Error())
			continue
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			fmt.Printf("Skipping %s: %s\n", name, err.Error())
			continue
		}
		content, metadata, err := extractMetadata(b)
		if err != nil {
			fmt.Printf("Failed to extract metadata from %s: %s", name, err.Error())
			continue
		}
		fl.Append(metadata, content)
	}
	sort.Sort(&fl)
	return
}
