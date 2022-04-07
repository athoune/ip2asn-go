package tree

import (
	"fmt"
	"io"
	"net"
	"sort"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type Node struct {
	Sons  map[byte]*Node
	Leafs []*Leaf
}

type Node2 struct {
	Name  byte
	Sons  Nodes
	Leafs []*Leaf
}

type Nodes []*Node2

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (n *Node2) Son(a byte) *Node2 {
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

func (n *Node2) SonOrNew(a byte) *Node2 {
	node := n.Son(a)
	if node != nil {
		return node
	}
	node = NewNode2(a)
	n.Sons = append(n.Sons, node)
	sort.Sort(n.Sons)
	return node
}

type Trunk2 struct {
	*Node2
}

func NewNode2(name byte) *Node2 {
	return &Node2{
		Name:  name,
		Sons:  make(Nodes, 0),
		Leafs: make([]*Leaf, 0),
	}
}

func (t *Trunk2) Append(nm *net.IPNet, data interface{}) {
	ones, _ := nm.Mask.Size()
	node := t.Node2
	for i := 0; i < ones/8; i++ {
		node = node.SonOrNew(nm.IP[i])
	}
	node.Leafs = append(node.Leafs, &Leaf{
		Netmask: nm,
		Data:    data,
	})
}

func (t *Trunk2) Get(ip net.IP) (interface{}, bool) {
	ip = ip.To4()
	node := t.Node2
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

func (t *Trunk) Get(_ip net.IP) (bool, interface{}) {
	chrono := time.Now()
	_ip = _ip.To4()
	key := _ip.String()
	value, ok := t.cache.Get(key)
	if ok {
		r := value.(response)
		fmt.Println("cache get", time.Now().Sub(chrono))
		return r.ok, r.value
	}
	node := t.Node
	cpt := 0
	for i := 0; i < len(_ip); i++ {
		var ok bool
		fmt.Println(len(node.Sons), "sons")
		node, ok = node.Sons[_ip[i]]
		if !ok {
			t.cache.Add(key, response{false, nil})
			fmt.Println(cpt, "tests for failing with", _ip)
			fmt.Println("No son", time.Now().Sub(chrono))
			return false, nil
		}
		for _, leaf := range node.Leafs {
			cpt++
			if leaf.Netmask.Contains(_ip) {
				t.cache.Add(key, response{true, leaf.Data})
				fmt.Println(cpt, "tests for", _ip)
				fmt.Println("subnet match", time.Now().Sub(chrono))
				return true, leaf.Data
			}
		}
	}
	fmt.Println(cpt, "tests for failing with", _ip)
	t.cache.Add(key, response{false, nil})
	fmt.Println("out of tree", time.Now().Sub(chrono))
	return false, nil
}

func (t *Trunk) Dump(w io.Writer) {
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
