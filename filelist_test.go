package main

import (
	"bytes"
	"testing"
	"time"
)

func TestFileListShouldReturnTheLengthOfTheInnerSlices(t *testing.T) {
	m := Metadata{}
	content := []byte("hi")
	list := FileList{}
	list.Append(m, content)
	if list.Len() != 1 {
		t.Errorf("Should have length 1 after append items")
	}
}

func TestLessShouldReturnTrueIfTheFirstDateIsLesserThanTheSecond(t *testing.T) {
	time1, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-28 23:56:00 -0300")
	if err != nil {
		t.Error(err)
	}
	time2, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-29 12:00:00 -0300")
	if err != nil {
		t.Error(err)
	}
	list := FileList{}
	list.metas = []Metadata{
		Metadata{
			Title: "Python get's old",
			Date:  Time{time1},
			Tags:  []string{"python"},
		},
		Metadata{
			Title: "Gopher got a lady",
			Date:  Time{time2},
			Tags:  []string{"gopher"},
		},
	}
	list.contents = [][]byte{[]byte("hi"), []byte("there!")}
	if !list.Less(0, 1) {
		t.Errorf("Item 0 should be lesser than item 1 in the list")
	}
	if list.Less(1, 0) {
		t.Errorf("Item 1 should not be lesses than item 0 in the list")
	}
}

func TestLessShouldReturnFalseIfXOrYIsGreaterThanLen(t *testing.T) {
	list := FileList{}
	if list.Less(1, 1) {
		t.Errorf("Should return false if x or y is greater than length")
	}
}

func TestLessShouldReturnFalseIfXOrYIsNegative(t *testing.T) {
	list := FileList{}
	if list.Less(-1, 0) {
		t.Errorf("Should return false if x or y is negative")
	}
}

func TestSwapShouldDoNothingIfXOrYIsGreaterThanLen(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not fail, but failed with: %s", r)
		}
	}()
	list := FileList{}
	list.Swap(1, 2)
}

func TestSwapShoulDoNothingIfXOrYIsNegative(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not fail, but failed with: %s", r)
		}
	}()
	list := FileList{}
	list.Swap(-1, 0)
}

func TestSwapShouldDoNothingIfXIsEqualToY(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not fail, but failed with: %s", r)
		}
	}()
	m := Metadata{}
	content := []byte("hi")
	list := FileList{}
	list.Append(m, content)
	list.Swap(0, 0)
}

func TestSwapShouldSwapBothInnerSlices(t *testing.T) {
	time1, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-28 23:56:00 -0300")
	if err != nil {
		t.Error(err)
	}
	time2, err := time.Parse("2006-01-02 15:04:05 -0700", "2012-05-29 12:00:00 -0300")
	if err != nil {
		t.Error(err)
	}
	list := FileList{}
	list.metas = []Metadata{
		Metadata{
			Title: "Python get's old",
			Date:  Time{time1},
			Tags:  []string{"python"},
		},
		Metadata{
			Title: "Gopher got a lady",
			Date:  Time{time2},
			Tags:  []string{"gopher"},
		},
	}
	list.contents = [][]byte{[]byte("hi"), []byte("there!")}
	list.Swap(0, 1)
	if list.metas[0].Title == "Python get's old" {
		t.Errorf("Should swap items in the metadata slice")
	}
	if bytes.Compare(list.contents[0], []byte("hi")) == 0 {
		t.Errorf("Should swap items in the contents slice")
	}
}
