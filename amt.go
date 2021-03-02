package immutable

import (
	"bytes"
	"math/bits"
)

func mask(prefix uint32, shift uint8) uint8 {
	return uint8(prefix>>shift) & 0x1f
}

func index(bitmap uint32, bitpos uint32) int {
	return bits.OnesCount32(bitmap & (bitpos - 1))
}

func present(bitmap uint32, bitpos uint32) bool {
	return bitmap&bitpos != 0
}

type any = interface{}

// amt takes 32 bytes on 64bit archs
type amt struct {
	bits    uint32 // 8 bytes on 64bit archs
	entries []any  // 24 bytes on 64bit archs
}

// item takes 48 bytes on 64bit archs
type item struct {
	prefix uint32 // 8 bytes on 64bit archs
	key    []byte // 24 bytes on 64bit archs
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

func (n *amt) get(prefix uint32, shift uint8, k []byte) (any, bool) {
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		d := n.entries[index(n.bits, bitpos)]
		return d.(item).value, true
	}
	return nil, false
}

func (n *amt) foreach(f func(k []byte, v any)) {
	for _, e := range n.entries {
		switch e := e.(type) {
		case item:
			f(e.key, e.value)
		case *amt:
			e.foreach(f)
		}
	}
}

func (n amt) put(prefix uint32, shift uint8, key []byte, value any) *amt {
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		// replace
		index := index(n.bits, bitpos)
		e := make([]any, len(n.entries))
		copy(e, n.entries)
		n.entries = e
		entrynode := n.entries[index]
		if entry, ok := entrynode.(item); ok {
			if entry.prefix == prefix && bytes.Equal(entry.key, key) {
				// replace
				n.entries[index] = item{prefix: prefix, key: key, value: value}
			} else {
				if entry.prefix == prefix && shift > 32 {
					// prefix collision, replace or insert an entry by enumerating entries
					for index, entry := range n.entries {
						if e := entry.(item); bytes.Equal(e.key, key) {
							n.entries[index] = item{prefix: prefix, key: key, value: value}
							return &n
						}
					}
					n.entries = append(n.entries, item{prefix: prefix, key: key, value: value})
				} else {
					// prefix different or not at collision level, replace entry with node
					node := &amt{}
					node = node.put(entry.prefix, shift+5, entry.key, entry.value)
					node = node.put(prefix, shift+5, key, value)
					n.entries[index] = node
				}
			}
		} else {
			node := entrynode.(*amt)
			n.entries[index] = node.put(prefix, shift+5, key, value)
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

func (n amt) remove(prefix uint32, shift uint8, key []byte) *amt {
	bitpos := uint32(1) << mask(prefix, shift)
	if present(n.bits, bitpos) {
		// remove
		index := index(n.bits, bitpos)
		entries := make([]any, len(n.entries)-1)
		n.bits &= ^bitpos
		copy(entries, n.entries[:index])
		copy(entries[index:], n.entries[index+1:])
		n.entries = entries
	}
	return &n
}
