package main

import (
	"log"
	"os"

	"github.com/athoune/go-deb-diff/archive"
)

func main() {
	err := archive.UnXZ(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}
