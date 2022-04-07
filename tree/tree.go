package tree

import (
	"fmt"
	"net"
	"time"

	lru "github.com/hashicorp/golang-lru"
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
	cache *lru.Cache
}

func NewNode() *Node {
	return &Node{
		Sons:  make(map[byte]*Node),
		Leafs: make([]*Leaf, 0),
	}
}

func New(size int) (*Trunk, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &Trunk{
		&Node{
			Sons: make(map[byte]*Node),
		},
		cache,
	}, nil
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

type response struct {
	ok    bool
	value interface{}
}

func (t *Trunk) Get(ip net.IP) (bool, interface{}) {
	chrono := time.Now()
	ip = ip.To4()
	key := ip.String()
	value, ok := t.cache.Get(key)
	if ok {
		r := value.(response)
		fmt.Println("cache get", time.Now().Sub(chrono))
		return r.ok, r.value
	}
	node := t.Node
	cpt := 0
	for i := 0; i < 4; i++ {
		var ok bool
		fmt.Println(len(node.Sons), "sons")
		node, ok = node.Sons[ip[i]]
		if !ok {
			fmt.Println(cpt, "tests for failing with", ip)
			t.cache.Add(key, response{false, nil})
			return false, nil
		}
		for _, leaf := range node.Leafs {
			cpt++
			if leaf.Netmask.Contains(ip) {
				fmt.Println(cpt, "tests for", ip)
				t.cache.Add(key, response{true, leaf.Data})
				return true, leaf.Data
			}
		}
	}
	fmt.Println(cpt, "tests for failing with", ip)
	t.cache.Add(key, response{false, nil})
	return false, nil
}
