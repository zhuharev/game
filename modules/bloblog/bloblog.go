package bloblog

import (
	"github.com/zhuharev/bloblog"
)

var (
	bl *bloblog.BlobLog
)

// NewContext inti bloblog database
func NewContext() error {
	var err error
	bl, err = bloblog.Open("_state/bloblog", 8*2000000)
	return err
}

// Save save file to bloblog
func Save(data []byte) (int64, error) {
	return bl.Insert(data)
}

// Get returns file from bloblog
func Get(id int64) ([]byte, error) {
	return bl.Get(id)
}
