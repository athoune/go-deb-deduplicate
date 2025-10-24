package chunker

import (
	"bytes"
	"crypto/sha256"
	"io"
	"log"

	chunkers "github.com/PlakarKorp/go-cdc-chunkers"
	_ "github.com/PlakarKorp/go-cdc-chunkers/chunkers/fastcdc"
	"github.com/athoune/go-deb-deduplicate/store"
)

const CHUNK_MIN_SIZE = 300
const CHUNK_NORMAL_SIZE = 500
const CHUNK_MAX_SIZE = 700

func Chunk(reader io.Reader, putter store.Putter) ([]byte, error) {
	chunker, err := chunkers.NewChunker("fastcdc", reader,
		&chunkers.ChunkerOpts{
			MinSize:    CHUNK_MIN_SIZE,
			NormalSize: CHUNK_NORMAL_SIZE,
			MaxSize:    CHUNK_MAX_SIZE,
		})
	if err != nil {
		return nil, err
	}
	ids := &bytes.Buffer{}

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
		h := hasher.Sum(nil)
		if putter != nil {
			err = putter.Put(h, chunk)
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err = ids.Write(h)
		if err != nil {
			log.Fatal(err)
		}
	}
	return ids.Bytes(), nil
}
