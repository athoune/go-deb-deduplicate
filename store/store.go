package store

type Putter interface {
	Put([]byte, []byte) error
}

type Getter interface {
	Get(key []byte) []byte
}

type GetterPutter interface {
	Getter
	Putter
}
