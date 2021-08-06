package snake

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
)

type Ziplib struct {
	Buffer      *bytes.Buffer
	FS          *zip.Writer
	ZipFileName string
}

func Zip(zipfile string) *Ziplib {
	z := new(Ziplib)
	z.Buffer = new(bytes.Buffer)
	z.FS = zip.NewWriter(z.Buffer)
	z.ZipFileName = zipfile
	return z
}

func (z *Ziplib) Add(path string, body []byte) bool {
	if file, err := z.FS.Create(path); err == nil {
		_, err := file.Write(body)
		return err == nil
	}
	return true
}

func (z *Ziplib) Close() error {
	err := z.FS.Close()
	if err == nil {
		ioutil.WriteFile(z.ZipFileName, z.Buffer.Bytes(), 0644)
	}
	return err
}
