// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

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
	if !reflect.DeepEqual(*metadata, expected) {
		t.Errorf("Expected metadata: %q\nGot metadata: %q", expected, metadata)
	}
}

func TestExtractMetadataFromAFileThatDoesNotHaveMetadata(t *testing.T) {
	content := getContent("two_paragraphs.html")
	content, metadata, err := extractMetadata(content)
	if err != nil {
		t.Error(nil)
	}
	if metadata != nil {
		t.Errorf("Metadata should be nil, but it is %q", metadata)
	}
}
