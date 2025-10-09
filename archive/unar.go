package archive

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/blakesmith/ar"
)

type ArReader interface {
	Next(io.Reader)
}

func UnAr(arFile string, destFolder string) error {
	source, err := os.Open(arFile)
	if err != nil {
		return err
	}
	err = os.Mkdir(destFolder, 0750)
	if err != nil {
		return err
	}
	idx := &ArIndex{
		Headers: make([]*ar.Header, 0),
	}
	reader := ar.NewReader(source)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		idx.Headers = append(idx.Headers, header)
		dst, err := os.OpenFile(fmt.Sprintf("%s/%s", destFolder, header.Name), os.O_CREATE|os.O_WRONLY, 0440)
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, reader)
		if err != nil {
			return err
		}
	}
	dst, err := os.OpenFile(fmt.Sprintf("%s/index.gob", destFolder), os.O_CREATE|os.O_WRONLY, 0440)
	if err != nil {
		return err
	}
	defer dst.Close()
	enc := gob.NewEncoder(dst)
	return enc.Encode(idx)
}
