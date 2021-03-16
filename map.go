package immutable

import "hash/maphash"

type Any = interface{}

type MapError string

func (e MapError) Error() string {
	return string(e)
}

const UnhashableKeyType = MapError("Unhashable Key Type")

// Hamt is a persistent immutable hash array mapped trie (HAMT) with an
// internal hash function that uses the standard "hash/maphash" package for
// hashing. The key types it supports are either string or any integer type.
// Keys are directly compared using the '==' operator. Key types other than
// string or integers need an external hasher. Use the method WithHasher to
// create a HAMT with an external hasher.
type Hamt struct{ amt }

// Map is an empty Hamt
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

// Len returns the number of entries that are present.
func (a Hamt) Len() int {
	return a.len()
}

// Depth returns the number of levels in the Hamt. Calling Depth on an empty
// Hamt returns 1.
func (a Hamt) Depth() int {
	return a.depth()
}

// Size returns the number of bytes used for storing the entries, not
// including the size of the actual keys and the values.
func (a Hamt) Size() int {
	return a.size()
}

// Lookup returns the value of an entry associated with a given key along with
// the value true when the key is present. Otherwise it returns (nil, false).
func (a Hamt) Lookup(key Any) (Any, bool) {
	return a.lookup(hash(key), 0, key)
}

// Has returns true when an entry with the given key is present.
func (a Hamt) Has(key Any) bool {
	_, b := a.lookup(hash(key), 0, key)
	return b
}

// Get returns the value for the entry with the given key or nil when it is
// not present.
func (a Hamt) Get(key Any) Any {
	v, _ := a.lookup(hash(key), 0, key)
	return v
}

// Range calls the given function for every key,value pair present.
func (a Hamt) Range(f func(Any, Any) bool) {
	a.foreach(f)
}

// String returns a string representation of the key,value pairs present.
func (a Hamt) String() string {
	return a.string()
}

// Set returns a copy of the Hamt with the given key,value pair inserted.
func (a Hamt) Set(key, value Any) Hamt {
	return Hamt{a.set(hash(key), 0, key, value)}
}

// Put returns a copy of the Hamt with the key,key pair inserted. So the key
// is also inserted as the value.
func (a Hamt) Put(key Any) Hamt {
	return Hamt{a.set(hash(key), 0, key, key)}
}

// Del returns a copy of the Hamt with the entry for the key removed.
func (a Hamt) Del(key Any) Hamt {
	return Hamt{a.delete(hash(key), 0, key)}
}

// WithHasher returns an empty HamtX with the given hasher function. The
// hasher function is used to convert a key into a hash an a key. So it allows
// for key transformation.
func (a Hamt) WithHasher(hash func(Any) (uint32, Any)) HamtX {
	return HamtX{hash: hash}
}

// HamtX is a Hash Array Mapped Trie with an external hash function.
type HamtX struct {
	amt
	hash func(Any) (uint32, Any)
}

// Len returns the number of entries that are present.
func (a HamtX) Len() int {
	return a.len()
}

// Depth returns the number of levels in the Hamt. Calling Depth on an empty
// Hamt returns 1.
func (a HamtX) Depth() int {
	return a.depth()
}

// Size returns the number of bytes used for storing the entries, not
// including the size of the actual keys and the values.
func (a HamtX) Size() int {
	return 8 + a.size()
}

// Lookup returns the value of an entry associated with a given key along with
// the value true when the key is present. Otherwise it returns (nil, false).
func (a HamtX) Lookup(key Any) (Any, bool) {
	h, k := a.hash(key)
	return a.lookup(h, 0, k)
}

// Has returns true when an entry with the given key is present.
func (a HamtX) Has(key Any) bool {
	h, k := a.hash(key)
	_, b := a.lookup(h, 0, k)
	return b
}

// Get returns the value for the entry with the given key or nil when it is
// not present.
func (a HamtX) Get(key Any) Any {
	h, k := a.hash(key)
	v, _ := a.lookup(h, 0, k)
	return v
}

// Range calls the given function for every key,value pair present.
func (a HamtX) Range(f func(Any, Any) bool) {
	a.foreach(f)
}

// String returns a string representation of the key,value pairs present.
func (a HamtX) String() string {
	return a.string()
}

// Set returns a copy of the HamtX with the given key,value pair inserted.
func (a HamtX) Set(key, value Any) HamtX {
	h, k := a.hash(key)
	return HamtX{a.set(h, 0, k, value), a.hash}
}

// Put returns a copy of the HamtX with the key,key pair inserted. So the key
// is also inserted as the value.
func (a HamtX) Put(key Any) HamtX {
	h, k := a.hash(key)
	return HamtX{a.set(h, 0, k, key), a.hash}
}

// Del returns a copy of the HamtX with the entry for the key removed.
func (a HamtX) Del(key Any) HamtX {
	h, k := a.hash(key)
	return HamtX{a.delete(h, 0, k), a.hash}
}
