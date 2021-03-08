package immutable

import (
	"fmt"
	"hash/maphash"
	"strings"
)

type MapError string

func (e MapError) Error() string {
	return string(e)
}

const InvalidKeyType = MapError("Invalid Key Type")

type Hasher interface {
	Sum32() uint32
}

var seed = maphash.MakeSeed()

func hash(key any) uint32 {
	if k, ok := key.(string); ok {
		var h maphash.Hash
		h.SetSeed(seed)
		h.WriteString(k)
		return uint32(h.Sum64() & 0xFFFFFFFF)
	} else if h, ok := key.(Hasher); ok {
		return h.Sum32()
	} else {
		panic(InvalidKeyType)
	}
}

type Hamt struct {
	*amt
}

var Map = &Hamt{&amt{}}

func (a *Hamt) Len() int {
	return a.len()
}

func (a *Hamt) Depth() int {
	return a.depth()
}

func (a *Hamt) Size() int {
	return 8 + a.size()
}

func (a *Hamt) Lookup(key any) (any, bool) {
	return a.get(hash(key), 0, key)
}

func (a *Hamt) Get(key any) any {
	v, _ := a.get(hash(key), 0, key)
	return v
}

func (a *Hamt) Range(f func(any, any) bool) {
	a.foreach(f)
}

func (a *Hamt) String() string {
	var b strings.Builder
	b.WriteByte('{')
	sep := ""
	a.foreach(func(k, v any) bool {
		if sep == "" {
			sep = ", "
		} else {
			b.WriteString(sep)
		}
		_, err := fmt.Fprintf(&b, "%#v: %#v", k, v)
		return err == nil
	})
	b.WriteByte('}')
	return b.String()
}

func (a Hamt) Put(key, value any) *Hamt {
	a.amt = a.put(hash(key), 0, key, value)
	return &a
}

func (a Hamt) Delete(key any) *Hamt {
	a.amt = a.delete(hash(key), 0, key)
	return &a
}
