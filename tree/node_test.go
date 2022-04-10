package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	a := NewNode(0, false)
	a.SonOrNew(10, false)
	aa := a.Son(10)
	assert.NotNil(t, aa)
	aa = a.Son(11)
	assert.Nil(t, aa)
}
