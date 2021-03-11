package immutable

import (
	"hash/maphash"
)

type Any = interface{}

type MapError string

func (e MapError) Error() string {
	return string(e)
}

const UnhashableKeyType = MapError("Unhashable Key Type")

// Hamt is a Hash Array Mapped Trie with an internal hash function.
type Hamt struct{ amt }

var Map = Hamt{}

var seed = maphash.MakeSeed()

func hash(key Any) uint32 {
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
	default:
		panic(UnhashableKeyType)
	}
}

func (a Hamt) Len() int {
	return a.len()
}

func (a Hamt) Depth() int {
	return a.depth()
}

func (a Hamt) Size() int {
	return 8 + a.size()
}

func (a Hamt) Lookup(key Any) (Any, bool) {
	return a.lookup(hash(key), 0, key)
}

func (a Hamt) Has(key Any) bool {
	_, b := a.lookup(hash(key), 0, key)
	return b
}

func (a Hamt) Get(key Any) Any {
	v, _ := a.lookup(hash(key), 0, key)
	return v
}

func (a Hamt) Range(f func(Any, Any) bool) {
	a.foreach(f)
}

func (a Hamt) String() string {
	return "Hamt" + a.string()
}

func (a Hamt) Put(key, value Any) Hamt {
	return Hamt{a.put(hash(key), 0, key, value)}
}

func (a Hamt) Set(key Any) Hamt {
	return Hamt{a.put(hash(key), 0, key, nil)}
}

func (a Hamt) Del(key Any) Hamt {
	return Hamt{a.delete(hash(key), 0, key)}
}

func (a Hamt) WithHasher(hash func(Any) (uint32, Any)) HamtX {
	return HamtX{hash: hash}
}

// HamtX is a Hash Array Mapped Trie with an eXternal hash function.
type HamtX struct {
	amt
	hash func(Any) (uint32, Any)
}

func (a HamtX) Len() int {
	return a.len()
}

func (a HamtX) Depth() int {
	return a.depth()
}

func (a HamtX) Size() int {
	return 8 + a.size()
}

func (a HamtX) Lookup(key Any) (Any, bool) {
	h, k := a.hash(key)
	return a.lookup(h, 0, k)
}

func (a HamtX) Has(key Any) bool {
	h, k := a.hash(key)
	_, b := a.lookup(h, 0, k)
	return b
}

func (a HamtX) Get(key Any) Any {
	h, k := a.hash(key)
	v, _ := a.lookup(h, 0, k)
	return v
}

func (a HamtX) Range(f func(Any, Any) bool) {
	a.foreach(f)
}

func (a HamtX) String() string {
	return "HamtX" + a.string()
}

func (a HamtX) Put(key, value Any) HamtX {
	h, k := a.hash(key)
	return HamtX{a.put(h, 0, k, value), a.hash}
}

func (a HamtX) Set(key Any) HamtX {
	h, k := a.hash(key)
	return HamtX{a.put(h, 0, k, nil), a.hash}
}

func (a HamtX) Del(key Any) HamtX {
	h, k := a.hash(key)
	return HamtX{a.delete(h, 0, k), a.hash}
}
