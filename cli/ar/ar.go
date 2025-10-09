package main

import (
	"encoding/gob"
	"os"

	"github.com/blakesmith/ar"
)

func UnAr(source string, archive string) error {
	f_source, err := os.OpenFile(source, os.O_RDONLY, 0440)
	if err != nil {
		return err
	}
	f_archive, err := os.OpenFile(archive, os.O_WRONLY|os.O_CREATE, 0440)
	if err != nil {
		return err
	}
	reader := ar.NewReader(f_source)
	header, err := reader.Next()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(f_archive)
	err = enc.Encode(header)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := UnAr(os.Args[2], os.Args[1])
	if err != nil {
		panic(err)
	}
}
