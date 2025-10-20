package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	chunker_ "github.com/athoune/go-deb-deduplicate/chunker"
	"github.com/athoune/go-deb-deduplicate/warehouse"
	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

func Read(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	store, err := warehouse.New("test_deb")
	if err != nil {
		return err
	}
	tx, err := store.Transaction()
	if err != nil {
		return err
	}
	reader := ar.NewReader(f)
	header, err := reader.Next()
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", header)

	for {
		h, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("\t%#v\n", h)
		xzReader, err := xz.NewReader(reader)
		if err != nil {
			return err
		}
		tReader := tar.NewReader(xzReader)
		for {
			_, err := tReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			err = chunker_.ChunkAndStore(tReader, tx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	err := Read(os.Args[1])
	if err != nil {
		panic(err)
	}

}
