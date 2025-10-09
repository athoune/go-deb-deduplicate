package archive

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/blakesmith/ar"
)

func ReAr(dest string, srcFolder string) error {
	idxReader, err := os.Open(fmt.Sprintf("%s/index.gob", srcFolder))
	if err != nil {
		return err
	}
	defer idxReader.Close()
	dec := gob.NewDecoder(idxReader)
	idx := &ArIndex{
		Headers: make([]*ar.Header, 0),
	}
	err = dec.Decode(idx)
	if err != nil {
		return err
	}
	archiveWriter, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0440)
	if err != nil {
		return err
	}
	archive := ar.NewWriter(archiveWriter)
	err = archive.WriteGlobalHeader()
	if err != nil {
		return err
	}
	for _, f := range idx.Headers {
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", srcFolder, f.Name))
		if err != nil {
			return err
		}
		err = archive.WriteHeader(f)
		if err != nil {
			return err
		}
		_, err = archive.Write(content)
		if err != nil {
			return err
		}
	}
	return nil
}
