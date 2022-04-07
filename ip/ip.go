package ip

import (
	"net"
)

func Mask(a, b net.IP) net.IPMask {
	var s int
	if a.To4() == nil {
		s = 16
		a = a.To16()
		b = b.To16()
	} else {
		a = a.To4()
		b = b.To4()
		s = 4
	}
	r := make([]byte, s)
	for i := 0; i < s; i++ {
		r[i] = (b[i] - a[i]) ^ 255
	}
	return net.IPMask(r)
}

func Net(a, b net.IP) net.IPNet {
	m := Mask(a, b)
	if len(m) == 4 {
		a = a.To4()
	} else {
		a = a.To16()
	}
	return net.IPNet{
		IP:   a,
		Mask: m,
	}
}

func To8(a net.IP) [8]byte {
	a = a.To4()
	r := [8]byte{}
	for i := 0; i < 4; i++ {
		r[i*2] = a[i] & 0xF0 / 0x10
		r[i*2+1] = a[i] & 0x0F
	}
	return r
}
