package deb

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	f, err := os.Open("../fixtures/golang_1.24~2_arm64.deb")
	assert.NoError(t, err)
	d, err := ReadPackage(f, "golang_1.24~2_arm64.deb", nil)
	assert.NoError(t, err)
	b, err := d.ToBin()
	assert.NoError(t, err)
	fmt.Println("Weight:", len(b))
	dd, err := FromBin(b)
	assert.NoError(t, err)
	assert.Equal(t, len(dd.Files), len(d.Files))
	for _, file := range dd.Files {
		fmt.Println(file.Header.Name, len(file.Contents))
	}
	//assert.False(t, true)
}
