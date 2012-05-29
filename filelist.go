package main

type FileList struct {
	metas    []Metadata
	contents [][]byte
}

func (fl *FileList) Append(m Metadata, content []byte) {
	fl.metas = append(fl.metas, m)
	fl.contents = append(fl.contents, content)
}

func (fl *FileList) Len() int {
	return len(fl.metas)
}

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

func (fl *FileList) Swap(x, y int) {
	length := fl.Len()
	if x >= 0 && x <= length && y >= 0 && y <= length && x != y {
		fl.metas[x], fl.metas[y] = fl.metas[y], fl.metas[x]
		fl.contents[x], fl.contents[y] = fl.contents[y], fl.contents[x]
	}
}
