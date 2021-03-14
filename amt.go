package immutable

import (
	"fmt"
	"math/bits"
	"strings"
)

// any takes 16 bytes on 64bit archs
type any = interface{}

// amt takes 32 + len(entries) * 8 bytes on 64bit archs
type amt struct {
	bits    uint32   // 8 bytes on 64bit archs
	entries []*entry // 24 + len(entries) * 8 bytes on 64bit archs
}

// entry takes 40 bytes on 64bit archs
type entry struct {
	prefix uint32 // 8 bytes on 64bit archs
	ref    any    // 16 bytes on 64bit archs
	value  any    // 16 bytes on 64bit archs
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

func (n amt) len() int {
	len := len(n.entries)
	for _, e := range n.entries {
		if a, ok := e.ref.(amt); ok {
			len += a.len() - 1
		}
	}
	return len
}

func (n amt) depth() int {
	depth := 0
	for _, e := range n.entries {
		if a, ok := e.ref.(amt); ok {
			d := a.depth()
			if d > depth {
				depth = d
			}
		}
	}
	return 1 + depth
}

func (n amt) size() int {
	const amtsize = 32
	const ptrsize = 8
	const entrysize = 40
	size := amtsize + (ptrsize + entrysize)*len(n.entries)
	for _, e := range n.entries {
		if a, ok := e.ref.(amt); ok {
			size += a.size()
		}
	}
	return size
}

func (n amt) lookup(prefix uint32, shift uint8, key any) (any, bool) {
	for {
		bitpos := bitpos(prefix, shift)
		if present(n.bits, bitpos) {
			e := n.entries[index(n.bits, bitpos)]
			if a, ok := e.ref.(amt); ok {
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
		return nil, false
	}
}

func (n amt) foreach(f func(key, value any) bool) {
	for _, e := range n.entries {
		if a, ok := e.ref.(amt); ok {
			a.foreach(f)
		} else {
			if !f(e.ref, e.value) {
				return
			}
		}
	}
}

func (n amt) string() string {
	var b strings.Builder
	b.WriteByte('{')
	sep := ""
	n.foreach(func(k, v any) bool {
		if sep == "" {
			sep = ", "
		} else {
			b.WriteString(sep)
		}
		_, err := fmt.Fprintf(&b, "%#v:%#v", k, v)
		return err == nil
	})
	b.WriteByte('}')
	return b.String()
}

func (n amt) set(prefix uint32, shift uint8, key, value any) amt {
	bitpos := bitpos(prefix, shift)
	if present(n.bits, bitpos) {
		index := index(n.bits, bitpos)
		entries := make([]*entry, len(n.entries))
		copy(entries, n.entries)
		n.entries = entries
		e := n.entries[index]
		if a, ok := e.ref.(amt); ok {
			n.entries[index] = &entry{ref: a.set(prefix, shift+nextlevel, key, value)}
		} else {
			if e.prefix == prefix && e.ref == key {
				n.entries[index] = &entry{prefix, key, value}
			} else {
				// replace item with a new amt node holding the 2 items
				n.entries[index] = &entry{ref: amt{}.
					set(e.prefix, shift+nextlevel, e.ref, e.value).
					set(prefix, shift+nextlevel, key, value)}
			}
		}
	} else if shift < collision {
		index := index(n.bits, bitpos)
		entries := make([]*entry, len(n.entries)+1)
		n.bits |= bitpos
		copy(entries, n.entries[:index])
		copy(entries[index+1:], n.entries[index:])
		entries[index] = &entry{prefix, key, value}
		n.entries = entries
	} else {
		entries := make([]*entry, len(n.entries))
		copy(entries, n.entries)
		n.entries = entries
		for index, e := range n.entries {
			if e.ref == key {
				n.entries[index] = &entry{prefix, key, value}
				return n
			}
		}
		n.entries = append(n.entries, &entry{prefix, key, value})
	}
	return n
}

func (n amt) delete(prefix uint32, shift uint8, key any) amt {
	bitpos := bitpos(prefix, shift)
	if present(n.bits, bitpos) {
		index := index(n.bits, bitpos)
		e := n.entries[index]
		if a, ok := e.ref.(amt); ok {
			entries := make([]*entry, len(n.entries))
			copy(entries, n.entries)
			if a = a.delete(prefix, shift+nextlevel, key); a.len() == 1 {
				entries[index] = a.entries[0]
			} else {
				entries[index] = &entry{ref: a}
			}
			n.entries = entries
		} else {
			if e.prefix == prefix && e.ref == key {
				entries := make([]*entry, len(n.entries)-1)
				copy(entries, n.entries[:index])
				copy(entries[index:], n.entries[index+1:])
				n.bits &= ^bitpos
				n.entries = entries
			}
		}
	} else if shift == collision {
		for index, e := range n.entries {
			if e.prefix == prefix && e.ref == key {
				if index+1 == len(n.entries) {
					n.entries = n.entries[:index]
				} else {
					entries := make([]*entry, len(n.entries)-1)
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
