// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// FileList represents a list of files in memory.
//
// Each file contains a Metadata and a slice of bytes that represents the
// content, so FileList encapsulates these two informations as a file, using
// two internal slices.
//
// FileList implements sort.Interface, so it is sortable using sort.Sort.
type FileList struct {
	metas    []Metadata
	contents [][]byte
}

// Append adds a file to the internal list.
//
// The file is identified by its metadata and content.
func (fl *FileList) Append(m Metadata, content []byte) {
	fl.metas = append(fl.metas, m)
	fl.contents = append(fl.contents, content)
}

// Len is the number of elements in the list.
func (fl *FileList) Len() int {
	return len(fl.metas)
}

// Less returns whether the element with index x should sort before the
// element with index y.
func (fl *FileList) Less(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	length := fl.Len()
	if x > length || y > length {
		return false
	}
	meta1 := fl.metas[x]
	meta2 := fl.metas[y]
	return meta1.Date.Before(meta2.Date.Time)
}

// Swap swaps the elements with indexes x and y.
func (fl *FileList) Swap(x, y int) {
	length := fl.Len()
	if x >= 0 && x <= length && y >= 0 && y <= length && x != y {
		fl.metas[x], fl.metas[y] = fl.metas[y], fl.metas[x]
		fl.contents[x], fl.contents[y] = fl.contents[y], fl.contents[x]
	}
}
