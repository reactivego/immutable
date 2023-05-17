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
	// {Hello:World!}
}

func ExampleHamt() {
	m := immutable.Map

	m = m.Set("first", 123).Set("second", 456).Set("first", 789)

	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})

	// Unordered Output:
	// first 789
	// second 456
}

func ExampleHamtX() {
	// Setup a hasher to allow the use of []byte values as a key
	seed := maphash.MakeSeed()
	hasher := func(key any) (uint32, any) {
		var h maphash.Hash
		h.SetSeed(seed)
		k := key.([]byte)
		h.Write(k)
		// Return a hash of the bytes and map the bytes to string to make it comparable.
		return uint32(h.Sum64() & 0xFFFFFFFF), string(k)
	}
	m := immutable.Map.WithHasher(hasher)

	m = m.Set([]byte{1, 2, 3}, "Mammalia is the class of mammals.")
	m = m.Set([]byte{4, 5, 6}, "Aves is the class of birds.")

	fmt.Println(m.Get([]byte{1, 2, 3}))
	fmt.Println(m.Get([]byte{4, 5, 6}))
	fmt.Println(m.Get([]byte{7, 8, 9}))
	// Output:
	// Mammalia is the class of mammals.
	// Aves is the class of birds.
	// <nil>
}

func ExampleHamtX_Set() {
	// Key is comparable (i.e. == and !=) but not hashable.
	type Key struct{ K1, K2 int64 }
	type Topic struct{ Name, Description string }

	// Setup a hasher to hash the key
	seed := maphash.MakeSeed()
	hasher := func(key any) (uint32, any) {
		var h maphash.Hash
		h.SetSeed(seed)
		k := key.(Key)
		binary.Write(&h, binary.LittleEndian, k.K1)
		binary.Write(&h, binary.LittleEndian, k.K2)
		// Return a hash of the key and the key itself verbatim as it is comparable
		return uint32(h.Sum64() & 0xFFFFFFFF), key
	}
	m := immutable.Map.WithHasher(hasher)

	m = m.Set(Key{1, 2}, Topic{"Theme", "This is a topic about theme"})
	fmt.Println(m)
	fmt.Println(m.Get(Key{1, 2}))
	// Output:
	// {{K1:1 K2:2}:{Name:Theme Description:This is a topic about theme}}
	// {Theme This is a topic about theme}
}

func ExampleHamtX_Put() {
	// Topic is an example of data where the key (i.e. Name) is part of the data.
	type Topic struct{ Name, Description string }

	// Setup a hasher to index a Topic on the Name field
	seed := maphash.MakeSeed()
	hasher := func(key any) (uint32, any) {
		var h maphash.Hash
		Name := key.(Topic).Name
		h.SetSeed(seed)
		h.WriteString(Name)
		// Return a 32 bit hash of Name and also the Name itself as it is comparable
		return uint32(h.Sum64() & 0xFFFFFFFF), Name
	}
	m := immutable.Map.WithHasher(hasher)

	m = m.Put(Topic{"Aves", "This topic is about birds."})
	m = m.Put(Topic{"Mammalia", "This topic is about mammals"})
	fmt.Printf("%+v\n", m.Get(Topic{Name: "Mammalia"}))
	fmt.Printf("%+v\n", m.Get(Topic{Name: "Aves"}))
	// Output:
	// {Name:Mammalia Description:This topic is about mammals}
	// {Name:Aves Description:This topic is about birds.}
}
