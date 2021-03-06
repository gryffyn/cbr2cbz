package ar

import (
	"io"

	"github.com/klauspost/compress/zip"
)

type Zip struct {
	zr *zip.Writer
}

func (z *Zip) Open(in io.Writer) {
	z.zr = zip.NewWriter(in)
}

func (z *Zip) WriteFile(file File) error {
	f, err := z.zr.Create(file.Name())
	if err != nil { return err }
	buf, err := io.ReadAll(file.ReadCloser)
	if err != nil { return err }
	_, err = f.Write(buf)
	return err
}

func (z *Zip) Close() error {
	return z.zr.Close()
}