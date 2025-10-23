package deduplicate

import (
	"github.com/athoune/go-deb-deduplicate/bytes"
	"github.com/athoune/go-deb-deduplicate/deb"
	"github.com/athoune/go-deb-deduplicate/warehouse"
	bolt "go.etcd.io/bbolt"
)

type Transaction struct {
	debManager  *deb.DebManager
	tx          *warehouse.Transaction
	bucket_meta *bolt.Bucket
}

func (d *Deduplicate) Transaction() (*Transaction, error) {
	var err error
	t := &Transaction{}
	t.tx, err = d.warehouse.Transaction()
	if err != nil {
		return nil, err
	}
	meta_b, err := t.tx.Tx.CreateBucketIfNotExists([]byte("meta"))
	if err != nil {
		return nil, err
	}
	t.debManager, err = deb.New(t.tx, meta_b)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transaction) AddPackage(path string) error {
	return t.debManager.AddPath(path)
}

func (t *Transaction) Close() error {
	return t.tx.Close()
}

func (t *Transaction) MissingChunks(old, fresh string) ([]byte, error) {
	o, err := t.debManager.GetDeb(old)
	if err != nil {
		return nil, err
	}
	f, err := t.debManager.GetDeb(fresh)
	if err != nil {
		return nil, err
	}
	return bytes.WhatsUp(o.Chunks(), f.Chunks(), 32), nil
}
