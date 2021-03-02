package immutable

import (
	"fmt"
	"strings"

	"github.com/reactivego/immutable/byteorder"
)

type Map struct {
	*amt
}

func NewMap() *Map {
	return &Map{&amt{}}
}

func (n *Map) Len() int {
	return n.len()
}

func (n *Map) Depth() int {
	return n.depth()
}

func (n *Map) Lookup(key []byte) (any, bool) {
	return n.get(byteorder.LittleEndian.Uint32(key), 0, key)
}

func (n *Map) Get(key []byte) any {
	v, _ := n.get(byteorder.LittleEndian.Uint32(key), 0, key)
	return v
}

func (n *Map) Range(f func([]byte, any) bool) {
	n.foreach(f)
}

func (n Map) Put(key []byte, value any) *Map {
	n.amt = n.put(byteorder.LittleEndian.Uint32(key), 0, key, value)
	return &n
}

func (n Map) Remove(key []byte) *Map {
	n.amt = n.remove(byteorder.LittleEndian.Uint32(key), 0, key)
	return &n
}

func (n Map) String() string {
	var b strings.Builder
	b.WriteByte('{')
	sep := byte(0)
	n.amt.foreach(func(k []byte, v any) bool {
		if sep == 0 {
			sep = ','
		} else {
			b.WriteByte(sep)
		}
		_, err := fmt.Fprintf(&b, "%q: %#v", k, v)
		return err == nil
	})
	b.WriteByte('}')
	return b.String()
}
