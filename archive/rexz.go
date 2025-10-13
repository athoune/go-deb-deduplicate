package archive

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

func ReXZ(file string) error {
	idxFile, err := os.Open(fmt.Sprintf("%s.xz.index", file))
	if err != nil {
		return err
	}
	defer idxFile.Close()
	decoder := gob.NewDecoder(idxFile)
	idx := &XZindex{
		Header: make([]byte, 0),
		Footer: make([]byte, 0),
	}
	err = decoder.Decode(idx)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	xzFile, err := os.OpenFile(fmt.Sprintf("%s.xz", file), os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return err
	}
	defer xzFile.Close()
	writer, err := xz.NewWriter(xzFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, srcFile)
	if err != nil {
		return err
	}
	err = idx.PatchArchive(xzFile)
	if err != nil {
		return err
	}
	return nil
}
