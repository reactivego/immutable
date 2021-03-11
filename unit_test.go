package immutable

import (
	// "encoding/binary"

	"testing"
	"unsafe"
)

func TestDelDeep(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	k1 := "He1lo"
	v1 := "World!"

	k2 := "He2lo"
	v2 := "There!"

	t1 := t0.Put(k1, v1)
	t2 := t1.Put(k2, v2)
	t3 := t2.Delete(k1)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 1, t1.Depth(), "t1.Depth()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 4, t2.Depth(), "t2.Depth()")
	assert.EqualInt(t, 1, t3.Len(), "t3.Len()")
	assert.EqualInt(t, 1, t3.Depth(), "t3.Depth()")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, v2, t3.Get(k2), "t3.Get(k2)")
}

func TestGetDeep(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	k1 := "He1lo"
	v1 := "World!"

	k2 := "He2lo"
	v2 := "There!"

	t1 := t0.Put(k1, v1)
	t2 := t1.Put(k2, v2)
	r2 := t2.Get(k2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 4, t2.Depth(), "t2.Depth()")
	assert.Equal(t, v2, r2, "t2.Get(k2)")
}

func TestGetCollision(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	k1 := "Hello1"
	v1 := "World!"

	k2 := "Hello2"
	v2 := "There!"

	t1 := t0.Put(k1, v1)
	t2 := t1.Put(k2, v2)
	r2 := t2.Get(k2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 8, t2.Depth(), "t2.Depth()")
	assert.Equal(t, v2, r2, "t2.Get(k2)")
}

func TestPutCollision(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	k1 := "Hello"
	v1 := "World!"

	k2 := k1
	v2 := "There!"

	k3 := "Hela"
	v3 := "All!"

	k4 := "Hella"
	v4 := "Strange!"

	t1 := t0.Put(k1, v1)
	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	t2 := t1.Put(k2, v2)
	assert.EqualInt(t, 1, t2.Len(), "t2.Len()")
	t3 := t2.Put(k3, v3)
	assert.EqualInt(t, 2, t3.Len(), "t3.Len()")
	t4 := t3.Put(k4, v4)
	assert.EqualInt(t, 3, t4.Len(), "t4.Len()")
}

func TestBasicPutGetDelete(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	key := "hello"
	val := "world"

	t1 := t0.Put(key, val)
	assert.EqualInt(t, 0, t0.Len(), "t0.Len()")
	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	got, ok := t1.Lookup(key)
	assert.Equal(t, true, ok, "_,ok := t1.Lookup(%q); ok", key)
	assert.Equal(t, val, got, "v,_ := t1.Lookup(%q); v", key)
	t2 := t1.Delete(key)
	assert.EqualInt(t, 0, t0.Len(), "t0.Len()")
	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 0, t2.Len(), "t2.Len()")
	assert.Equal(t, true, t1.Get(key) != nil, "t1.Get(%q) != nil", key)
	assert.Equal(t, true, t2.Get(key) == nil, "t2.Get(%q) == nil", key)
}

func TestEntryPresent(t *testing.T) {
	tests := []struct{ exp, got bool }{
		/*0*/
		{exp: true, got: present(0b0001, bitpos(0, 0))},
		{exp: true, got: present(0b0010, bitpos(1, 0))},
		{exp: true, got: present(0b0100, bitpos(2, 0))},
		{exp: true, got: present(0b1000, bitpos(3, 0))},
		{exp: false, got: present(0b0010, bitpos(0, 0))},

		/*5*/
		{exp: false, got: present(0b0100, bitpos(0, 0))},
		{exp: true, got: present(0b0100, bitpos(2, 0))},
		{exp: false, got: present(0b0100, bitpos(3, 0))},
		{exp: true, got: present(0b1100, bitpos(3, 0))},
		{exp: false, got: present(0b10100, bitpos(3, 0))},

		/*10*/
		{exp: true, got: present(0b11100, bitpos(4, 0))},
		{exp: true, got: present(0b111100, bitpos(4, 0))},
		{exp: true, got: present(0b111100, bitpos(5, 0))},
		{exp: true, got: present(0b100100, bitpos(2, 0))},
		{exp: true, got: present(0b100100, bitpos(5, 0))},

		/*15*/
		{exp: false, got: present(0b1001010, bitpos(0, 0))},
		{exp: true, got: present(0b1001010, bitpos(1, 0))},
		{exp: false, got: present(0b1001010, bitpos(2, 0))},
		{exp: true, got: present(0b1001010, bitpos(3, 0))},
		{exp: false, got: present(0b1001010, bitpos(4, 0))},

		/*20*/
		{exp: false, got: present(0b1001010, bitpos(5, 0))},
		{exp: true, got: present(0b1001010, bitpos(6, 0))},
		{exp: false, got: present(0b1001010, bitpos(7, 0))},
		{exp: false, got: present(0b1001010, bitpos(7, 0))},
		{exp: false, got: present(0b1001010, bitpos(240, 0))},
	}
	for i, test := range tests {
		assert.Equal(t, test.exp, test.got, "test #%d", i)
	}
}

func TestEntryIndex(t *testing.T) {
	tests := []struct{ exp, got int }{
		/*0*/
		{exp: 0, got: index(0b0000, bitpos(0, 0))},
		{exp: 0, got: index(0b0000, bitpos(1, 0))},
		{exp: 0, got: index(0b0001, bitpos(0, 0))},
		{exp: 1, got: index(0b0001, bitpos(1, 0))},
		{exp: 1, got: index(0b0001, bitpos(2, 0))},

		/*5*/
		{exp: 0, got: index(0b0100, bitpos(0, 0))},
		{exp: 0, got: index(0b0100, bitpos(1, 0))},
		{exp: 0, got: index(0b0100, bitpos(2, 0))},
		{exp: 1, got: index(0b0100, bitpos(3, 0))},
		{exp: 1, got: index(0b0100, bitpos(4, 0))},

		/*10*/
		{exp: 3, got: index(0b011110, bitpos(4, 0))},
		{exp: 4, got: index(0b011110, bitpos(5, 0))},
		{exp: 4, got: index(0b011110, bitpos(6, 0))},
		{exp: 4, got: index(0b011110, bitpos(7, 0))},
		{exp: 4, got: index(0b011110, bitpos(31, 0))},

		/*15*/
		{exp: 0, got: index(0b1001010, bitpos(0, 0))},
		{exp: 0, got: index(0b1001010, bitpos(1, 0))},
		{exp: 1, got: index(0b1001010, bitpos(2, 0))},
		{exp: 1, got: index(0b1001010, bitpos(3, 0))},
		{exp: 2, got: index(0b1001010, bitpos(4, 0))},

		/*20*/
		{exp: 2, got: index(0b1001010, bitpos(5, 0))},
		{exp: 2, got: index(0b1001010, bitpos(6, 0))},
		{exp: 3, got: index(0b1001010, bitpos(7, 0))},
		{exp: 3, got: index(0b1001010, bitpos(9, 0))},
		{exp: 3, got: index(0b1001010, bitpos(240, 0))},
	}
	for i, test := range tests {
		assert.EqualInt(t, test.exp, test.got, "test #%d", i)
	}
}

func TestSize(t *testing.T) {
	t0 := &amt{}
	t1 := t0.put(0, 0, "Hello", "World!")
	assert.EqualInt(t, 32, int(unsafe.Sizeof(amt{})), "unsafe.Sizeof(amt{})")
	assert.EqualInt(t, 32, t0.size(), "t0.size()")
	assert.EqualInt(t, 16, int(unsafe.Sizeof(t1.entries[0])), "unsafe.Sizeof(t1.entries[0])")
	assert.EqualInt(t, 40, int(unsafe.Sizeof(item{})), "unsafe.Sizeof(item{})")
	assert.EqualInt(t, 32+16+40, t1.size(), "t1.size()")

	m0 := Map.WithHasher(Bole32)
	m1 := m0.Put("Hello", "World!")
	m2 := m1.Put("He11o", "There!")

	assert.EqualInt(t, 32, int(unsafe.Sizeof(Map)), "unsafe.Sizeof(Map)")
	assert.EqualInt(t, 40, m0.Size(), "m0.Size()")

	assert.EqualInt(t, 96, m1.Size(), "m1.Size()")

	assert.EqualInt(t, 2, m2.Len(), "m2.Len()")
	assert.EqualInt(t, 4, m2.Depth(), "m2.Depth()")
	assert.EqualInt(t, 8+48+48+48+64+80, m2.Size(), "m2.Size()")
}

func TestBole32(t *testing.T) {
	h := func(h uint32, k Any) uint32 {
		return h
	}
	tests := []struct{ exp, got uint32 }{
		{exp: 0, got: h(Bole32(string([]byte{})))},
		{exp: 0, got: h(Bole32(string([]byte{0})))},
		{exp: 256, got: h(Bole32(string([]byte{0, 1})))},
		{exp: 1, got: h(Bole32(string([]byte{0x01, 0x00})))},
		{exp: 0x4321, got: h(Bole32(string([]byte{0x21, 0x43})))},
		{exp: 0x654321, got: h(Bole32(string([]byte{0x21, 0x43, 0x65})))},
		{exp: 0x87654321, got: h(Bole32(string([]byte{0x21, 0x43, 0x65, 0x87})))},
		{exp: 0x87654321, got: h(Bole32(string([]byte{0x21, 0x43, 0x65, 0x87, 0x09})))},
	}
	for i, test := range tests {
		assert.Equal(t, test.exp, test.got, "test #%d", i)
	}

	hash, _ := Bole32([]byte{1, 2, 3})
	assert.EqualInt(t, 0, int(hash), "hash, _ := Bole32(); hash")
}

// Bole32 returns the head of a string as a uint32 in Little Endian Byte
// Order.
func Bole32(key Any) (uint32, Any) {
	switch k := key.(type) {
	case string:
		switch len(k) {
		case 0:
			return 0, key
		case 1:
			return uint32(k[0]), key
		case 2:
			return uint32(k[0]) | uint32(k[1])<<8, key
		case 3:
			return uint32(k[0]) | uint32(k[1])<<8 | uint32(k[2])<<16, key
		default:
			return uint32(k[0]) | uint32(k[1])<<8 | uint32(k[2])<<16 | uint32(k[3])<<24, key
		}
	default:
		return 0, key
	}
}

var assert = struct {
	Equal       func(t *testing.T, exp, got interface{}, msg string, info ...interface{})
	EqualInt    func(t *testing.T, exp, got int, msg string, info ...interface{})
	EqualString func(t *testing.T, exp, got string, msg string, info ...interface{})
}{
	Equal: func(t *testing.T, exp, got interface{}, msg string, info ...interface{}) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %#v got %#v", append(append(info, exp), got)...)
		}
	},
	EqualInt: func(t *testing.T, exp, got int, msg string, info ...interface{}) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %d got %d", append(append(info, exp), got)...)
		}
	},
}
