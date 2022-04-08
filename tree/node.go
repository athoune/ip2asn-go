package tree

import "sort"

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
	node = NewNode(a)
	n.Sons = append(n.Sons, node)
	sort.Sort(n.Sons)
	return node
}

func NewNode(name byte) *Node {
	return &Node{
		Name:  name,
		Sons:  make(Nodes, 0),
		Leafs: make([]*Leaf, 0),
	}
}
