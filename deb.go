package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	chunker_ "github.com/athoune/go-deb-deduplicate/chunker"
	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

func Read(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	chunker := chunker_.New("chunks")
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
			th, err := tReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			fmt.Printf("\t\t%v %d\n", th.Name, th.Size)
			err = chunker.Chunk(tReader)
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
