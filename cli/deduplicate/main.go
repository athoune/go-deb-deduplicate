package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-deb-deduplicate/deduplicate"
)

func add(paths ...string) error {
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
		err = tx.AddPackage(path)
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

func similarity(old string, fresh string) error {
	dedup, err := deduplicate.New("test_deb")
	if err != nil {
		return err
	}
	tx, err := dedup.Transaction()
	if err != nil {
		return err
	}
	defer tx.Close()
	s, err := tx.Similarity(old, fresh)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s : %v %%\n", old, fresh, s.Chunks*100)

	return nil
}

func ratio(old, fresh string) error {
	dedup, err := deduplicate.New("test_deb")
	if err != nil {
		return err
	}
	tx, err := dedup.Transaction()
	if err != nil {
		return err
	}
	defer tx.Close()
	r, err := tx.PatchRatio(old, fresh)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s : %v %%\n", old, fresh, r*100)
	return nil

}

func main() {
	switch os.Args[1] {
	case "add":
		err := add(os.Args[2:]...)
		if err != nil {
			panic(err)
		}
	case "similarity":
		err := similarity(os.Args[2], os.Args[3])
		if err != nil {
			panic(err)
		}
	case "ratio":
		err := ratio(os.Args[2], os.Args[3])
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("unknown action: ", os.Args[1])
	}
}
