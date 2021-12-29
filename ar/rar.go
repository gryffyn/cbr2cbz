/*
https://github.com/mholt/archiver

MIT License

Copyright (c) 2016 Matthew Holt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */


package ar

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nwaples/rardecode"
)

type Rar struct {
	rr *rardecode.Reader     // underlying stream reader
}

func (r *Rar) Open(in io.Reader, size int64) error {
	if r.rr != nil {
		return fmt.Errorf("rar archive is already open for reading")
	}
	var err error
	r.rr, err = rardecode.NewReader(in, "")
	return err
}

func (r *Rar) Read() (File, error) {
	if r.rr == nil {
		return File{}, fmt.Errorf("rar archive is not open")
	}

	hdr, err := r.rr.Next()
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	file := File{
		FileInfo:   rarFileInfo{hdr},
		Header:     hdr,
		ReadCloser: ReadFakeCloser{r.rr},
	}

	return file, nil
}

// Close closes the rar archive(s) opened by Create and Open.
func (r *Rar) Close() error {
	var err error
	if r.rr != nil {
		r.rr = nil
	}
	return err
}

// Walk calls walkFn for each visited item in archive.
func (r *Rar) Walk(archive string, walkFn func(f File) error) error {
	file, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("opening archive file: %v", err)
	}
	defer file.Close()

	err = r.Open(file, 0)
	if err != nil {
		return fmt.Errorf("opening archive: %v", err)
	}
	defer r.Close()

	for {
		f, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("opening next file: %v", err)
		}
		err = walkFn(f)
		if err != nil {
			if err == ErrStopWalk {
				break
			}
			return fmt.Errorf("walking %s: %v", f.Name(), err)
		}
	}

	return nil
}

// CheckRAR checks if the file is a RAR, ZIP, or unidentified file.
// returns file type and error.
func (r *Rar) CheckRAR(archive string) (string, error) {
	RAR_HEADER := []byte{0x52, 0x61, 0x72, 0x21}
	ZIP_HEADER := []byte{0x50, 0x4b, 0x03, 0x04}

	file, err := os.Open(archive)
	if err != nil {
		return "", fmt.Errorf("opening archive file: %v", err)
	}

	buf := make([]byte, 4)
	_, err = file.Read(buf)
	if err != nil { return "", err }
	file.Close()
	if bytes.Equal(buf, ZIP_HEADER) {
		return "zip", errors.New("archive is already CBZ (extension incorrect?)")
	} else if !bytes.Equal(buf, RAR_HEADER) {
		return "und", errors.New("input is not a RAR archive")
	}
	return "rar", nil
}

type rarFileInfo struct {
	fh *rardecode.FileHeader
}

func (rfi rarFileInfo) Name() string       { return rfi.fh.Name }
func (rfi rarFileInfo) Size() int64        { return rfi.fh.UnPackedSize }
func (rfi rarFileInfo) Mode() os.FileMode  { return rfi.fh.Mode() }
func (rfi rarFileInfo) ModTime() time.Time { return rfi.fh.ModificationTime }
func (rfi rarFileInfo) IsDir() bool        { return rfi.fh.IsDir }
func (rfi rarFileInfo) Sys() interface{}   { return nil }