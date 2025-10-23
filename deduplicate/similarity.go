package deduplicate

import (
	_bytes "github.com/athoune/go-deb-deduplicate/bytes"
)

type Similarity struct {
	Chunks float64
}

func (t *Transaction) Similarity(old, fresh string) (*Similarity, error) {
	a, err := t.debManager.PathToChunks(old)
	if err != nil {
		return nil, err
	}
	b, err := t.debManager.PathToChunks(fresh)
	if err != nil {
		return nil, err
	}
	u := _bytes.Union(a, b, 32)
	return &Similarity{
		Chunks: float64(len(u)) / float64(len(b)),
	}, nil
}
