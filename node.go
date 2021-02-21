package immutable

import	"github.com/reactivego/immutable/byteorder"

type Node struct {
	*amt
}

func New() *Node {
	return &Node{&amt{}}
}

func (n *Node) Len() int {
	return n.len()
}

func (n *Node) Depth() int {
	return n.depth()
}

func (n *Node) Lookup(key []byte) (any, bool) {
	return n.get(byteorder.LittleEndian.Uint32(key), 0, key)
}

func (n *Node) Get(key []byte) any {
	v, _ := n.get(byteorder.LittleEndian.Uint32(key), 0, key)
	return v
}

func (n *Node) Range(f func([]byte, any)) {
	n.enum(f)
}

func (n Node) Put(key []byte, value any) *Node {
	n.amt = n.put(byteorder.LittleEndian.Uint32(key), 0, key, value)
	return &n
}

func (n Node) Remove(key []byte) *Node {
	n.amt = n.remove(byteorder.LittleEndian.Uint32(key), 0, key)
	return &n
}
