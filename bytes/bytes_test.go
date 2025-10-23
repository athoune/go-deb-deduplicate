package bytes

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhatsUp(t *testing.T) {
	old := []byte{0, 1, 0, 2, 0, 3}
	fresh := []byte{0, 2, 0, 4, 1, 0}
	diff := WhatsUp(old, fresh, 2)
	assert.Equal(t, []byte{0, 4, 1, 0}, diff)
}

func TestSort(t *testing.T) {
	fmt.Println("ab" > "ba")
	fmt.Println(string([]byte{0, 2}))
	a := ById{Id("\x00\x02"), Id("\x00\x04"), Id("\x00\x01")}
	sort.Sort(a)
	assert.Equal(t, ById{Id("\x00\x01"), Id("\x00\x02"), Id("\x00\x04")}, a)
}
