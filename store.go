package immutable

// Store is a Hash Array Mapped Trie with an external split function.
type Store[D any, K comparable, V any] struct {
	amt[K, V]
	split func(D) (K, V)
}

func StoreWith[D any, K comparable, V any](splitter func(D) (K, V)) Store[D, K, V] {
	return Store[D, K, V]{amt[K, V]{}, splitter}
}

// Len returns the number of entries that are present.
func (a Store[D, K, V]) Len() int {
	return a.len()
}

// Depth returns the number of levels in the Hamt. Calling Depth on an empty
// Hamt returns 1.
func (a Store[D, K, V]) Depth() int {
	return a.depth()
}

// Has returns true when an entry with the given key is present.
func (a Store[D, K, V]) Has(data D) bool {
	k, _ := a.split(data)
	_, b := a.lookup(hash(k), 0, k)
	return b
}

// Get returns the value for the entry with the given key or nil when it is
// not present.
func (a Store[D, K, V]) Get(data D) any {
	k, _ := a.split(data)
	v, _ := a.lookup(hash(k), 0, k)
	return v
}

// Range calls the given function for every key,value pair present.
func (a Store[D, K, V]) Range(f func(K, V) bool) {
	a.foreach(f)
}

// String returns a string representation of the key,value pairs present.
func (a Store[D, K, V]) String() string {
	return a.string()
}

// Put returns a copy of the Set with the key as part of the set.
func (a Store[D, K, V]) Put(data D) Store[D, K, V] {
	k, v := a.split(data)
	return Store[D, K, V]{a.set(hash(k), 0, k, v), a.split}
}

// Del returns a copy of the Store with the key removed from the set.
func (a Store[D, K, V]) Del(data D) Store[D, K, V] {
	k, _ := a.split(data)
	return Store[D, K, V]{a.delete(hash(k), 0, k), a.split}
}
