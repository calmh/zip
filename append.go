// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Append files to an existing Zip file.

package zip

import (
	"os"
)

type AppendWriter struct {
	Writer
}

func NewAppendWriter(rw *os.File) (*AppendWriter, error) {
	fi, err := rw.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()

	if size > 0 {
		zr, err := NewReader(rw, size)
		if err != nil {
			return nil, err
		}

		_, err = rw.Seek(size, os.SEEK_SET)
		if err != nil {
			return nil, err
		}

		zw := NewWriter(rw)
		for _, f := range zr.File {
			zw.dir = append(zw.dir, &header{
				&f.FileHeader, uint64(f.headerOffset),
			})
		}
		zw.cw.count = size

		return &AppendWriter{*zw}, nil
	} else {
		return &AppendWriter{*NewWriter(rw)}, nil
	}
}

func (w *AppendWriter) Remove(name string) {
	var dir []*header
	for _, f := range w.dir {
		if f.Name != name {
			dir = append(dir, f)
		}
	}
	w.dir = dir
}
