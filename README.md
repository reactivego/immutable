# immutable

    import "github.com/reactivego/immutable"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/immutable.svg)](https://pkg.go.dev/github.com/reactivego/immutable#section-documentation)

Package `immutable` provides an immutable persistent Map type.
A Map is a collection of unordered `key:value` pairs.


```go
package main

import (
    "fmt"
    "github.com/reactivego/immutable"
)

func main() {
	var m immutable.Map[string, string]

	m = m.Set("Hello", "World!")

	fmt.Println(m)

	// Output:
	// {Hello:World!}
}
```

``` go
package main

import (
    "fmt"
    "github.com/reactivego/immutable"
)

func main() {
	var m immutable.Map[string, int]

	m = m.Set("first", 123).Set("second", 456).Set("first", 789)

	m.Range(func(key string, value int) bool {
		fmt.Println(key, value)
		return true
	})

	// Unordered Output:
	// first 789
	// second 456
}
```

``` go
package main

import (
    "fmt"
    "github.com/reactivego/immutable"
)

func main() {
	// Key is comparable (i.e. == and !=) but not hashable.
	// Notice that the key is a struct with unexported fields that are not
	// marshalled by the default json.Marshal function.
	type key struct{ A, B, c int }

	// Map with a marshal function to convert the key to []byte for hashing.
	m := immutable.MapWith[key, string](json.Marshal)

	m = m.Set(key{1, 2, 3}, "Mammalia is the class of mammals.")
	m = m.Set(key{1, 2, 4}, "Aves is the class of birds.")

	b, e := json.Marshal(key{1, 2, 3})
	fmt.Println(string(b), e)
	fmt.Println(m.Get(key{1, 2, 3}))
	fmt.Println(m.Get(key{1, 2, 4})) // same hash as key{1, 2, 3} because 'c' is not exported.
	fmt.Println(m.Get(key{7, 8, 9}))

	// Output:
	// {"A":1,"B":2} <nil>
	// Mammalia is the class of mammals.
	// Aves is the class of birds.
}
```

``` go
package main

import (
    "fmt"
    "github.com/reactivego/immutable"
)

func main() {
	// Key is comparable (i.e. == and !=) but not hashable.
	type Key struct{ K1, K2 int64 }
	type Topic struct{ Name, Description string }

	// Map with a marshal function to convert the key to []byte for hashing.
	m := immutable.MapWith[Key, Topic](json.Marshal)

	m = m.Set(Key{1, 2}, Topic{"Theme", "This is a topic about theme"})
	fmt.Println(m)
	fmt.Println(m.Get(Key{1, 2}))

	// Output:
	// {{K1:1 K2:2}:{Name:Theme Description:This is a topic about theme}}
	// {Theme This is a topic about theme}
}
```

``` go
package main

import (
    "fmt"
    "github.com/reactivego/immutable"
)

func main() {
	// Topic is an example of data where the key (i.e. Name) is part of the data.
	type Topic struct{ Name, Description string }

	store := immutable.StoreWith(func(t Topic) (string, Topic) {
		return t.Name, t
	})

	store = store.Put(Topic{Name: "Aves", Description: "This topic is about birds."})
	store = store.Put(Topic{Name: "Mammalia", Description: "This topic is about mammals"})
	fmt.Printf("%+v\n", store.Get(Topic{Name: "Mammalia"}))
	fmt.Printf("%+v\n", store.Get(Topic{Name: "Aves"}))

	// Output:
	// {Name:Mammalia Description:This topic is about mammals}
	// {Name:Aves Description:This topic is about birds.}
}
```
