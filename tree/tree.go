package tree

import (
	"fmt"
	"net"
)

type Node struct {
	Sons    map[byte]*Node
	Netmask *net.IPNet
	Data    interface{}
}

type Trunk struct {
	*Node
}

func NewNode() *Node {
	return &Node{
		Sons: make(map[byte]*Node),
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
	node.Netmask = nm
	node.Data = data
}

func (t *Trunk) Get(ip net.IP) (bool, interface{}) {
	ip = ip.To4()
	node := t.Node
	for i := 0; i < 4; i++ {
		var ok bool
		fmt.Println(len(node.Sons), "sons")
		node, ok = node.Sons[ip[i]]
		if !ok {
			return false, nil
		}
		if node.Netmask != nil {
			if node.Netmask.Contains(ip) {
				return true, node.Data
			}
		}
	}
	return false, nil
}
