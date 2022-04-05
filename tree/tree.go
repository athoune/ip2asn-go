package tree

import (
	"fmt"
	"net"
)

type Node struct {
	Sons  map[byte]*Node
	Leafs []*Leaf
}

type Leaf struct {
	Netmask *net.IPNet
	Data    interface{}
}

type Trunk struct {
	*Node
}

func NewNode() *Node {
	return &Node{
		Sons:  make(map[byte]*Node),
		Leafs: make([]*Leaf, 0),
	}
}

func New() *Trunk {
	return &Trunk{
		&Node{
			Sons: make(map[byte]*Node),
		},
	}
}

func (t *Trunk) Append(nm *net.IPNet, data interface{}) {
	ones, _ := nm.Mask.Size()
	node := t.Node
	for i := 0; i < ones/8; i++ {
		n, ok := node.Sons[nm.IP[i]]
		if !ok {
			n = NewNode()
			node.Sons[nm.IP[i]] = n
		}
		node = n
	}
	node.Leafs = append(node.Leafs, &Leaf{
		Netmask: nm,
		Data:    data,
	})
}

func (t *Trunk) Get(ip net.IP) (bool, interface{}) {
	ip = ip.To4()
	node := t.Node
	cpt := 0
	for i := 0; i < 4; i++ {
		var ok bool
		fmt.Println(len(node.Sons), "sons")
		node, ok = node.Sons[ip[i]]
		if !ok {
			fmt.Println(cpt, "tests for failing with", ip)
			return false, nil
		}
		for _, leaf := range node.Leafs {
			cpt++
			if leaf.Netmask.Contains(ip) {
				fmt.Println(cpt, "tests for", ip)
				return true, leaf.Data
			}
		}
	}
	fmt.Println(cpt, "tests for failing with", ip)
	return false, nil
}
