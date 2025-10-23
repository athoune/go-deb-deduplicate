package deb

import (
	"bytes"
	"fmt"
	"io"
	"os"

	bytes_ "github.com/athoune/go-deb-deduplicate/bytes"
	"github.com/athoune/go-deb-deduplicate/store"
)

const DEB_BUCKET_NAME = "deb"

type DebManager struct {
	dataStore store.GetterPutter
	metaStore store.GetterPutter
}

func New(dataStore store.GetterPutter, metaStore store.GetterPutter) (*DebManager, error) {
	var err error
	d := &DebManager{
		dataStore: dataStore,
		metaStore: metaStore,
	}
	return d, err
}

// AddPath add a deb package path
func (dr *DebManager) AddPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return dr.Add(f, f.Name())
}

// Add read a de b package
func (dr *DebManager) Add(r io.Reader, name string) error {
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

func (dr *DebManager) GetDeb(name string) (*Deb, error) {
	return nil, nil
}

func (dr *DebManager) pathToDeb(path string) (*Deb, error) {
	meta_raw := dr.metaStore.Get([]byte(path))
	if meta_raw == nil {
		return nil, fmt.Errorf("unknown package name: %s", path)
	}
	meta, err := FromBin(meta_raw)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (dr *DebManager) PathToChunks(path string) ([]byte, error) {
	meta, err := dr.pathToDeb(path)
	if err != nil {
		return nil, err
	}
	buff := &bytes.Buffer{}
	for _, file := range meta.Files {
		_, err = buff.Write(file.Contents)
		if err != nil {
			return nil, err
		}
	}
	return buff.Bytes(), nil
}

func (dr *DebManager) Union(path_a string, path_b string) ([]byte, error) {
	a, err := dr.PathToChunks(path_a)
	if err != nil {
		return nil, err
	}
	b, err := dr.PathToChunks(path_b)
	if err != nil {
		return nil, err
	}
	return bytes_.Union(a, b, 32), nil
}
