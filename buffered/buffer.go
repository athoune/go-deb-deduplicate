package buffered

import (
	"bytes"
	"io"
)

const MIN_BUFFER_SIZE = 256 * 1024

type DocumentMeta interface {
	SetPosition(whence, size int) error
}

type position struct {
	meta   DocumentMeta
	whence int
	size   int
}

type BufferedDocumentWriter struct {
	buffer    *bytes.Buffer
	writer    io.WriteCloser
	positions []*position
	poz       int
}

func New(writer io.WriteCloser) *BufferedDocumentWriter {
	b := &bytes.Buffer{}
	b.Grow(MIN_BUFFER_SIZE)
	return &BufferedDocumentWriter{
		buffer:    b,
		writer:    writer,
		positions: make([]*position, 0),
	}
}

func (b *BufferedDocumentWriter) WriteDocument(doc DocumentMeta, data []byte) (int, error) {
	poz := &position{
		meta: doc,
		size: len(data),
	}
	if len(data) >= MIN_BUFFER_SIZE {
		_, err := b.writer.Write(data)
		if err != nil {
			return 0, err
		}
		err = doc.SetPosition(b.poz, len(data))
		if err != nil {
			return 0, err
		}
		b.poz += len(data)
		return 0, nil
	}

	_, err := b.buffer.Write(data)
	if err != nil {
		return 0, err
	}
	b.positions = append(b.positions, poz)
	if b.buffer.Len() > MIN_BUFFER_SIZE {
		_, err = b.flush()
		if err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

func (b *BufferedDocumentWriter) flush() (int, error) {
	if len(b.positions) == 0 { // nothing to flush
		return 0, nil
	}
	n := 0
	_, err := b.writer.Write(b.buffer.Bytes())
	if err != nil {
		return 0, err
	}
	for _, pz := range b.positions {
		err = pz.meta.SetPosition(b.poz, pz.size)
		if err != nil {
			return 0, err
		}
		b.poz += pz.size
	}
	// Reset
	b.buffer.Reset()
	b.positions = make([]*position, 0)
	return n, err
}

func (b *BufferedDocumentWriter) Close() error {
	_, err := b.flush()
	if err != nil {
		return err
	}
	return b.writer.Close()
}
