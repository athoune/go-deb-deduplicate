package chunker

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	chunkers "github.com/PlakarKorp/go-cdc-chunkers"
	_ "github.com/PlakarKorp/go-cdc-chunkers/chunkers/fastcdc"
	"github.com/klauspost/compress/zstd"
)

type Chunker struct {
	dir string
}

func New(dir string) *Chunker {
	return &Chunker{dir}
}

func (c Chunker) Chunk(rd io.Reader) error {
	chunker, err := chunkers.NewChunker("fastcdc", rd,
		&chunkers.ChunkerOpts{})
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(c.dir)
	if err != nil {
		return err
	}
	offset := 0
	for {
		chunk, err := chunker.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if len(chunk) == 0 {
			panic("empty chunk.")
		}
		hasher := sha256.New()
		hasher.Write(chunk)
		h := hex.EncodeToString(hasher.Sum(nil))
		_, err = root.Stat(h)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				f, err := root.OpenFile(h,
					os.O_WRONLY|os.O_CREATE, 0640)
				if err != nil {
					return err
				}
				defer f.Close()
				zw, err := zstd.NewWriter(f)
				if err != nil {
					return err
				}
				defer zw.Close()
				_, err = zw.Write(chunk)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		chunkLen := len(chunk)
		fmt.Println(offset, chunkLen)

		offset += chunkLen
	}
	return nil
}
