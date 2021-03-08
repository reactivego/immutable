package immutable

import (
	"math/bits"
)

const collision = 35
const nextlevel = 5

func mask(prefix uint32, shift uint8) uint8 {
	return uint8(prefix>>shift) & 0x1f
}

func bitpos(pos uint8) uint32 {
	return uint32(1) << pos
}

func index(bitmap uint32, bitpos uint32) int {
	return bits.OnesCount32(bitmap & (bitpos - 1))
}

func present(bitmap uint32, bitpos uint32) bool {
	return bitmap&bitpos != 0
}

// any takes 16 bytes on 64bit archs
type any = interface{}

// amt takes 32 + len(entries) * 16 bytes on 64bit archs
type amt struct {
	bits    uint32 // 8 bytes on 64bit archs
	entries []any  // 24 + len(entries) * 16 bytes on 64bit archs
}

// item takes 40 bytes on 64bit archs
type item struct {
	prefix uint32 // 8 bytes on 64bit archs
	key    any    // 16 bytes on 64bit archs
	value  any    // 16 bytes on 64bit archs
}

func (n *amt) len() int {
	size := len(n.entries)
	for _, entry := range n.entries {
		if e, ok := entry.(*amt); ok {
			size += e.len() - 1
		}
	}
	return size
}

func (n *amt) depth() int {
	depth := 0
	for _, entry := range n.entries {
		if e, ok := entry.(*amt); ok {
			d := e.depth()
			if d > depth {
				depth = d
			}
		}
	}
	return 1 + depth
}

func (n *amt) size() int {
	const amtsize = 32
	const anysize = 16
	const itemsize = 40
	size := amtsize + anysize*len(n.entries)
	for _, entry := range n.entries {
		switch e := entry.(type) {
		case item:
			size += itemsize
		case *amt:
			size += e.size()
		}
	}
	return size
}

func (n *amt) get(prefix uint32, shift uint8, key any) (any, bool) {
	if shift == collision {
		for _, entry := range n.entries {
			if e := entry.(item); e.prefix == prefix && e.key == key {
				return e.value, true
			}
		}
		return nil, false
	}
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		switch e := n.entries[index(n.bits, bitpos)].(type) {
		case item:
			if e.prefix == prefix && e.key == key {
				return e.value, true
			}
		case *amt:
			return e.get(prefix, shift+nextlevel, key)
		}
	}
	return nil, false
}

func (n *amt) foreach(f func(key, value any) bool) {
	for _, entry := range n.entries {
		switch e := entry.(type) {
		case item:
			if !f(e.key, e.value) {
				return
			}
		case *amt:
			e.foreach(f)
		}
	}
}

func (n amt) put(prefix uint32, shift uint8, key, value any) *amt {
	if shift == collision {
		entries := make([]any, len(n.entries))
		copy(entries, n.entries)
		n.entries = entries
		for index, entry := range n.entries {
			if e := entry.(item); e.key == key {
				n.entries[index] = item{prefix: prefix, key: key, value: value}
				return &n
			}
		}
		n.entries = append(n.entries, item{prefix: prefix, key: key, value: value})
		return &n
	}
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		index := index(n.bits, bitpos)
		entries := make([]any, len(n.entries))
		copy(entries, n.entries)
		n.entries = entries
		switch e := n.entries[index].(type) {
		case item:
			if e.prefix == prefix && e.key == key {
				n.entries[index] = item{prefix: prefix, key: key, value: value}
			} else {
				// replace item with a new amt node holding the 2 items
				n.entries[index] = amt{}.
					put(e.prefix, shift+nextlevel, e.key, e.value).
					put(prefix, shift+nextlevel, key, value)
			}
		case *amt:
			n.entries[index] = e.put(prefix, shift+nextlevel, key, value)
		}
	} else {
		index := index(n.bits, bitpos)
		e := make([]any, len(n.entries)+1)
		n.bits |= bitpos
		copy(e, n.entries[:index])
		copy(e[index+1:], n.entries[index:])
		e[index] = item{prefix: prefix, key: key, value: value}
		n.entries = e
	}
	return &n
}

func (n amt) delete(prefix uint32, shift uint8, key any) *amt {
	if shift == collision {
		for index, entry := range n.entries {
			if e := entry.(item); e.prefix == prefix && e.key == key {
				if index+1 == len(n.entries) {
					n.entries = n.entries[:index]
				} else {
					entries := make([]any, len(n.entries)-1)
					copy(entries, n.entries[:index])
					copy(entries[index:], n.entries[index+1:])
					n.entries = entries
				}
				return &n
			}
		}
		return &n
	}
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		index := index(n.bits, bitpos)
		switch e := n.entries[index].(type) {
		case item:
			if e.prefix == prefix && e.key == key {
				entries := make([]any, len(n.entries)-1)
				copy(entries, n.entries[:index])
				copy(entries[index:], n.entries[index+1:])
				n.bits &= ^bitpos
				n.entries = entries
			}
		case *amt:
			entries := make([]any, len(n.entries))
			copy(entries, n.entries)
			if e = e.delete(prefix, shift+nextlevel, key); len(e.entries) == 1 {
				entries[index] = e.entries[0]
			} else {
				entries[index] = e
			}
			n.entries = entries
		}
	}
	return &n
}
