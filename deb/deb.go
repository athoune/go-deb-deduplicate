package deb

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/athoune/go-deb-deduplicate/chunker"
	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

const DEB_BUCKET_NAME = "deb"

type DebReader struct {
	dataStore chunker.ChunkPutter
	metaStore chunker.ChunkPutter
}

func New(dataStore chunker.ChunkPutter, metaStore chunker.ChunkPutter) (*DebReader, error) {
	var err error
	d := &DebReader{
		dataStore: dataStore,
		metaStore: metaStore,
	}
	//d.metaStore, err = meta_tx.CreateBucketIfNotExists([]byte(DEB_BUCKET_NAME))
	return d, err
}

func (dr *DebReader) AddPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return dr.Add(f, f.Name())
}

func (dr *DebReader) Add(r io.Reader, name string) error {
	d, err := ReadPackage(r, name, dr.dataStore)
	bin, err := d.ToBin()
	if err != nil {
		return err
	}
	err = dr.dataStore.Put([]byte(name), bin)
	if err != nil {
		return err
	}
	return dr.metaStore.Put([]byte(name), bin)
}

// ReadPackage a deb package
func ReadPackage(r io.Reader, name string, putter chunker.ChunkPutter) (*Deb, error) {
	d := &Deb{
		Files: make([]*File, 0),
	}
	arReader := ar.NewReader(r)
	// FIXME handle .zstd format
	for {
		arHeader, err := arReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch arHeader.Name {
		case "debian-binary":
			d.Binary_h = arHeader
			continue
		case "control.tar.xz":
			d.Control_h = arHeader
		case "data.tar.xz":
			d.Data_h = arHeader
		default:
			return nil, fmt.Errorf("unknown root file in a .deb package : %v", arHeader.Name)
		}

		xzReader, err := xz.NewReader(arReader)
		if err != nil {
			return nil, err
		}
		tReader := tar.NewReader(xzReader)
		for {
			th, err := tReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			f := &File{
				Header: th,
				path:   path.Join(arHeader.Name, th.Name),
			}
			f.Contents, err = chunker.Chunk(tReader, putter)
			if err != nil {
				return nil, err
			}
			d.Files = append(d.Files, f)
		}
	}
	return d, nil
}
