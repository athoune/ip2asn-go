package ip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	for _, test := range []struct {
		a net.IP
		b net.IP
		n string
	}{
		{
			a: net.ParseIP("192.168.1.0"),
			b: net.ParseIP("192.168.1.255"),
			n: "192.168.1.0/24",
		},
		{
			a: net.ParseIP("1.0.0.0"),
			b: net.ParseIP("1.0.0.255"),
			n: "1.0.0.0/24",
		},
	} {
		n := Net(test.a, test.b)
		assert.Equal(t, test.n, n.String())
		assert.Len(t, n.Mask, 4)
		_, nn, err := net.ParseCIDR(test.n)
		assert.NoError(t, err)
		assert.Equal(t, nn, &n)
	}
}
