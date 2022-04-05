package tsv

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/athoune/ip2asn-go/ip"
)

type Line struct {
	RangeStart    net.IP
	RangeEnd      net.IP
	network       net.IPNet
	ASNumber      int
	CountryCode   string
	ASDescription string
}

func (l Line) Network() net.IPNet {
	return l.network
}

type Source struct {
	scanner *bufio.Scanner
}

func New(r io.Reader) *Source {
	return &Source{
		scanner: bufio.NewScanner(r),
	}
}

func (s *Source) Next() bool {
	return s.scanner.Scan()
}

func (s *Source) Values() (*Line, error) {
	line := strings.Split(s.scanner.Text(), "\t")
	var err error
	v := &Line{
		RangeStart:    net.ParseIP(line[0]),
		RangeEnd:      net.ParseIP(line[1]),
		CountryCode:   line[3],
		ASDescription: line[4],
	}
	v.ASNumber, err = strconv.Atoi(line[2])
	v.network = ip.Net(v.RangeStart, v.RangeEnd)
	return v, err
}
