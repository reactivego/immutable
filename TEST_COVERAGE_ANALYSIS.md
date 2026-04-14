# Test Coverage Analysis

**Overall: 83.2% statement coverage**

Generated via `go test -coverprofile=coverage.out && go tool cover -func=coverage.out`

## Per-Function Coverage

| File | Function | Coverage |
|------|----------|----------|
| amt.go | bitpos | 100% |
| amt.go | index | 100% |
| amt.go | present | 100% |
| amt.go | len | 100% |
| amt.go | depth | 100% |
| amt.go | lookup | 100% |
| amt.go | **foreach** | **60%** |
| amt.go | string | 100% |
| amt.go | set | 100% |
| amt.go | delete | 100% |
| error.go | Error | 100% |
| hash.go | hash | 93.8% |
| map.go | Map.Len | 100% |
| map.go | Map.Depth | 100% |
| map.go | Map.Lookup | 100% |
| map.go | Map.Has | 100% |
| map.go | Map.Get | 100% |
| map.go | Map.Range | 100% |
| map.go | Map.String | 100% |
| map.go | Map.Set | 100% |
| map.go | Map.Del | 100% |
| map.go | MapWith | 100% |
| map.go | MapX.Len | 100% |
| map.go | MapX.Depth | 100% |
| map.go | **MapX.Lookup** | **0%** |
| map.go | MapX.Has | 80% |
| map.go | MapX.Get | 80% |
| map.go | **MapX.Range** | **0%** |
| map.go | MapX.String | 100% |
| map.go | MapX.Set | 75% |
| map.go | MapX.Del | 75% |
| set.go | **Set.Len** | **0%** |
| set.go | **Set.Depth** | **0%** |
| set.go | Set.Has | 100% |
| set.go | **Set.Range** | **0%** |
| set.go | **Set.String** | **0%** |
| set.go | Set.Put | 100% |
| set.go | **Set.Del** | **0%** |
| store.go | StoreWith | 100% |
| store.go | **Store.Len** | **0%** |
| store.go | **Store.Depth** | **0%** |
| store.go | Store.Has | 100% |
| store.go | Store.Get | 100% |
| store.go | **Store.Range** | **0%** |
| store.go | **Store.String** | **0%** |
| store.go | Store.Put | 100% |
| store.go | **Store.Del** | **0%** |

## Recommended Improvements

### 1. Bug: `foreach` doesn't propagate early termination through nested nodes

In `amt.go:86-96`, when `f` returns `false` inside a recursive `foreach` call, the
outer loop continues iterating. `Range` with early termination is broken for maps
deeper than one level. A test would expose this, and the fix requires `foreach` to
propagate its cancellation signal.

### 2. Set type (5 of 7 methods at 0%)

Only `Has` and `Put` are tested. Needs tests for: `Del`, `Len`, `Depth`, `Range`,
`String`.

### 3. Store type (5 of 8 methods at 0%)

Only `Has`, `Get`, and `Put` are tested. Needs tests for: `Del`, `Len`, `Depth`,
`Range`, `String`.

### 4. MapX gaps

`MapX.Lookup` and `MapX.Range` are at 0%. Marshal-error panic paths in `Has`, `Get`,
`Set`, and `Del` are untested.

### 5. Immutability invariant

No explicit tests verify that `Set`/`Store` operations leave the original unchanged.

### 6. Hash function edge cases

The `uint` key type and non-collision `[]byte` path are untested.

### 7. Empty collection edge cases

`Del`/`Range`/`Get`/`Has`/`String` on empty collections are not tested for `Set` or
`Store`.
