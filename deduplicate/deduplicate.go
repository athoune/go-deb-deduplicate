package deduplicate

import (
	"github.com/athoune/go-deb-deduplicate/warehouse"
)

type Deduplicate struct {
	warehouse *warehouse.Warehouse
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
