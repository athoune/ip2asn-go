package tsv

import (
	"compress/gzip"
	"net"
	"os"
	"testing"

	_tree "github.com/athoune/iptree/tree"
	"github.com/stretchr/testify/assert"
)

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
	tree := _tree.NewTrunk(2)
	err = FeedTrunk(tree, New(r))
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
	tree, err := _tree.NewCachedTrunk(256, 2)
	assert.NoError(b, err)
	err = FeedTrunk(tree, New(r))
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
