package main

import (
	"fmt"
	"os"

	"github.com/athoune/go-deb-deduplicate/archive"
)

func main() {
	name := os.Args[1]
	in, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	out, err := os.OpenFile(fmt.Sprintf("%s.xz", name), os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	err = archive.XZCompress(in, out)
	if err != nil {
		panic(err)
	}
}
