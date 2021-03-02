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
	if k, ok := key.([]byte); ok {
		return n.get(byteorder.LittleEndian.Uint32(k), 0, k)
	} else if k, ok := key.(string); ok {
		return n.get(byteorder.LittleEndian.Uint32([]byte(k)), 0, []byte(k))
	} else {
		panic(InvalidKeyType)
	}
}

func (n *Map) Get(key any) any {
	if k, ok := key.([]byte); ok {
		v, _ := n.get(byteorder.LittleEndian.Uint32(k), 0, k)
		return v
	} else if k, ok := key.(string); ok {
		v, _ := n.get(byteorder.LittleEndian.Uint32([]byte(k)), 0, []byte(k))
		return v
	} else {
		panic(InvalidKeyType)
	}
}

func (n *Map) Range(f func([]byte, any) bool) {
	n.foreach(f)
}

func (n Map) Put(key, value any) *Map {
	if k, ok := key.([]byte); ok {
		n.amt = n.put(byteorder.LittleEndian.Uint32(k), 0, k, value)
	} else if k, ok := key.(string); ok {
		n.amt = n.put(byteorder.LittleEndian.Uint32([]byte(k)), 0, []byte(k), value)
	} else {
		panic(InvalidKeyType)
	}
	return &n
}

func (n Map) Delete(key any) *Map {
	if k, ok := key.([]byte); ok {
		n.amt = n.delete(byteorder.LittleEndian.Uint32(k), 0, k)
	} else if k, ok := key.(string); ok {
		n.amt = n.delete(byteorder.LittleEndian.Uint32([]byte(k)), 0, []byte(k))
	} else {
		panic(InvalidKeyType)
	}
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
