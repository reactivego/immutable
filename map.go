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

const UnhashableKeyType = MapError("Unhashable Key Type")

type Hasher interface {
	Sum32() uint32
}

var seed = maphash.MakeSeed()

func hash(key any) uint32 {
	switch k := key.(type) {
	case string:
		var h maphash.Hash
		h.SetSeed(seed)
		h.WriteString(k)
		return uint32(h.Sum64()) & 0xFFFFFFFF
	case int8:
		return uint32(k)
	case uint8:
		return uint32(k)
	case int16:
		return uint32(k)
	case uint16:
		return uint32(k)
	case int32:
		return uint32(k)
	case uint32:
		return k
	case int64:
		return uint32(k) & 0xFFFFFFFF
	case uint64:
		return uint32(k) & 0xFFFFFFFF
	case int:
		return uint32(k) & 0xFFFFFFFF
	case uint:
		return uint32(k) & 0xFFFFFFFF
	case Hasher:
		return k.Sum32()
	default:
		panic(UnhashableKeyType)
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
