package tree

import (
	"io"
	"net"

	"github.com/athoune/ip2asn-go/tsv"
)

type Trunk struct {
	*Node
	size int
}

func NewTrunk() *Trunk {
	return &Trunk{
		NewNode(0),
		0,
	}
}

func (t *Trunk) Append(nm *net.IPNet, data interface{}) {
	ones, _ := nm.Mask.Size()
	node := t.Node
	for i := 0; i < ones/8; i++ {
		node = node.SonOrNew(nm.IP[i])
	}
	node.Leafs = append(node.Leafs, &Leaf{
		Netmask: nm,
		Data:    data,
	})
	t.size++
}

func (t *Trunk) Size() int {
	return t.size
}

func (t *Trunk) FeedWithTSV(r io.Reader) error {
	src := tsv.New(r)
	for src.Next() {
		line, err := src.Values()
		if err != nil {
			return err
		}
		n := line.Network()
		//fmt.Println(n, line.ASNumber, line.CountryCode, line.ASDescription)
		if line.ASNumber != 0 {
			if n.IP.To4() != nil {
				t.Append(&n, line)
			}
		}
	}
	return nil
}

func (t *Trunk) Get(ip net.IP) (interface{}, bool) {
	ip = ip.To4()
	node := t.Node
	for i := 0; i < 4; i++ {
		n := node.Son(ip[i])
		if n == nil {
			return nil, false
		}
		for _, leaf := range n.Leafs {
			if leaf.Netmask.Contains(ip) {
				return leaf.Data, true
			}
		}
		node = n
	}
	return nil, false
}

type Leaf struct {
	Netmask *net.IPNet
	Data    interface{}
}

type response struct {
	ok    bool
	value interface{}
}

/*
func (t *Trunk2) Dump(w io.Writer) {
	dump(w, 0, t.Node)
}

func dump(w io.Writer, tabs int, node *Node) {
	for key, son := range node.Sons {
		for i := 0; i < tabs; i++ {
			fmt.Fprint(w, "-")
		}
		fmt.Fprintf(w, "%x", key)
		for _, leaf := range son.Leafs {
			fmt.Fprintf(w, " %v", leaf.Netmask)
		}
		fmt.Fprint(w, "\n")
		dump(w, tabs+1, son)
	}
}

*/
