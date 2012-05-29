// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
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

func extractMetadata(content []byte) ([]byte, *Metadata, error) {
	if !bytes.HasPrefix(content, jsonStart) {
		return content, nil, nil
	}
	tail := content
	meta := new(Metadata)
	end := bytes.Index(content, jsonEnd)
	content = content[len(jsonStart)-1 : end+1]
	err := json.Unmarshal(content, meta)
	return tail[end+len(jsonEnd):], meta, err
}

func readConfig(content []byte) (config map[string]string, err error) {
	err = json.Unmarshal(content, &config)
	return
}

type Generator struct {
	config map[string]string
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
	return
}
