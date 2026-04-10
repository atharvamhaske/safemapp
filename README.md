# safemapp - thread safe concurrent map on top of native go's map.

Safemapp package is an intentionally simple attempt at a thread-safe wrapper built on top of native Go maps, using `sync.RWMutex`.

This package provides a simple API for thread-safe reads/writes with helpers like existence checks, iteration, and atomic compare-and-swap updates.

### Motivation

Go's built-in map is fast and ergonomic, but it is **not safe for concurrent access** without external synchronization. In real services, map access often spreads across multiple goroutines (request handlers, background workers, queues), and ad-hoc locking around each access quickly becomes repetitive and error-prone.

`safemapp` wraps a native map with a clear API and internal locking strategy:

- **Correctness first**: one shared lock policy for all map operations.
- **Predictable behavior**: each method is atomic at method scope.
- **Generic typing**: compile-time type safety for keys and values.
- **Low cognitive overhead**: no repeated lock/unlock code in business logic.

This package is intentionally small. It is meant for cases where you want readability and safety before reaching for more specialized concurrent structures.

### Getting Started

### Install

```bash
go get github.com/atharvamhaske/safemapp
```

### Usage

```go
package main

import (
	"fmt"

	"github.com/atharvamhaske/safemapp"
)

func main() {
	m := safemap.New[string, int]()

	m.Set("count", 1)

	v, ok := m.Get("count")
	fmt.Println(v, ok) // 1 true

	swapped := m.CompareAndSwap("count", 1, 2)
	fmt.Println(swapped) // true

	fmt.Println(m.Exists("count")) // true
	fmt.Println(m.Len())           // 1
}
```

## Generic Usage

`SafeMap` is parameterized as `SafeMap[K comparable, V comparable]`.

- `K` must be `comparable` because Go maps require comparable keys.
- `V` is currently `comparable` so `CompareAndSwap` can safely evaluate `val == old`.

Examples:

```go
// string -> int
scores := safemap.New[string, int]()

// int -> string
labels := safemap.New[int, string]()

// struct key/value (both comparable)
type Key struct{ ID int }
type Value struct{ State int }
states := safemap.New[Key, Value]()
```

## API

- `New[K, V]()` creates a new initialized map.
- `Set(k, v)` inserts or updates a value.
- `Get(k)` returns `(value, found)`.
- `Delete(k)` removes a key.
- `Exists(k)` checks key presence.
- `Len()` returns current size.
- `ForEach(func(K, V))` iterates over all entries.
- `CompareAndSwap(k, old, new)` performs atomic conditional update.

## Design Choices

- **`sync.RWMutex` over `sync.Map`**: favors explicit API, typed generics, and predictable semantics for common workloads.
- **Method-level atomicity**: each method call locks internally and completes as one synchronized operation.
- **`CompareAndSwap` built in**: enables lock-safe conditional updates without exposing internals.
- **Simple iteration model**: `ForEach` runs under read lock to avoid races while traversing.
- **Minimal surface area**: keep API compact and easy to reason about.

## Concurrency Notes

- Reads use `RLock`; writes use `Lock`.
- Operations on a single method call are thread-safe.
- `ForEach` holds a read lock while iterating; avoid slow/blocking work in the callback.

## Limitations

- `V` must be `comparable` in the current design (because of `CompareAndSwap`).
- No context-aware or cancellable operations.
- No snapshot iterator; `ForEach` holds a read lock during callback execution.
- No built-in TTL/expiration, eviction, or size limits.
- No persistence/serialization layer.
- No lock-free guarantees; high write contention can still bottleneck.

## Roadmap

- [ ] Add examples and benchmark suite.
- [ ] Add unit tests and race-detector coverage in CI.
- [ ] Consider split APIs:
  - `SafeMap[K, V any]` base map
  - optional CAS extension for `V comparable`
- [ ] Add utility helpers (`Clear`, `Keys`, `Values`, `Clone`).
- [ ] Evaluate optional snapshot iteration mode.
- [ ] Publish tagged releases and changelog.

## License

This project is licensed under the [MIT License](./LICENSE).
