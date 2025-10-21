package deduplicate

import (
	"github.com/athoune/go-deb-deduplicate/deb"
	"github.com/athoune/go-deb-deduplicate/warehouse"
)

type Deduplicate struct {
	warehouse *warehouse.Warehouse
}

type Transaction struct {
	debReader *deb.DebReader
	tx        *warehouse.Transaction
}

func New(path string) (*Deduplicate, error) {
	var err error
	d := &Deduplicate{}
	d.warehouse, err = warehouse.New(path)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Deduplicate) Close() error {
	return d.warehouse.Close()
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
	t.debReader, err = deb.New(t.tx, meta_b)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transaction) Add(path string) error {
	return t.debReader.AddPath(path)
}

func (t *Transaction) Close() error {
	return t.tx.Close()
}
