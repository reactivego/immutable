package immutable

import (
	"fmt"
	"testing"
)

func TestDelDeep(t *testing.T) {
	EnableHashCollision = true
	t0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})

	k1 := "He1lo"
	v1 := "World!"

	k2 := "He2lo"
	v2 := "There!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

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
	EnableHashCollision = true
	t0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})

	k1 := "He1lo"
	v1 := "World!"

	k2 := "He2lo"
	v2 := "There!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	r2 := t2.Get(k2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 4, t2.Depth(), "t2.Depth()")
	assert.Equal(t, v2, r2, "t2.Get(k2)")
}

func TestPutDeep(t *testing.T) {
	EnableHashCollision = true
	t0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})

	k1 := "He1lo"
	v1 := "World!"

	k2 := "He2lo"
	v2 := "There!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t1.Set(k1, v2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 1, t1.Depth(), "t1.Depth()")
	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")

	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 4, t2.Depth(), "t2.Depth()")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")

	assert.EqualInt(t, 1, t3.Len(), "t3.Len()")
	assert.EqualInt(t, 1, t3.Depth(), "t3.Depth()")
	assert.Equal(t, v2, t3.Get(k1), "t3.Get(k1)")
}

func TestDelCollision(t *testing.T) {
	EnableHashCollision = true
	t0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})

	k1 := "Hello1"
	v1 := "World!"

	k2 := "Hello2"
	v2 := "There!"

	k3 := "Hello3"
	v3 := "Gophers!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k2)
	t4 := t2.Del(k1)
	t5 := t2.Set(k3, v3)
	t6 := t5.Del(k2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")

	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 8, t2.Depth(), "t2.Depth()")
	assert.Equal(t, true, t2.Has(k1), "t2.Has(k1)")
	assert.Equal(t, true, t2.Has(k2), "t2.Has(k2)")

	assert.EqualInt(t, 1, t3.Depth(), "t3.Depth()")
	assert.Equal(t, true, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, false, t3.Has(k2), "t3.Has(k2)")

	assert.EqualInt(t, 1, t4.Depth(), "t4.Depth()")
	assert.Equal(t, false, t4.Has(k1), "t4.Has(k1)")
	assert.Equal(t, true, t4.Has(k2), "t4.Has(k2)")

	assert.EqualInt(t, 8, t5.Depth(), "t5.Depth()")
	assert.Equal(t, v1, t5.Get(k1), "t5.Get(k1)")
	assert.Equal(t, v2, t5.Get(k2), "t5.Get(k2)")
	assert.Equal(t, v3, t5.Get(k3), "t5.Get(k3)")

	assert.EqualInt(t, 2, t6.Len(), "t6.Len()")
	assert.EqualInt(t, 8, t6.Depth(), "t6.Depth()")
	assert.Equal(t, v1, t6.Get(k1), "t6.Get(k1)")
	assert.Equal(t, false, t6.Has(k2), "t6.Has(k2)")
	assert.Equal(t, v3, t6.Get(k3), "t6.Get(k3)")
}

func TestGetCollision(t *testing.T) {
	EnableHashCollision = true
	t0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})

	k1 := "Hello1"
	v1 := "World!"

	k2 := "Hello2"
	v2 := "There!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	r2 := t2.Get(k2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 8, t2.Depth(), "t2.Depth()")
	assert.Equal(t, v2, r2, "t2.Get(k2)")
}

func TestPutCollision(t *testing.T) {
	EnableHashCollision = true
	t0 := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string))[:4], nil
	})

	k1 := "Hello1"
	v1 := "World!"

	k2 := "Hello2"
	v2 := "There!"

	k3 := "Hela"
	v3 := "All!"

	k4 := "Hella"
	v4 := "Strange!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Set(k3, v3)
	t4 := t3.Set(k4, v4)
	t5 := t4.Set(k1, v2)

	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 1, t1.Depth(), "t1.Depth()")

	assert.EqualInt(t, 2, t2.Len(), "t2.Len()")
	assert.EqualInt(t, 8, t2.Depth(), "t2.Depth()")

	assert.EqualInt(t, 3, t3.Len(), "t3.Len()")
	assert.EqualInt(t, 8, t3.Depth(), "t3.Depth()")

	assert.EqualInt(t, 4, t4.Len(), "t4.Len()")
	assert.EqualInt(t, 8, t4.Depth(), "t4.Depth()")

	assert.EqualInt(t, 4, t5.Len(), "t5.Len()")
	assert.EqualInt(t, 8, t5.Depth(), "t5.Depth()")
	assert.Equal(t, v2, t5.Get(k1), "t5.Get(k1)")
	assert.Equal(t, v2, t5.Get(k2), "t5.Get(k2)")
	assert.Equal(t, v3, t5.Get(k3), "t5.Get(k3)")
	assert.Equal(t, v4, t5.Get(k4), "t5.Get(k4)")
}

func TestRange(t *testing.T) {
	var t0 Map[string, string]

	k1 := "Hello"
	v1 := "World!"

	k2 := "He11o"
	v2 := "There!"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)

	c1 := 0
	t1.Range(func(key string, value string) bool {
		c1 += 1
		return key != k2
	})

	c2 := 0
	t2.Range(func(key string, value string) bool {
		c2 += 1
		return true
	})

	assert.EqualInt(t, 1, c1, "t1.Range()")
	assert.EqualInt(t, 2, c2, "t2.Range()")
}

func TestSet(t *testing.T) {
	var s0 Set[string]

	k1 := "first"
	k2 := "second"
	k3 := "third"
	k4 := "fourth"

	s1 := s0.Put(k1)
	s2 := s1.Put(k2)
	s3 := s2.Put(k3)

	assert.Equal(t, true, s3.Has(k1), "s3.Has(k1)")
	assert.Equal(t, true, s3.Has(k2), "s3.Has(k2)")
	assert.Equal(t, true, s3.Has(k3), "s3.Has(k3)")
	assert.Equal(t, false, s3.Has(k4), "s3.Has(k4)")
}

func TestStore(t *testing.T) {
	type composite struct{ name, value string }

	x0 := StoreWith(func(data composite) (string, string) {
		return data.name, data.value
	})

	k := []composite{
		{"first", "clown"},
		{"second", "joker"},
		{"third", "jester"},
		{"fourth", "fool"},
	}

	x1 := x0.Put(k[0])
	x2 := x1.Put(k[1])
	x3 := x2.Put(k[2])
	x4 := x3.Put(k[3])

	assert.Equal(t, true, x3.Has(k[0]), "x3.Has(k[0])")
	assert.Equal(t, true, x3.Has(k[1]), "x3.Has(k[1])")
	assert.Equal(t, true, x3.Has(k[2]), "x3.Has(k[2])")
	assert.Equal(t, k[2].value, x3.Get(k[2]), "x3.Get(k[2])")
	assert.Equal(t, false, x3.Has(k[3]), "x3.Has(k[3])")
	assert.Equal(t, false, x3.Has(k[3]), "x3.Has(k[3])")
	assert.Equal(t, true, x4.Has(k[3]), "x4.Has(k[3])")
	assert.Equal(t, k[3].value, x4.Get(k[3]), "x4.Has(k[3])")
}

func TestUnhashableKey(t *testing.T) {
	defer func() {
		assert.Equal(t, UnhashableKeyType, recover(), "err == UnhashableKeyType")
	}()
	assert.Equal(t, "Unhashable Key Type", UnhashableKeyType.Error(), "err == UnhashableKeyType")
	Map[any, int]{}.Set(struct{ a int }{123}, 456)
	assert.Equal(t, false, true, "Unreachable")
}

func TestPutGetDelInt(t *testing.T) {
	var t0 Map[int, string]

	k1 := uint(120120)
	v1 := "value1"

	k2 := int(120120)
	v2 := "value2"

	t1 := t0.Set(int(k1), v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(int(k1))

	count := 0
	t2.Range(func(key int, value string) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count, "t2.Range()")

	assert.Equal(t, v1, t1.Get(int(k1)), "t1.Get(k1)")
	assert.Equal(t, v1, t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v2, t2.Get(int(k1)), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(int(k1)), "t3.Has(k1)")
	assert.Equal(t, false, t3.Has(k2), "t3.Has(k2)")

	val, ok := t3.Lookup(k2)
	assert.Equal(t, false, ok, "_, ok := t3.Lookup(k2)")
	assert.Equal(t, "", val, "val, _ := t3.Lookup(k2)")
}

func TestPutGetDelUint64(t *testing.T) {
	var t0 Map[uint64, string]

	k1 := uint64(120120)
	v1 := "value1"

	k2 := uint64(240240)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelInt64(t *testing.T) {
	var t0 Map[int64, string]

	k1 := int64(120120)
	v1 := "value1"

	k2 := int64(-120120)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelUint32(t *testing.T) {
	var t0 Map[uint32, string]

	k1 := uint32(120)
	v1 := "value1"

	k2 := uint32(240)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelInt32(t *testing.T) {
	var t0 Map[int32, string]

	k1 := int32(120)
	v1 := "value1"

	k2 := int32(-120)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelUint16(t *testing.T) {
	var t0 Map[uint16, string]

	k1 := uint16(120)
	v1 := "value1"

	k2 := uint16(240)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelInt16(t *testing.T) {
	var t0 Map[int16, string]

	k1 := int16(120)
	v1 := "value1"

	k2 := int16(-120)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelUint8(t *testing.T) {
	var t0 Map[uint8, string]

	k1 := uint8(120)
	v1 := "value1"

	k2 := uint8(240)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelInt8(t *testing.T) {
	var t0 Map[int8, string]

	k1 := int8(120)
	v1 := "value1"

	k2 := int8(-120)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDelString(t *testing.T) {
	var t0 Map[string, string]

	k1 := "First Key"
	v1 := "value1"

	k2 := "Second Key"
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestPutGetDel(t *testing.T) {
	var t0 Map[string, string]

	key := "hello"
	val := "world"

	t1 := t0.Set(key, val)
	assert.EqualInt(t, 0, t0.Len(), "t0.Len()")
	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	got, ok := t1.Lookup(key)
	assert.Equal(t, true, ok, "_,ok := t1.Lookup(%q); ok", key)
	assert.Equal(t, val, got, "v,_ := t1.Lookup(%q); v", key)
	t2 := t1.Del(key)
	assert.EqualInt(t, 0, t0.Len(), "t0.Len()")
	assert.EqualInt(t, 1, t1.Len(), "t1.Len()")
	assert.EqualInt(t, 0, t2.Len(), "t2.Len()")
	assert.Equal(t, true, t1.Get(key) != "", "t1.Get(%q) != nil", key)
	assert.Equal(t, true, t2.Get(key) == "", "t2.Get(%q) == nil", key)
}

func TestStringer(t *testing.T) {
	m := Map[string, string]{}.Set("Hello", "World!")
	x := Map[string, string]{}.Set("Hello", "World!").Set("Hi", "There!")
	x = x.Del("Hello")

	assert.EqualString(t, `{Hello:World!}`, m.String(), "m.String()")
	assert.EqualString(t, `{Hi:There!}`, x.String(), "x.String()")
}

func TestPresent(t *testing.T) {
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

func TestIndex(t *testing.T) {
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

func TestRangeEarlyTermination(t *testing.T) {
	// Keys 0 and 32 share the same lower 5 bits, so they collide at
	// level 0 and create a nested node (depth > 1). Key 1 sits at a
	// separate slot in the top-level node.
	var m Map[int, string]
	m = m.Set(0, "zero").Set(1, "one").Set(32, "thirty-two")

	assert.Equal(t, true, m.Depth() > 1, "m.Depth() > 1")
	assert.EqualInt(t, 3, m.Len(), "m.Len()")

	// Early termination: callback returns false on the first call.
	// Without the foreach fix this would visit 2 entries because the
	// outer loop continued after the nested foreach returned.
	count := 0
	m.Range(func(key int, value string) bool {
		count++
		return false
	})
	assert.EqualInt(t, 1, count, "Range early termination")

	// Also verify that stopping after 2 yields exactly 2.
	count = 0
	m.Range(func(key int, value string) bool {
		count++
		return count < 2
	})
	assert.EqualInt(t, 2, count, "Range stop after 2")
}

func TestSetComplete(t *testing.T) {
	var s0 Set[string]

	// Empty set
	assert.EqualInt(t, 0, s0.Len(), "s0.Len()")
	assert.EqualInt(t, 1, s0.Depth(), "s0.Depth()")
	assert.EqualString(t, "{}", s0.String(), "s0.String()")

	c := 0
	s0.Range(func(k string) bool { c++; return true })
	assert.EqualInt(t, 0, c, "s0.Range()")

	// Del on empty is no-op
	s0d := s0.Del("nonexistent")
	assert.EqualInt(t, 0, s0d.Len(), "s0.Del(nonexistent).Len()")

	// Build set
	s1 := s0.Put("first")
	s2 := s1.Put("second")
	s3 := s2.Put("third")

	// Len
	assert.EqualInt(t, 1, s1.Len(), "s1.Len()")
	assert.EqualInt(t, 2, s2.Len(), "s2.Len()")
	assert.EqualInt(t, 3, s3.Len(), "s3.Len()")

	// Depth
	assert.Equal(t, true, s3.Depth() >= 1, "s3.Depth() >= 1")

	// Range
	count := 0
	s3.Range(func(k string) bool { count++; return true })
	assert.EqualInt(t, 3, count, "s3.Range() count")

	// String (single-entry set is deterministic)
	assert.EqualString(t, "{first}", s1.String(), "s1.String()")

	// Del
	s4 := s3.Del("second")
	assert.EqualInt(t, 2, s4.Len(), "s4.Len()")
	assert.Equal(t, true, s4.Has("first"), "s4.Has(first)")
	assert.Equal(t, false, s4.Has("second"), "s4.Has(second)")
	assert.Equal(t, true, s4.Has("third"), "s4.Has(third)")

	// Immutability: s3 unchanged after Del produced s4
	assert.EqualInt(t, 3, s3.Len(), "s3.Len() after Del")
	assert.Equal(t, true, s3.Has("second"), "s3.Has(second) after Del")
}

func TestStoreComplete(t *testing.T) {
	type composite struct{ name, value string }

	x0 := StoreWith(func(data composite) (string, string) {
		return data.name, data.value
	})

	// Empty store
	assert.EqualInt(t, 0, x0.Len(), "x0.Len()")
	assert.EqualInt(t, 1, x0.Depth(), "x0.Depth()")
	assert.EqualString(t, "{}", x0.String(), "x0.String()")

	c := 0
	x0.Range(func(k string, v string) bool { c++; return true })
	assert.EqualInt(t, 0, c, "x0.Range()")

	// Del on empty is no-op
	x0d := x0.Del(composite{"nonexistent", ""})
	assert.EqualInt(t, 0, x0d.Len(), "x0.Del(nonexistent).Len()")

	// Build store
	x1 := x0.Put(composite{"first", "clown"})
	x2 := x1.Put(composite{"second", "joker"})
	x3 := x2.Put(composite{"third", "jester"})

	// Len
	assert.EqualInt(t, 1, x1.Len(), "x1.Len()")
	assert.EqualInt(t, 2, x2.Len(), "x2.Len()")
	assert.EqualInt(t, 3, x3.Len(), "x3.Len()")

	// Depth
	assert.Equal(t, true, x3.Depth() >= 1, "x3.Depth() >= 1")

	// Range
	count := 0
	x3.Range(func(k string, v string) bool { count++; return true })
	assert.EqualInt(t, 3, count, "x3.Range() count")

	// String (single-entry store is deterministic)
	assert.EqualString(t, "{first:clown}", x1.String(), "x1.String()")

	// Del
	x4 := x3.Del(composite{"second", ""})
	assert.EqualInt(t, 2, x4.Len(), "x4.Len()")
	assert.Equal(t, true, x4.Has(composite{"first", ""}), "x4.Has(first)")
	assert.Equal(t, false, x4.Has(composite{"second", ""}), "x4.Has(second)")
	assert.Equal(t, true, x4.Has(composite{"third", ""}), "x4.Has(third)")

	// Immutability: x3 unchanged after Del produced x4
	assert.EqualInt(t, 3, x3.Len(), "x3.Len() after Del")
	assert.Equal(t, true, x3.Has(composite{"second", ""}), "x3.Has(second) after Del")
}

func TestMapXLookup(t *testing.T) {
	EnableHashCollision = false
	m := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string)), nil
	})

	m = m.Set("hello", "world")

	v, ok := m.Lookup("hello")
	assert.Equal(t, true, ok, "Lookup found")
	assert.Equal(t, "world", v, "Lookup value")

	v, ok = m.Lookup("missing")
	assert.Equal(t, false, ok, "Lookup not found")
	assert.Equal(t, "", v, "Lookup zero value")
}

func TestMapXRange(t *testing.T) {
	EnableHashCollision = false
	m := MapWith[string, string](func(a any) ([]byte, error) {
		return []byte(a.(string)), nil
	})

	m = m.Set("a", "1").Set("b", "2").Set("c", "3")

	count := 0
	m.Range(func(k string, v string) bool {
		count++
		return true
	})
	assert.EqualInt(t, 3, count, "MapX.Range() count")
}

func TestMapXMarshalError(t *testing.T) {
	errMarshal := func(a any) ([]byte, error) {
		return nil, fmt.Errorf("marshal error")
	}
	m := MapWith[string, string](errMarshal)

	// Has panic
	func() {
		defer func() { assert.Equal(t, UnhashableKeyType, recover(), "Has panic") }()
		m.Has("key")
	}()

	// Get panic
	func() {
		defer func() { assert.Equal(t, UnhashableKeyType, recover(), "Get panic") }()
		m.Get("key")
	}()

	// Set panic
	func() {
		defer func() { assert.Equal(t, UnhashableKeyType, recover(), "Set panic") }()
		m.Set("key", "val")
	}()

	// Del panic
	func() {
		defer func() { assert.Equal(t, UnhashableKeyType, recover(), "Del panic") }()
		m.Del("key")
	}()

	// Lookup panic
	func() {
		defer func() { assert.Equal(t, UnhashableKeyType, recover(), "Lookup panic") }()
		m.Lookup("key")
	}()
}

func TestPutGetDelUint(t *testing.T) {
	var t0 Map[uint, string]

	k1 := uint(120)
	v1 := "value1"

	k2 := uint(240)
	v2 := "value2"

	t1 := t0.Set(k1, v1)
	t2 := t1.Set(k2, v2)
	t3 := t2.Del(k1)

	assert.Equal(t, v1, t1.Get(k1), "t1.Get(k1)")
	assert.Equal(t, "", t1.Get(k2), "t1.Get(k2)")
	assert.Equal(t, v1, t2.Get(k1), "t2.Get(k1)")
	assert.Equal(t, v2, t2.Get(k2), "t2.Get(k2)")
	assert.Equal(t, false, t3.Has(k1), "t3.Has(k1)")
	assert.Equal(t, true, t3.Has(k2), "t3.Has(k2)")
}

func TestEmptyMapOperations(t *testing.T) {
	var m Map[string, string]

	// Del on empty
	m2 := m.Del("nonexistent")
	assert.EqualInt(t, 0, m2.Len(), "empty.Del().Len()")

	// Range on empty
	count := 0
	m.Range(func(k string, v string) bool { count++; return true })
	assert.EqualInt(t, 0, count, "empty.Range()")

	// Has/Get/Lookup on empty
	assert.Equal(t, false, m.Has("x"), "empty.Has()")
	assert.Equal(t, "", m.Get("x"), "empty.Get()")
	v, ok := m.Lookup("x")
	assert.Equal(t, false, ok, "empty.Lookup() ok")
	assert.Equal(t, "", v, "empty.Lookup() val")

	// String on empty
	assert.EqualString(t, "{}", m.String(), "empty.String()")
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
	EqualString: func(t *testing.T, exp, got string, msg string, info ...interface{}) {
		t.Helper()
		if exp != got {
			t.Errorf(msg+" expected %#q got %#q", append(append(info, exp), got)...)
		}
	},
}
