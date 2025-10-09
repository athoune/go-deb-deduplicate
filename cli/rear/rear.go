package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-deb-diff/archive"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("rear dest srcFolder")
		return
	}
	err := archive.ReAr(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}
