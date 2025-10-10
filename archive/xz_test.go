package archive

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXZReadIndex(t *testing.T) {
	f, err := os.Open("../fixtures/README.md.xz")
	assert.NoError(t, err)
	idx, err := XZReadIndex(f)
	assert.NoError(t, err)
	assert.NotEmpty(t, idx.Header)
	assert.NotEmpty(t, idx.Footer)
}
