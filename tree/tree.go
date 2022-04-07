package tree

import (
	"net"
	"sort"
)

type Node struct {
	Name  byte
	Sons  Nodes
	Leafs []*Leaf
}

type Nodes []*Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (n *Node) Son(a byte) *Node {
	if len(n.Sons) == 0 {
		return nil
	}
	r := sort.Search(len(n.Sons), func(i int) bool {
		return n.Sons[i].Name >= a
	})
	if r < len(n.Sons) {
		return n.Sons[r]
	}
	return nil
}

func (n *Node) SonOrNew(a byte) *Node {
	node := n.Son(a)
	if node != nil {
		return node
	}
	node = NewNode2(a)
	n.Sons = append(n.Sons, node)
	sort.Sort(n.Sons)
	return node
}

type Trunk struct {
	*Node
}

func NewNode2(name byte) *Node {
	return &Node{
		Name:  name,
		Sons:  make(Nodes, 0),
		Leafs: make([]*Leaf, 0),
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
