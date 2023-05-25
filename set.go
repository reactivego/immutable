package immutable

import (
	"fmt"
	"strings"
)

type Set[K comparable] struct{ amt[K, struct{}] }

// Len returns the number of entries that are present.
func (a Set[K]) Len() int {
	return a.len()
}

// Depth returns the number of levels in the Set.
// Calling Depth on an empty amt returns 1.
func (a Set[K]) Depth() int {
	return a.depth()
}

// Has returns true when an entry with the given key is present.
func (a Set[K]) Has(key K) bool {
	_, b := a.lookup(hash(key), 0, key)
	return b
}

// Range calls the given function for every (key,value) pair present.
func (a Set[K]) Range(f func(K) bool) {
	a.foreach(func(key K, _ struct{}) bool { return f(key) })
}

// String returns a string representation of the keys pairs present.
func (a Set[K]) String() string {
	var b strings.Builder
	b.WriteByte('{')
	f := "%+v"
	a.foreach(func(k K, _ struct{}) bool {
		_, err := fmt.Fprintf(&b, f, k)
		f = ", %+v"
		return err == nil
	})
	b.WriteByte('}')
	return b.String()
}

// Put returns a copy of the Set with the key added to it.
func (a Set[K]) Put(key K) Set[K] {
	return Set[K]{a.set(hash(key), 0, key, struct{}{})}
}

// Del returns a copy of the Set with the key removed from it.
func (a Set[K]) Del(key K) Set[K] {
	return Set[K]{a.delete(hash(key), 0, key)}
}
