package tree

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {

	tree, err := New(1)
	assert.NoError(t, err)
	_, nm, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(t, err)
	tree.Append(nm, "Hello")
	ok, _ := tree.Get(net.ParseIP("192.168.1.42"))
	assert.True(t, ok)
	ok, _ = tree.Get(net.ParseIP("192.168.2.42"))
	assert.False(t, ok)
}
