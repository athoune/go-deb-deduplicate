package main

import (
	"log"
	"os"

	"github.com/athoune/go-deb-deduplicate/archive"
)

func main() {
	name := os.Args[1]
	in, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.OpenFile(name[:len(name)-3], os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatal(err)
	}
	err = archive.XZDecompress(in, out)
	if err != nil {
		log.Fatal(err)
	}
}
