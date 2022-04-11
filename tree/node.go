package tree

import "sort"

type Node struct {
	Name     byte
	Sons     Nodes
	Leafs    []*Leaf
	fullList bool
}

type Nodes []*Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (n *Node) Son(a byte) *Node {
	if n.fullList {
		return n.Sons[a]
	}
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

func (n *Node) SonOrNew(a byte, full bool) *Node {
	node := n.Son(a)
	if node != nil {
		return node
	}
	node = NewNode(a, full)
	if n.fullList {
		n.Sons[a] = node
	} else {
		n.Sons = append(n.Sons, node)
		sort.Sort(n.Sons)
	}
	return node
}

func NewNode(name byte, full bool) *Node {
	var size int
	if full {
		size = 256
	} else {
		size = 0
	}
	return &Node{
		Name:     name,
		Sons:     make(Nodes, size),
		Leafs:    make([]*Leaf, 0),
		fullList: full,
	}
}
