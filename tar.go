package snake

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
)

type Tarlib struct {
	Buffer   *bytes.Buffer
	FS       *tar.Writer
	Gzip     *gzip.Writer
	FileName string
}

func Tar(tarfile string) *Tarlib {
	t := new(Tarlib)
	t.Buffer = new(bytes.Buffer)
	t.Gzip = gzip.NewWriter(t.Buffer)
	t.FS = tar.NewWriter(t.Gzip)
	t.FileName = tarfile
	return t
}

func (t *Tarlib) Add(path string, body []byte, typeflag byte) bool {
	if !String(path).Find(".DS_Store", true) && !String(path).Find("__MACOSX", true) {
		header := &tar.Header{
			Typeflag: typeflag,
			Name:     path,
			Mode:     0644,
			Size:     int64(len(body)),
		}
		if err := t.FS.WriteHeader(header); err == nil {
			_, err := t.FS.Write(body)
			return err == nil
		}
	}
	return false
}

func (t *Tarlib) Close() error {
	t.FS.Close()
	t.Gzip.Close()
	_, err := FS(t.FileName).ByteWriter(t.Buffer.Bytes())
	return err
}
