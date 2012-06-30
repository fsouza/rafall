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

// File represents a file in memory.
//
// Each file contains a Metadata and a slice of byte that represents the
// content.
type File struct {
	meta    Metadata
	content []byte
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

// Iter returns an instance of the Iter type, and should be used to iterate the
// list.
func (fl *FileList) Iter() *Iter {
	ch := make(chan File)
	go func(fl *FileList, fch chan File) {
		for i, meta := range fl.metas {
			fch <- File{
				meta:    meta,
				content: fl.contents[i],
			}
		}
		close(fch)
	}(fl, ch)
	return &Iter{
		fch: ch,
	}
}

// Iter represents a iterator that iterates through a filelist.
//
// To get an interator, you should call the Iter method on FileList type.
type Iter struct {
	fch chan File
}

// Next returns the next element in the iterator, or false in the third return
// value, indicating that there are no more elements to iterate over.
func (i *Iter) Next() (meta Metadata, content []byte, present bool) {
	var f File
	f, present = <-i.fch
	meta = f.meta
	content = f.content
	return
}
