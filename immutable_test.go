package immutable

import (
	// "encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/reactivego/immutable/byteorder"
)

// TestKey implements a "hash" that can easily be forced into collisions.
type TestKey string

// Sum32 returns a hash value for the key, in this case just the head of the
// string little endian encoded.
func (k TestKey) Sum32() uint32 {
	return byteorder.LittleEndian.Uint32([]byte(k))
}

func TestDelDeep(t *testing.T) {
	t0 := NewMap()

	k1 := TestKey("He1lo")
	v1 := "World!"

	k2 := TestKey("He2lo")
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
	t0 := NewMap()

	k1 := TestKey("He1lo")
	v1 := "World!"

	k2 := TestKey("He2lo")
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
	t0 := NewMap()

	k1 := TestKey("Hello1")
	v1 := "World!"

	k2 := TestKey("Hello2")
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
	t0 := NewMap()

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
	t0 := NewMap()

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
		{exp: true, got: present(0b0001, bitpos(0))},
		{exp: true, got: present(0b0010, bitpos(1))},
		{exp: true, got: present(0b0100, bitpos(2))},
		{exp: true, got: present(0b1000, bitpos(3))},
		{exp: false, got: present(0b0010, bitpos(0))},

		/*5*/
		{exp: false, got: present(0b0100, bitpos(0))},
		{exp: true, got: present(0b0100, bitpos(2))},
		{exp: false, got: present(0b0100, bitpos(3))},
		{exp: true, got: present(0b1100, bitpos(3))},
		{exp: false, got: present(0b10100, bitpos(3))},

		/*10*/
		{exp: true, got: present(0b11100, bitpos(4))},
		{exp: true, got: present(0b111100, bitpos(4))},
		{exp: true, got: present(0b111100, bitpos(5))},
		{exp: true, got: present(0b100100, bitpos(2))},
		{exp: true, got: present(0b100100, bitpos(5))},

		/*15*/
		{exp: false, got: present(0b1001010, bitpos(0))},
		{exp: true, got: present(0b1001010, bitpos(1))},
		{exp: false, got: present(0b1001010, bitpos(2))},
		{exp: true, got: present(0b1001010, bitpos(3))},
		{exp: false, got: present(0b1001010, bitpos(4))},

		/*20*/
		{exp: false, got: present(0b1001010, bitpos(5))},
		{exp: true, got: present(0b1001010, bitpos(6))},
		{exp: false, got: present(0b1001010, bitpos(7))},
		{exp: false, got: present(0b1001010, bitpos(7))},
		{exp: false, got: present(0b1001010, bitpos(240))},
	}
	for i, test := range tests {
		assert.EqualBool(t, test.exp, test.got, fmt.Sprintf("present(n, bitpos(m)) test:%d", i))
	}
}

func TestEntryIndex(t *testing.T) {
	tests := []struct{ exp, got int }{
		/*0*/
		{exp: 0, got: index(0b0000, bitpos(0))},
		{exp: 0, got: index(0b0000, bitpos(1))},
		{exp: 0, got: index(0b0001, bitpos(0))},
		{exp: 1, got: index(0b0001, bitpos(1))},
		{exp: 1, got: index(0b0001, bitpos(2))},

		/*5*/
		{exp: 0, got: index(0b0100, bitpos(0))},
		{exp: 0, got: index(0b0100, bitpos(1))},
		{exp: 0, got: index(0b0100, bitpos(2))},
		{exp: 1, got: index(0b0100, bitpos(3))},
		{exp: 1, got: index(0b0100, bitpos(4))},

		/*10*/
		{exp: 3, got: index(0b011110, bitpos(4))},
		{exp: 4, got: index(0b011110, bitpos(5))},
		{exp: 4, got: index(0b011110, bitpos(6))},
		{exp: 4, got: index(0b011110, bitpos(7))},
		{exp: 4, got: index(0b011110, bitpos(31))},

		/*15*/
		{exp: 0, got: index(0b1001010, bitpos(0))},
		{exp: 0, got: index(0b1001010, bitpos(1))},
		{exp: 1, got: index(0b1001010, bitpos(2))},
		{exp: 1, got: index(0b1001010, bitpos(3))},
		{exp: 2, got: index(0b1001010, bitpos(4))},

		/*20*/
		{exp: 2, got: index(0b1001010, bitpos(5))},
		{exp: 2, got: index(0b1001010, bitpos(6))},
		{exp: 3, got: index(0b1001010, bitpos(7))},
		{exp: 3, got: index(0b1001010, bitpos(9))},
		{exp: 3, got: index(0b1001010, bitpos(240))},
	}
	for i, test := range tests {
		assert.EqualInt(t, test.exp, test.got, fmt.Sprintf("index() test:%d", i))
	}
}

func TestMask(t *testing.T) {
	tests := []struct{ exp, got uint8 }{
		{exp: 1, got: mask(0x1, 5*0)},
		{exp: 16, got: mask(0x10, 5*0)},
		{exp: 0, got: mask(0x1, 5*1)},
		{exp: 1, got: mask(0x20, 5*1)},
		{exp: 1, got: mask(0x40000000, 5*6)},
		{exp: 2, got: mask(0x80000000, 5*6)},
		{exp: 3, got: mask(0xC0000000, 5*6)},
	}
	for i, test := range tests {
		assert.EqualUint8(t, test.exp, test.got, fmt.Sprintf("mask() test:%d", i))
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

	m0 := NewMap()
	m1 := m0.Put("Hello", "World!")
	m2 := m1.Put("He11o", "There!")

	assert.EqualInt(t, 8, int(unsafe.Sizeof(Map{})), "unsafe.Sizeof(Map{})")
	assert.EqualInt(t, 40, m0.Size(), "m0.Size()")

	assert.EqualInt(t, 96, m1.Size(), "m1.Size()")

	assert.EqualInt(t, 2, m2.Len(), "m2.Len()")
	assert.EqualInt(t, 4, m2.Depth(), "m2.Depth()")
	assert.EqualInt(t, 8+48+48+48+64+80, m2.Size(), "m2.Size()")
}

func TestLittleEndianUint32(t *testing.T) {
	tests := []struct{ exp, got uint32 }{
		{exp: 0, got: byteorder.LittleEndian.Uint32([]byte{})},
		{exp: 0, got: byteorder.LittleEndian.Uint32([]byte{0})},
		{exp: 256, got: byteorder.LittleEndian.Uint32([]byte{0, 1})},
		{exp: 1, got: byteorder.LittleEndian.Uint32([]byte{0x01, 0x00})},
		{exp: 0x4321, got: byteorder.LittleEndian.Uint32([]byte{0x21, 0x43})},
		{exp: 0x654321, got: byteorder.LittleEndian.Uint32([]byte{0x21, 0x43, 0x65})},
		{exp: 0x87654321, got: byteorder.LittleEndian.Uint32([]byte{0x21, 0x43, 0x65, 0x87})},
		{exp: 0x87654321, got: byteorder.LittleEndian.Uint32([]byte{0x21, 0x43, 0x65, 0x87, 0x09})},
	}
	for i, test := range tests {
		assert.EqualUint32(t, test.exp, test.got, fmt.Sprintf("LittleEndian.Uint32() test:%d", i))
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
