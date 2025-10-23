package deb

import (
	"archive/tar"
	"fmt"
	"io"
	"path"

	"github.com/athoune/go-deb-deduplicate/chunker"
	"github.com/athoune/go-deb-deduplicate/store"
	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

// ReadPackage read a deb package, chunk it and store chunks with a chunker.ChunkPutter
func ReadPackage(r io.Reader, name string, putter store.Putter) (*Deb, error) {
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
