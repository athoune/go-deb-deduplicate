package archive

import "github.com/blakesmith/ar"

type ArIndex struct {
	Headers []*ar.Header
}
