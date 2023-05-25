package immutable

import (
	"testing"
	"unsafe"
)

func TestSize(t *testing.T) {
	const arch = int(2 - uint64(^uint(0))>>63)

	ints0 := amt[uint8, uint8]{}
	ints1 := ints0.set(hash(123), 0, 123, 42)
	ints2 := ints1.set(hash(124), 0, 124, 69)
	assert.EqualInt(t, 32/arch, SizeAMT(ints0), "SizeAMT(ints0)")
	assert.EqualInt(t, 64/arch, SizeAMT(ints1), "SizeAMT(ints1)")
	assert.EqualInt(t, 96/arch, SizeAMT(ints2), "SizeAMT(ints2)")

	assert.EqualInt(t, 32/arch, int(unsafe.Sizeof(amt[string, any]{})), "unsafe.Sizeof(amt{})")
	assert.EqualInt(t, 40/arch, int(unsafe.Sizeof(entry[any]{})), "unsafe.Sizeof(entry{})")
	assert.EqualInt(t, 32/arch, int(unsafe.Sizeof(Map[string, any]{})), "unsafe.Sizeof(Map{})")
	assert.EqualInt(t, 40/arch, int(unsafe.Sizeof(MapX[any, any]{})), "unsafe.Sizeof(MapX{})")

	t0 := &amt[any, string]{}
	t1 := t0.set(0, 0, "Hello", "World!")

	assert.EqualInt(t, 8/arch, int(unsafe.Sizeof(t1.entries[0])), "unsafe.Sizeof(t1.entries[0])")
	assert.EqualInt(t, 32/arch, SizeAMT(*t0), "t0.size()")
	assert.EqualInt(t, (32+8+40)/arch, SizeAMT(t1), "t1.size()")
	assert.EqualInt(t, 1, t1.len(), "t1.Len()")
	assert.EqualInt(t, 1, t1.depth(), "t1.Depth()")

	EnableHashCollision = true
	m0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})
	m1 := m0.Set("Hello", "World!")
	m2 := m1.Set("He11o", "There!")

	assert.EqualInt(t, (8+32)/arch, SizeMapX(m0), "SizeMapX(m0)")
	assert.EqualInt(t, 1, m1.Len(), "m1.Len()")
	assert.EqualInt(t, 1, m1.Depth(), "m1.Depth()")
	assert.EqualInt(t, (8+(32+8+40))/arch, SizeMapX(m1), "SizeMapX(m1)")
	assert.EqualInt(t, 2, m2.Len(), "m2.Len()")
	assert.EqualInt(t, 4, m2.Depth(), "m2.Depth()")
	assert.EqualInt(t, (8+(32+8+40)+(32+8+40)+(32+8+40)+(32+2*(8+40)))/arch, SizeMapX(m2), "SizeMapX(m2)")
}

// SizeMap returns the number of bytes used for storing the entries, not
// including the size of the actual keys and the values.
func SizeMap[K comparable, V any](m Map[K, V]) int {
	return SizeAMT(m.amt)
}

// SizeMapX returns the number of bytes used for storing the entries, not
// including the size of the actual keys and the values.
func SizeMapX[K comparable, V any](m MapX[K, V]) int {
	const arch = int(2 - uint64(^uint(0))>>63)
	return 8/arch + SizeAMT(m.amt)
}

func SizeAMT[K comparable, V any](a amt[K, V]) int {
	// const arch = int(2 - uint64(^uint(0))>>63)
	// const amtsize = 32 / arch
	amtsize := int(unsafe.Sizeof(a))
	// const ptrsize = 8 / arch
	ptrsize := int(unsafe.Sizeof(&a))
	// const entrysize = 40 / arch
	entrysize := int(unsafe.Sizeof(entry[V]{}))
	size := amtsize + len(a.entries)*(ptrsize+entrysize)
	for _, e := range a.entries {
		if a, ok := e.ref.(amt[K, V]); ok {
			size += SizeAMT(a)
		}
	}
	return size
}
