package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-deb-deduplicate/archive"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("unar arFile destFolder")
		return
	}
	err := archive.UnAr(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}
