package immutable

import (
	"fmt"
	"math/bits"
	"strings"
)

// amt takes 32 + len(entries) * 8 bytes on 64bit archs
type amt[K comparable, V any] struct {
	bits    uint32      // 8 bytes on 64bit archs
	entries []*entry[V] // 24 + len(entries) * 8 bytes on 64bit archs
}

// entry takes 24..40 bytes on 64bit archs depending on
// actual type V of the value
type entry[V any] struct {
	prefix uint32 // 4..8 bytes on 64bit archs
	value  V      // 4..16 bytes on 64bit archs
	ref    any    // 16 bytes on 64bit archs
}

const collision = 35
const nextlevel = 5

func bitpos(prefix uint32, shift uint8) uint32 {
	return 1 << (prefix >> shift & 0x1f)
}

func index(bitmap uint32, bitpos uint32) int {
	return bits.OnesCount32(bitmap & (bitpos - 1))
}

func present(bitmap uint32, bitpos uint32) bool {
	return bitmap&bitpos != 0
}

func (n amt[K, V]) len() int {
	len := len(n.entries)
	for _, e := range n.entries {
		if a, ok := e.ref.(amt[K, V]); ok {
			len += a.len() - 1
		}
	}
	return len
}

func (n amt[K, V]) depth() int {
	depth := 0
	for _, e := range n.entries {
		if a, ok := e.ref.(amt[K, V]); ok {
			d := a.depth()
			if d > depth {
				depth = d
			}
		}
	}
	return 1 + depth
}

func (n amt[K, V]) lookup(prefix uint32, shift uint8, key K) (V, bool) {
	for {
		bitpos := bitpos(prefix, shift)
		if present(n.bits, bitpos) {
			e := n.entries[index(n.bits, bitpos)]
			if a, ok := e.ref.(amt[K, V]); ok {
				n = a
				shift += nextlevel
				continue
			}
			if e.prefix == prefix && e.ref == key {
				return e.value, true
			}
		} else if shift == collision {
			for _, e := range n.entries {
				if e.prefix == prefix && e.ref == key {
					return e.value, true
				}
			}
		}
		var zero V
		return zero, false
	}
}

func (n amt[K, V]) foreach(f func(K, V) bool) {
	for _, e := range n.entries {
		if a, ok := e.ref.(amt[K, V]); ok {
			a.foreach(f)
		} else {
			if !f(e.ref.(K), e.value) {
				return
			}
		}
	}
}

func (n amt[K, V]) string() string {
	var b strings.Builder
	b.WriteByte('{')
	f := "%+v:%+v"
	n.foreach(func(k K, v V) bool {
		_, err := fmt.Fprintf(&b, f, k, v)
		f = ", %+v:%+v"
		return err == nil
	})
	b.WriteByte('}')
	return b.String()
}

func (n amt[K, V]) set(prefix uint32, shift uint8, key K, value V) amt[K, V] {
	bitpos := bitpos(prefix, shift)
	if present(n.bits, bitpos) {
		index := index(n.bits, bitpos)
		entries := make([]*entry[V], len(n.entries))
		copy(entries, n.entries)
		n.entries = entries
		e := n.entries[index]
		if a, ok := e.ref.(amt[K, V]); ok {
			n.entries[index] = &entry[V]{ref: a.set(prefix, shift+nextlevel, key, value)}
		} else {
			if e.prefix == prefix && e.ref == key {
				n.entries[index] = &entry[V]{prefix, value, key}
			} else {
				// replace item with a new amt node holding the 2 items
				n.entries[index] = &entry[V]{ref: amt[K, V]{}.
					set(e.prefix, shift+nextlevel, e.ref.(K), e.value).
					set(prefix, shift+nextlevel, key, value)}
			}
		}
	} else if shift < collision {
		index := index(n.bits, bitpos)
		entries := make([]*entry[V], len(n.entries)+1)
		n.bits |= bitpos
		copy(entries, n.entries[:index])
		copy(entries[index+1:], n.entries[index:])
		entries[index] = &entry[V]{prefix, value, key}
		n.entries = entries
	} else {
		entries := make([]*entry[V], len(n.entries))
		copy(entries, n.entries)
		n.entries = entries
		for index, e := range n.entries {
			if e.ref == key {
				n.entries[index] = &entry[V]{prefix, value, key}
				return n
			}
		}
		n.entries = append(n.entries, &entry[V]{prefix, value, key})
	}
	return n
}

func (n amt[K, V]) delete(prefix uint32, shift uint8, key K) amt[K, V] {
	bitpos := bitpos(prefix, shift)
	if present(n.bits, bitpos) {
		index := index(n.bits, bitpos)
		e := n.entries[index]
		if a, ok := e.ref.(amt[K, V]); ok {
			entries := make([]*entry[V], len(n.entries))
			copy(entries, n.entries)
			if a = a.delete(prefix, shift+nextlevel, key); a.len() == 1 {
				entries[index] = a.entries[0]
			} else {
				entries[index] = &entry[V]{ref: a}
			}
			n.entries = entries
		} else {
			if e.prefix == prefix && e.ref == key {
				if index+1 == len(n.entries) {
					n.entries = n.entries[:index]
				} else {
					entries := make([]*entry[V], len(n.entries)-1)
					copy(entries, n.entries[:index])
					copy(entries[index:], n.entries[index+1:])
					n.entries = entries
				}
				n.bits &= ^bitpos
			}
		}
	} else if shift == collision {
		for index, e := range n.entries {
			if e.prefix == prefix && e.ref == key {
				if index+1 == len(n.entries) {
					n.entries = n.entries[:index]
				} else {
					entries := make([]*entry[V], len(n.entries)-1)
					copy(entries, n.entries[:index])
					copy(entries[index:], n.entries[index+1:])
					n.entries = entries
				}
				return n
			}
		}
	}
	return n
}
