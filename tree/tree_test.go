package tree

import (
	"net"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func TestNode2(t *testing.T) {
	a := NewNode2(0)
	a.SonOrNew(10)
	aa := a.Son(10)
	assert.NotNil(t, aa)
	aa = a.Son(11)
	assert.Nil(t, aa)
}

func TestTree2(t *testing.T) {
	tree := Trunk2{
		NewNode2(0),
	}
	_, nm, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(t, err)
	tree.Append(nm, "Hello")
	spew.Dump(tree)
	_, ok := tree.Get(net.ParseIP("192.168.1.42"))
	assert.True(t, ok)
	_, ok = tree.Get(net.ParseIP("192.168.2.42"))
	assert.False(t, ok)
}

func BenchmarkContains(b *testing.B) {
	_, nm, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(b, err)
	a := net.ParseIP("192.168.1.42")
	for i := 0; i < b.N; i++ {
		nm.Contains(a)
	}
}

func BenchmarkTree(b *testing.B) {
	trunk := NewNode()
	node := trunk
	_, nm, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(b, err)
	ones, _ := nm.Mask.Size()
	for i := 0; i < ones/8; i++ {
		n := NewNode()
		node.Sons[nm.IP[i]] = n
		node = n
	}
	node.Leafs = append(node.Leafs, &Leaf{
		Netmask: nm,
	})
	a := net.ParseIP("192.168.1.42").To4()
	for i := 0; i < b.N; i++ {
		node = trunk
		cpt := 0
		for j := 0; j < 4; j++ {
			n, ok := node.Sons[a[j]]
			if !ok {
				break
			}
			if len(n.Leafs) > 0 {
				break
			}
			cpt++
			node = n
		}
		assert.Equal(b, 2, cpt)
	}

}
