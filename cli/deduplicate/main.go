package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-deb-deduplicate/deduplicate"
)

func Read(paths ...string) error {

	dedup, err := deduplicate.New("test_deb")
	if err != nil {
		return err
	}
	tx, err := dedup.Transaction()
	if err != nil {
		return err
	}

	for _, path := range paths {
		fmt.Println(path)
		err = tx.Add(path)
		if err != nil {
			return err
		}
	}

	err = tx.Close()
	if err != nil {
		return err
	}
	return dedup.Close()
}

func main() {
	err := Read(os.Args[1:]...)
	if err != nil {
		panic(err)
	}
}
