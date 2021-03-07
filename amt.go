package immutable

import (
	"math/bits"
)

const collision = 35
const nextlevel = 5

func mask(prefix uint32, shift uint8) uint8 {
	return uint8(prefix>>shift) & 0x1f
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
		if node, ok := entry.(*amt); ok {
			size += node.len() - 1
		}
	}
	return size
}

func (n *amt) depth() int {
	depth := 0
	for _, entry := range n.entries {
		if node, ok := entry.(*amt); ok {
			ndepth := node.depth()
			if ndepth > depth {
				depth = ndepth
			}
		}
	}
	return 1 + depth
}

func (n *amt) get(prefix uint32, shift uint8, key any) (any, bool) {
	if shift == collision {
		for _, entry := range n.entries {
			if e := entry.(item); e.key == key {
				return e.value, true
			}
		}
		return nil, false
	}
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		d := n.entries[index(n.bits, bitpos)]
		switch e := d.(type) {
		case item:
			return e.value, true
		case *amt:
			return e.get(prefix, shift+nextlevel, key)
		}
	}
	return nil, false
}

func (n *amt) foreach(f func(key, value any) bool) {
	for _, e := range n.entries {
		switch e := e.(type) {
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
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		// replace
		index := index(n.bits, bitpos)
		e := make([]any, len(n.entries))
		copy(e, n.entries)
		n.entries = e
		entrynode := n.entries[index]
		if entry, ok := entrynode.(item); ok {
			if entry.prefix == prefix && entry.key == key {
				// replace
				n.entries[index] = item{prefix: prefix, key: key, value: value}
			} else {
				if entry.prefix == prefix && shift == collision {
					// prefix collision, replace or insert an entry by enumerating entries
					for index, entry := range n.entries {
						if e := entry.(item); e.key == key {
							n.entries[index] = item{prefix: prefix, key: key, value: value}
							return &n
						}
					}
					n.entries = append(n.entries, item{prefix: prefix, key: key, value: value})
				} else {
					// prefix different or not at collision level, replace entry with node
					node := &amt{}
					node = node.put(entry.prefix, shift+nextlevel, entry.key, entry.value)
					node = node.put(prefix, shift+nextlevel, key, value)
					n.entries[index] = node
				}
			}
		} else {
			node := entrynode.(*amt)
			n.entries[index] = node.put(prefix, shift+nextlevel, key, value)
		}
	} else {
		// insert
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
			if e := entry.(item); e.key == key {
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
		d := n.entries[index(n.bits, bitpos)]
		switch e := d.(type) {
		case item:
			// delete
			index := index(n.bits, bitpos)
			entries := make([]any, len(n.entries)-1)
			n.bits &= ^bitpos
			copy(entries, n.entries[:index])
			copy(entries[index:], n.entries[index+1:])
			n.entries = entries
		case *amt:
			// delete
			index := index(n.bits, bitpos)
			entries := make([]any, len(n.entries))
			copy(entries, n.entries)
			if e := e.delete(prefix, shift+nextlevel, key); len(e.entries) == 1 {
				entries[index] = e.entries[0]
			} else {
				entries[index] = e
			}
			n.entries = entries
		}
	}
	return &n
}
