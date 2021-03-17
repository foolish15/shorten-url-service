package hashid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenHash(t *testing.T) {
	st, err := Encrypt(1)
	assert.Nil(t, err)
	fmt.Println(st)
}
