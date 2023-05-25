package immutable

// Map is a persistent immutable hash array mapped trie (HAMT) with an
// internal hash function. The key types it supports are either string
// or any integer type. Keys are directly compared using the '==' operator.
// Key types other than string or integers need an external hasher.
type Map[K comparable, V any] struct{ amt[K, V] }

// Len returns the number of entries that are present.
func (a Map[K, V]) Len() int {
	return a.len()
}

// Depth returns the number of levels in the amt.
// Calling Depth on an empty amt returns 1.
func (a Map[K, V]) Depth() int {
	return a.depth()
}

// Lookup returns the value of an entry associated with a given key along with
// the value true when the key is present. Otherwise it returns (zero, false).
func (a Map[K, V]) Lookup(key K) (V, bool) {
	return a.lookup(hash(key), 0, key)
}

// Has returns true when an entry with the given key is present.
func (a Map[K, V]) Has(key K) bool {
	_, b := a.lookup(hash(key), 0, key)
	return b
}

// Get returns the value for the entry with the given key or zero value
// when it is not present.
func (a Map[K, V]) Get(key K) V {
	v, _ := a.lookup(hash(key), 0, key)
	return v
}

// Range calls the given function for every key,value pair present.
func (a Map[K, V]) Range(f func(K, V) bool) {
	a.foreach(f)
}

// String returns a string representation of the key,value pairs present.
func (a Map[K, V]) String() string {
	return a.string()
}

// Set returns a copy of the Map with the given key,value pair inserted.
func (a Map[K, V]) Set(key K, value V) Map[K, V] {
	return Map[K, V]{a.set(hash(key), 0, key, value)}
}

// Del returns a copy of the Map with the entry for the key removed.
func (a Map[K, V]) Del(key K) Map[K, V] {
	return Map[K, V]{a.delete(hash(key), 0, key)}
}

// MapX is a persistent immutable hash array mapped trie (HAMT) with an
// external key marshal function. The marshal function will map the key
// to a byte slice. The byteslice is passed to the internal hash function
// for the Map. The key itself is stored in the Map verbatim and actually
// used in compare operations. The marshaled key is only used for hashing.
type MapX[K comparable, V any] struct {
	amt[K, V]
	marshal func(any) ([]byte, error)
}

func MapWith[K comparable, V any](marshal func(any) ([]byte, error)) MapX[K, V] {
	return MapX[K, V]{amt[K, V]{}, marshal}
}

// Len returns the number of entries that are present.
func (a MapX[K, V]) Len() int {
	return a.len()
}

// Depth returns the number of levels in the Map.
// Calling Depth on an empty amt returns 1.
func (a MapX[K, V]) Depth() int {
	return a.depth()
}

// Lookup returns the value of an entry associated with a given key along with
// the value true when the key is present. Otherwise it returns (nil, false).
func (a MapX[K, V]) Lookup(key K) (V, bool) {
	k, e := a.marshal(key)
	if e != nil {
		panic(UnhashableKeyType)
	}
	return a.lookup(hash(k), 0, key)
}

// Has returns true when an entry with the given key is present.
func (a MapX[K, V]) Has(key K) bool {
	k, e := a.marshal(key)
	if e != nil {
		panic(UnhashableKeyType)
	}
	_, b := a.lookup(hash(k), 0, key)
	return b
}

// Get returns the value for the entry with the given key or nil when it is
// not present.
func (a MapX[K, V]) Get(key K) V {
	k, e := a.marshal(key)
	if e != nil {
		panic(UnhashableKeyType)
	}
	v, _ := a.lookup(hash(k), 0, key)
	return v
}

// Range calls the given function for every key,value pair present.
func (a MapX[K, V]) Range(f func(K, V) bool) {
	a.foreach(f)
}

// String returns a string representation of the key,value pairs present.
func (a MapX[K, V]) String() string {
	return a.string()
}

// Set returns a copy of the Map with the given key,value pair inserted.
func (a MapX[K, V]) Set(key K, value V) MapX[K, V] {
	k, e := a.marshal(key)
	if e != nil {
		panic(UnhashableKeyType)
	}
	return MapX[K, V]{a.set(hash(k), 0, key, value), a.marshal}
}

// Del returns a copy of the Map with the entry for the key removed.
func (a MapX[K, V]) Del(key K) MapX[K, V] {
	k, e := a.marshal(key)
	if e != nil {
		panic(UnhashableKeyType)
	}
	return MapX[K, V]{a.delete(hash(k), 0, key), a.marshal}
}
