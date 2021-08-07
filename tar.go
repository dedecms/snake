package snake

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
)

type Tarlib struct {
	Buffer      *bytes.Buffer
	FS          *tar.Writer
	Gzip        *gzip.Writer
	ZipFileName string
}

func Tar(tarfile string) *Tarlib {
	t := new(Tarlib)
	t.Buffer = new(bytes.Buffer)
	t.Gzip = gzip.NewWriter(t.Buffer)
	t.FS = tar.NewWriter(t.Buffer)
	t.ZipFileName = tarfile
	return t
}

func (t *Tarlib) Add(path string, body []byte) bool {

	header := &tar.Header{
		Name: path,
		Mode: 0644,
		Size: int64(len(body)),
	}
	if err := t.FS.WriteHeader(header); err != nil {
		return false
	}
	if _, err := t.FS.Write(body); err == nil {
		return true
	}

	return false
}

func (z *Tarlib) Close() error {
	z.Gzip.Close()
	err := z.FS.Close()
	if err == nil {
		FS(z.ZipFileName).ByteWriter(z.Buffer.Bytes())
	}
	return err
}
