package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

func Read(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
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
			var buf bytes.Buffer
			io.Copy(&buf, tReader)
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
