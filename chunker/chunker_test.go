package chunker

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	r := rand.New(rand.NewSource(99))
	noise := make([]byte, 100*1000)
	_, err := r.Read(noise)
	assert.NoError(t, err)

	buffer := bytes.NewBuffer(noise)
	ids, err := Chunk(buffer, nil)
	assert.NoError(t, err)
	assert.True(t, len(ids) > 0)
}
