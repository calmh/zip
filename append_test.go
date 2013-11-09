// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tests that involve both reading and writing.

package zip

import (
	"fmt"
	"os"
	"testing"
)

func TestAppend(t *testing.T) {
	testFile := "appendtest.zip"

	f, _ := os.Create(testFile)

	w, err := NewAppendWriter(f)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	nFiles := 10

	for i := 0; i < nFiles; i++ {
		fw, err := w.CreateHeader(&FileHeader{
			Name:   fmt.Sprintf("%d.dat", i),
			Method: Store, // avoid Issue 6136 and Issue 6138
		})
		fmt.Fprintf(fw, "content of file number %d here", i)
		if err != nil {
			t.Fatalf("creating file %d: %v", i, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}

	fi, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	size := fi.Size()

	f.Close()

	f, _ = os.OpenFile(testFile, os.O_RDWR, 0666)

	w, err = NewAppendWriter(f)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	for i := nFiles; i < 2*nFiles; i++ {
		fw, err := w.CreateHeader(&FileHeader{
			Name:   fmt.Sprintf("%d.dat", i),
			Method: Store, // avoid Issue 6136 and Issue 6138
		})
		fmt.Fprintf(fw, "content of file number %d here", i)
		if err != nil {
			t.Fatalf("creating file %d: %v", i, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}

	fi, err = f.Stat()
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	size = fi.Size()

	f.Close()

	f, _ = os.Open(testFile)

	zr, err := NewReader(f, size)

	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	if got := len(zr.File); got != 2*nFiles {
		t.Fatalf("File contains %d files, want %d", got, 2*nFiles)
	}
	for i := 0; i < 2*nFiles; i++ {
		want := fmt.Sprintf("%d.dat", i)
		if zr.File[i].Name != want {
			t.Fatalf("File(%d) = %q, want %q", i, zr.File[i].Name, want)
		}
	}

	os.Remove(testFile)
}
