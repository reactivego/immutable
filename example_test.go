package immutable_test

import (
	"fmt"

	"github.com/reactivego/immutable"

	"encoding/binary"
	"hash/maphash"
)

func Example_map() {
	m := immutable.Map
	m = m.Set("Hello", "World!")
	fmt.Println(m)
	// Output:
	// Hamt{"Hello":"World!"}
}

func ExampleHamt() {
	type any = interface{}
	type Key struct{ K1, K2 int64 }
	type Topic struct{ Name, Description string }

	seed := maphash.MakeSeed()
	hasher := func(key any) (uint32, any) {
		var h maphash.Hash
		h.SetSeed(seed)
		k := key.(Key)
		binary.Write(&h, binary.LittleEndian, k.K1)
		binary.Write(&h, binary.LittleEndian, k.K2)
		return uint32(h.Sum64() & 0xFFFFFFFF), key
	}

	m := immutable.Map.WithHasher(hasher)

	m = m.Set(Key{1, 2}, Topic{"Theme", "This is a topic about theme"})
	fmt.Println(m)
	// Output:
	// HamtX{immutable_test.Key{K1:1, K2:2}:immutable_test.Topic{Name:"Theme", Description:"This is a topic about theme"}}
}

func ExampleHamtX() {
	type any = interface{}

	// Topic Name is used as the key
	type Topic struct{ Name, Description string }

	seed := maphash.MakeSeed()
	hasher := func(key any) (uint32, any) {
		var h maphash.Hash
		h.SetSeed(seed)
		h.WriteString(key.(Topic).Name)
		return uint32(h.Sum64() & 0xFFFFFFFF), key.(Topic).Name
	}

	m := immutable.Map.WithHasher(hasher)

	m = m.Put(Topic{"Theme", "This is a topic about theme"})
	fmt.Println(m)
	fmt.Printf("%#v", m.Get(Topic{Name: "Theme"}))
	// Output:
	// HamtX{"Theme":immutable_test.Topic{Name:"Theme", Description:"This is a topic about theme"}}
	// immutable_test.Topic{Name:"Theme", Description:"This is a topic about theme"}
}
