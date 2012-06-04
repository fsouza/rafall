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

// Time is used to parse dates in RFC822 format from JSON.
type Time struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface, formatting time
// using RFC822 format.
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"` + time.RFC822Z + `"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaller interface, formatting time
// using RFC822 format.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var tim time.Time
	if tim, err = time.Parse(`"`+time.RFC822Z+`"`, string(data)); err == nil {
		*t = Time{tim}
	}
	return
}

// Metadata represents the metadata of a file.
//
// Each file should declare its metadata in JSON format in the beggining of
// the file, using HTML comments. For example:
//
//     <!--{
//     "Title": "My post",
//     "Date": "04 Jun 12 13:56 -0300",
//     "Tags": ["post", "blog"]
//     }-->
//
// Notice that this code should be the first thing in the code.
type Metadata struct {
	Title string
	Date  Time
	Tags  []string
}

// extractMetadata extracts Metadata from a slice of bytes.
//
// The slice is supposed to be the content of the file. It returns the
// content minus the metadata, the metadata and an error, if any happens.
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

// readConfig reads config from a slice of bytes. The slice of bytes should
// be in JSON format. It is supposed to be the content from the config file.
func readConfig(content []byte) (config map[string]string, err error) {
	err = json.Unmarshal(content, &config)
	return
}

// Generator represents the site generator.
//
// This type encapsulates the site configuration, and provide the methods
// that should be used to (re)generate the site.
type Generator struct {
	config    map[string]string
	metaFiles map[string]string
}

// NewGenerator returns a new generator instance, and an error, if any.
//
// It receives the configuration file name. It should be a path, that can be
// absolute or relative. If the file does not exist, or is not in JSON
// format, or other kind of unexpected error happen, this function returns
// error.
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

// isMetaFile indicates whether the given filename represents a metafile. A
// metafile is an special HTML file that can not be used as post. For
// example, layout.html is the metafile that defines the layout of the site.
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

// collectFiles collect files in the src directory, and returns a FileList
// and an error, if any happens. The FileList contains all files contents
// and metadata, and does not include metafiles.
//
// After using this method to collect the list of files, the generator
// should render all files and store the output.
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
