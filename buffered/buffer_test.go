package buffered

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummyList(b byte, size int) []byte {
	bb := make([]byte, size)
	for i := 0; i < size; i++ {
		bb[i] = b
	}
	return bb
}

type writeCloser struct {
	store *bytes.Buffer
}

func (w *writeCloser) Write(data []byte) (int, error) {
	return w.store.Write(data)
}

func (w *writeCloser) Close() error {
	return nil
}

type metaTest struct {
	id     byte
	whence int
	size   int
}

func (m *metaTest) SetPosition(whence int) error {
	m.whence = whence
	return nil
}

func TestBuffered(t *testing.T) {
	store := &writeCloser{&bytes.Buffer{}}
	docs := New(store)
	a := &metaTest{id: 'a'}
	_, err := docs.WriteDocument(a, dummyList('a', 128*1024))
	assert.NoError(t, err)
	b := &metaTest{id: 'b'}
	_, err = docs.WriteDocument(b, dummyList('b', 256*1024))
	assert.NoError(t, err)
	c := &metaTest{id: 'c'}
	_, err = docs.WriteDocument(c, dummyList('b', 92*1024))
	assert.NoError(t, err)
	err = docs.Close()
	assert.NoError(t, err)

	assert.Equal(t, (128+256+92)*1024, store.store.Len())
	assert.Equal(t, 256*1024, a.whence, "a")
	assert.Equal(t, 0, b.whence, "b")
	assert.Equal(t, (256+128)*1024, c.whence, "c")
}
