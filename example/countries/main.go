package main

import (
	"fmt"

	"github.com/reactivego/immutable"

	"encoding/json"
	"os"
)

func main() {
	m := immutable.Map

	if cfg, err := os.Open("countries.json"); err != nil {
		panic(err)
	} else {
		defer cfg.Close()
		countries := make(map[string]string)
		if err = json.NewDecoder(cfg).Decode(&countries); err != nil {
			panic(err)
		} else {
			for k,v := range countries {
				m = m.Put(k,v)
			}
		}
	}

	fmt.Println("Len:", m.Len())
	fmt.Println("Depth:", m.Depth())
	fmt.Println("Size:", m.Size())
	fmt.Println(m)
}
