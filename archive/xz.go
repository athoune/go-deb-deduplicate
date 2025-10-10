package archive

import "io"

// See https://en.wikipedia.org/wiki/XZ_Utils

type XZindex struct {
	Header []byte
	Footer []byte
}

func XZReadIndex(file io.ReadSeeker) (*XZindex, error) {
	x := &XZindex{
		Header: make([]byte, 6+2+4),
		Footer: make([]byte, 4+4+2+2),
	}
	_, err := file.Read(x.Header)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(int64(-len(x.Footer)), io.SeekEnd)
	if err != nil {
		return nil, err
	}
	_, err = file.Read(x.Footer)
	if err != nil {
		return nil, err
	}

	return x, nil
}

func (x XZindex) PatchArchive(file io.WriteSeeker) error {
	_, err := file.Write(x.Header)
	if err != nil {
		return err
	}
	file.Seek(int64(len(x.Footer)), io.SeekEnd)
	_, err = file.Write(x.Footer)
	return err
}
