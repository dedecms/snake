package snake

import (
	"archive/zip"
	"bytes"
	"io/fs"
)

type Ziplib struct {
	Buffer   *bytes.Buffer
	FS       *zip.Writer
	FileName string
}

func Zip(zipfile string) *Ziplib {
	z := new(Ziplib)
	z.Buffer = new(bytes.Buffer)
	z.FS = zip.NewWriter(z.Buffer)
	z.FileName = zipfile
	return z
}

func (z *Ziplib) Add(path string, stat fs.FileInfo, body []byte) bool {
	if !String(path).Find(".DS_Store", true) && !String(path).Find("__MACOSX", true) && !String(path).Find(".gitignore", true) && !String(path).Find(".index", true) {
		header, _ := zip.FileInfoHeader(stat)
		header.Name = path
		if stat.IsDir() {
			header.Name += "/"
			z.FS.CreateHeader(header)
		} else {
			if file, err := z.FS.Create(path); err == nil {
				_, err := file.Write(body)
				return err == nil
			}
		}
	}
	return false
}

func (z *Ziplib) Close() error {
	err := z.FS.Close()
	if err == nil {
		FS(z.FileName).ByteWriter(z.Buffer.Bytes())
	}
	return err
}
