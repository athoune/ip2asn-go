package tree

import (
	"compress/gzip"
	"net"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestTree2(t *testing.T) {
	tree := NewTrunk(2)
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
	f, err := os.Open("../ip2asn-v4.tsv.gz")
	assert.NoError(b, err)
	r, err := gzip.NewReader(f)
	assert.NoError(b, err)
	tree := NewTrunk(2)
	err = tree.FeedWithTSV(r)
	assert.NoError(b, err)
	freeS, err := net.LookupHost("free.fr")
	assert.NoError(b, err)
	var free net.IP
	for _, f := range freeS {
		i := net.ParseIP(f)
		if i.To4() != nil {
			free = i
			break
		}
	}
	assert.NotNil(b, free)
	for i := 0; i < b.N; i++ {
		_, ok := tree.Get(free)
		assert.True(b, ok)
	}
}

func BenchmarkCachedTree(b *testing.B) {
	f, err := os.Open("../ip2asn-v4.tsv.gz")
	assert.NoError(b, err)
	r, err := gzip.NewReader(f)
	assert.NoError(b, err)
	tree, err := NewCachedTrunk(256, 2)
	assert.NoError(b, err)
	err = tree.FeedWithTSV(r)
	assert.NoError(b, err)
	freeS, err := net.LookupHost("google.fr")
	assert.NoError(b, err)
	var free net.IP
	for _, f := range freeS {
		i := net.ParseIP(f)
		if i.To4() != nil {
			free = i
			break
		}
	}
	assert.NotNil(b, free)
	for i := 0; i < b.N; i++ {
		_, ok := tree.Get(free)
		assert.True(b, ok)
	}
}
