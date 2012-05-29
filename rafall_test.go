// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func buildGenerator() (g Generator) {
	g.config = map[string]string{
		"siteName":        "Rafall",
		"subtitle":        "Random stuff",
		"description":     "hi, my name is rafall",
		"disqusShortname": "rafall",
	}
	g.metaFiles = map[string]string{
		"archive": "archive.html",
		"layout":  "layout.html",
		"post":    "post.html",
	}
	return
}

func getContent(filename string) []byte {
	c, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		panic(err)
	}
	return c
}

func TestTimeMarshalJSON(t *testing.T) {
	expected := []byte(`"28 May 12 02:00 -0300"`)
	inputTime, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-28 02:00:00 -0300")
	if err != nil {
		t.Error(err)
	}
	tim := Time{inputTime}
	got, _ := tim.MarshalJSON()
	if bytes.Compare(expected, got) != 0 {
		t.Errorf("Failed to marshal time as json.\nExpected: %q\nGot: %q", expected, got)
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	tim, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-28 02:00:00 -0300")
	if err != nil {
		t.Error(err)
	}
	expected := Time{tim}
	input := []byte(`"28 May 12 02:00 -0300"`)
	got := Time{}
	err = got.UnmarshalJSON(input)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Failed to unmarshal %s.\nExpected: %q\nGot: %q", string(input), expected, got)
	}
}

func TestExtractMetadata(t *testing.T) {
	expectedTime, _ := time.Parse(time.RFC822Z, "27 May 12 01:50 -0300")
	expected := Metadata{
		Title: "Hello world",
		Tags:  []string{"post"},
		Date:  Time{expectedTime},
	}
	content := getContent("hello_world.html")
	content, metadata, err := extractMetadata(content)
	if err != nil {
		t.Error(nil)
	}
	if !reflect.DeepEqual(metadata, expected) {
		t.Errorf("Expected metadata: %q\nGot metadata: %q", expected, metadata)
	}
}

func TestExtractMetadataFromAFileThatDoesNotHaveMetadata(t *testing.T) {
	content := getContent("two_paragraphs.html")
	content, _, err := extractMetadata(content)
	if err != nil {
		t.Error(nil)
	}
}

func TestExtractMetadataReturnsTail(t *testing.T) {
	content := getContent("hello_world.html")
	content, _, _ = extractMetadata(content)
	expected := getContent("hello_world_without_metadata.html")
	if bytes.Compare(content, expected) != 0 {
		t.Errorf("Should extract metadata and return tail. Expected: %q\nGot: %q", expected, content)
	}
}

func TestNewGenerator(t *testing.T) {
	g, err := NewGenerator("testdata/config")
	if err != nil {
		t.Error(err)
	}
	if g.config["siteName"] != "Rafall" {
		t.Errorf("Should read config from given config file")
	}
}

func TestNewGeneratorReturnsErrorIfTheFileDoesNotExist(t *testing.T) {
	_, err := NewGenerator("something/that/does/not/exist")
	if err == nil {
		t.Error("Should return error if the file does not exist, but returned nil")
	}
}

func TestReadConfig(t *testing.T) {
	conf := getContent("config")
	expected := map[string]string{
		"siteName":        "Rafall",
		"subtitle":        "Random stuff",
		"description":     "hi, my name is rafall",
		"disqusShortname": "rafall",
	}
	config, err := readConfig(conf)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, config) {
		t.Errorf("Should read config from bytes.\nExpected: %q\nGot: %q", expected, config)
	}
}

func TestReadConfigReturnsErrorIfJSONIsInvalid(t *testing.T) {
	_, err := readConfig([]byte("invalid;json:"))
	if err == nil {
		t.Error("Should return error if the json is invalid, but returned nil")
	}
}

func TestIsMetaFile(t *testing.T) {
	input := map[string]bool{
		"archive.html":   true,
		"layout.html":    true,
		"post.html":      true,
		"francisco.html": false,
		"something.xml":  false,
	}
	g := buildGenerator()
	for k, v := range input {
		if g.isMetaFile(k) != v {
			t.Errorf("Should %s be metafile?\nExpected: %q\nGot: %q", k, v, !v)
		}
	}
}

func TestIsValid(t *testing.T) {
	input := map[string]bool{
		"archive.html":   false,
		"abc.xml":        false,
		"francisco.html": true,
	}
	g := buildGenerator()
	for k, v := range input {
		if g.isValid(k) != v {
			t.Errorf("Should %s be valid?\nExpected: %q\nGot: %q", k, v, !v)
		}
	}
}

func TestCollectFiles(t *testing.T) {
	g := buildGenerator()
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	defer os.Chdir(cwd)
	os.Chdir("testdata/sampleproject")
	fl, err := g.collectFiles()
	if err != nil {
		t.Error(err)
	}
	time1, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-28 23:56:00 -0300")
	if err != nil {
		t.Error(err)
	}
	time2, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-29 12:00:00 -0300")
	if err != nil {
		t.Error(err)
	}
	expectedMetas := []Metadata{
		Metadata{
			Title: "Hello world",
			Date:  Time{time1},
			Tags:  []string{"hello", "world"},
		},
		Metadata{
			Title: "Good bye cruel world",
			Date:  Time{time2},
			Tags:  []string{"goodbye", "world"},
		},
	}
	if !reflect.DeepEqual(expectedMetas, fl.metas) {
		t.Errorf("Should return all metadatas as expected, sorted by date.\nExpected: %q\nGot: %q", expectedMetas, fl.metas)
	}
}
