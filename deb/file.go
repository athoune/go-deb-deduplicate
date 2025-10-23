package deb

import (
	"archive/tar"
	"bytes"

	"github.com/blakesmith/ar"
	"github.com/hashicorp/go-msgpack/codec"
)

type Deb struct {
	Binary_h  *ar.Header
	Control_h *ar.Header
	Data_h    *ar.Header
	Files     []*File
}

type File struct {
	Header   *tar.Header
	path     string
	Contents []byte // sha256 of their chunks
}

func FromBin(data []byte) (*Deb, error) {
	h := &codec.MsgpackHandle{}
	decoder := codec.NewDecoderBytes(data, h)
	d := &Deb{
		Files: make([]*File, 0),
	}
	err := decoder.Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// ToBin serialize Deb object
func (d *Deb) ToBin() ([]byte, error) {
	w := &bytes.Buffer{}
	h := &codec.MsgpackHandle{}
	enc := codec.NewEncoder(w, h)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (d *Deb) Chunks() []byte {
	l := 0
	for _, f := range d.Files {
		l += len(f.Contents)
	}
	r := make([]byte, l)
	i := 0
	for _, f := range d.Files {
		ll := len(f.Contents)
		copy(r[i:i+ll], f.Contents)
		i += ll
	}
	return r
}
