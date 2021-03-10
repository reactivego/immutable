package immutable

import (
	// "encoding/binary"
	"fmt"
	"testing"
	"unsafe"
)

// Bole32 Byte Order Little Endian (32 bits)
func Bole32(key any) uint32 {
	return StringBOLE32(key.(string))
}

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
	assert.EqualString(t, v1, t2.Get(k1).(string), "t2.Get(k1)")
	assert.EqualString(t, v2, t2.Get(k2).(string), "t2.Get(k2)")
	assert.True(t, t3.Get(k1) == nil, "t3.Get(k1) == nil")
	assert.EqualString(t, v2, t3.Get(k2).(string), "t3.Get(k2)")
}

func TestGetDeep(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	k1 := "He1lo"
	v1 := "World!"

	k2 := "He2lo"
	v2 := "There!"

	t1 := t0.Put(k1, v1)
	t2 := t1.Put(k2, v2)
	r2 := t2.Get(k2).(string)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 4, t2.Depth(), "t2.Depth()")
	assert.EqualString(t, v2, r2, "t2.Get()")
}

func TestGetCollision(t *testing.T) {
	t0 := Map.WithHasher(Bole32)

	k1 := "Hello1"
	v1 := "World!"

	k2 := "Hello2"
	v2 := "There!"

	t1 := t0.Put(k1, v1)
	t2 := t1.Put(k2, v2)
	r2 := t2.Get(k2).(string)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 8, t2.Depth(), "t2.Depth()")
	assert.EqualString(t, v2, r2, "t2.Get()")
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
	got, ok := t1.Get(key).(string)
	assert.True(t, ok, "t1.Get() expected key %q to be present", key)
	assert.EqualString(t, val, got, "t1.Get()")
	t2 := t1.Delete(key)
	assert.EqualInt(t, 0, t0.Len(), "t0.Len()")
	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 0, t2.Len(), "t2.Len()")
	gotraw := t1.Get(key)
	assert.True(t, nil != gotraw, "t1.Get() expected key %q to be present", key)
	gotraw = t2.Get(key)
	assert.True(t, nil == gotraw, "t2.Get() expected key %q to be removed", key)
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
		assert.EqualBool(t, test.exp, test.got, fmt.Sprintf("present(n, bitpos(m)) test:%d", i))
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
		assert.EqualInt(t, test.exp, test.got, fmt.Sprintf("index() test:%d", i))
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

	assert.EqualInt(t, 40, int(unsafe.Sizeof(Map)), "unsafe.Sizeof(Map)")
	assert.EqualInt(t, 40, m0.Size(), "m0.Size()")

	assert.EqualInt(t, 96, m1.Size(), "m1.Size()")

	assert.EqualInt(t, 2, m2.Len(), "m2.Len()")
	assert.EqualInt(t, 4, m2.Depth(), "m2.Depth()")
	assert.EqualInt(t, 8+48+48+48+64+80, m2.Size(), "m2.Size()")
}

func TestStringBOLE32(t *testing.T) {
	tests := []struct{ exp, got uint32 }{
		{exp: 0, got: StringBOLE32(string([]byte{}))},
		{exp: 0, got: StringBOLE32(string([]byte{0}))},
		{exp: 256, got: StringBOLE32(string([]byte{0, 1}))},
		{exp: 1, got: StringBOLE32(string([]byte{0x01, 0x00}))},
		{exp: 0x4321, got: StringBOLE32(string([]byte{0x21, 0x43}))},
		{exp: 0x654321, got: StringBOLE32(string([]byte{0x21, 0x43, 0x65}))},
		{exp: 0x87654321, got: StringBOLE32(string([]byte{0x21, 0x43, 0x65, 0x87}))},
		{exp: 0x87654321, got: StringBOLE32(string([]byte{0x21, 0x43, 0x65, 0x87, 0x09}))},
	}
	for i, test := range tests {
		assert.EqualUint32(t, test.exp, test.got, fmt.Sprintf("StringBOLE32() test:%d", i))
	}
}

func StringBOLE32(k string) uint32 {
	switch len(k) {
	case 0:
		return 0
	case 1:
		return uint32(k[0])
	case 2:
		return uint32(k[0]) | uint32(k[1])<<8
	case 3:
		return uint32(k[0]) | uint32(k[1])<<8 | uint32(k[2])<<16
	default:
		return uint32(k[0]) | uint32(k[1])<<8 | uint32(k[2])<<16 | uint32(k[3])<<24
	}
}

var assert = struct {
	True        func(t *testing.T, correct bool, msg string, info ...interface{})
	EqualString func(t *testing.T, exp, got string, msg string)
	EqualInt    func(t *testing.T, exp, got int, msg string)
	EqualUint8  func(t *testing.T, exp, got uint8, msg string)
	EqualUint32 func(t *testing.T, exp, got uint32, msg string)
	EqualBool   func(t *testing.T, exp, got bool, msg string)
}{
	True: func(t *testing.T, correct bool, msg string, info ...interface{}) {
		t.Helper()
		if !correct {
			t.Errorf(msg, info...)
		}
	},
	EqualString: func(t *testing.T, exp, got string, msg string) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %q got %q", exp, got)
		}
	},
	EqualInt: func(t *testing.T, exp, got int, msg string) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %d got %d", exp, got)
		}
	},
	EqualUint8: func(t *testing.T, exp, got uint8, msg string) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %d got %d", exp, got)
		}
	},
	EqualUint32: func(t *testing.T, exp, got uint32, msg string) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %d got %d", exp, got)
		}
	},
	EqualBool: func(t *testing.T, exp, got bool, msg string) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %t got %t", exp, got)
		}
	},
}
