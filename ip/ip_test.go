package ip

import (
	"fmt"
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

func TestTo8(t *testing.T) {
	a := net.IPv4(192, 168, 1, 1).To4()
	fmt.Printf("%x %x %x %x\n", a[0], a[1], a[2], a[3])
	aa := To8(a)
	fmt.Printf("%x %x %x %x %x %x %x %x\n", aa[0], aa[1], aa[2], aa[3], aa[4], aa[5], aa[6], aa[7])
	assert.Equal(t, [8]byte{0xc, 0x0, 0xa, 0x8, 0x0, 0x1, 0x0, 0x1}, aa)
}
