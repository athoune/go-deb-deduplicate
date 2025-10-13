package archive

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

func UnXZ(xzFile string) error {
	f, err := os.Open(xzFile)
	if err != nil {
		return err
	}
	idx, err := XZReadIndex(f)
	if err != nil {
		return fmt.Errorf("ZX Read Index error: %v", err)
	}

	dst, err := os.OpenFile(fmt.Sprintf("%s.index", xzFile),
		os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer dst.Close()
	enc := gob.NewEncoder(dst)
	err = enc.Encode(idx)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	dstArchive, err := os.OpenFile(xzFile[:len(xzFile)-3],
		os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	reader, err := xz.NewReader(f)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstArchive, reader)
	return err
}
