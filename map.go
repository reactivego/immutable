package immutable

import (
	"fmt"
	"strings"

	"github.com/reactivego/immutable/byteorder"
)

type MapError string

func (e MapError) Error() string {
	return string(e)
}

const InvalidKeyType = MapError("Invalid Key Type")

type Hasher interface {
	Sum32() uint32
}

func prefix(key any) uint32 {
	if k, ok := key.(string); ok {
		return byteorder.LittleEndian.Uint32([]byte(k))
	} else if h, ok := key.(Hasher); ok {
		return h.Sum32()
	} else {
		panic(InvalidKeyType)
	}
}

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

func (n *Map) Lookup(key any) (any, bool) {
	return n.get(prefix(key), 0, key)
}

func (n *Map) Get(key any) any {
	v, _ := n.get(prefix(key), 0, key)
	return v
}

func (n Map) Put(key, value any) *Map {
	n.amt = n.put(prefix(key), 0, key, value)
	return &n
}

func (n Map) Delete(key any) *Map {
	n.amt = n.delete(prefix(key), 0, key)
	return &n
}

func (n *Map) Range(f func(any, any) bool) {
	n.foreach(f)
}

func (n Map) String() string {
	var b strings.Builder
	b.WriteByte('{')
	sep := byte(0)
	n.amt.foreach(func(k, v any) bool {
		if sep == 0 {
			sep = ','
		} else {
			b.WriteByte(sep)
		}
		_, err := fmt.Fprintf(&b, "%#v: %#v", k, v)
		return err == nil
	})
	b.WriteByte('}')
	return b.String()
}
