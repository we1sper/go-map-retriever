# go-map-retriever

A lightweight Go utility for reading values from nested map, slice, and array structures.

`go-map-retriever` is useful when you work with dynamic data such as decoded JSON, loosely typed config, or mixed `map[string]interface{}` payloads and want a cleaner way to navigate them.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/we1sper/go-map-retriever"
)

func main() {
    data := map[string]interface{}{
        "name": "alpha",
        "meta": map[string]interface{}{
            "region": "north",
            "nested": map[string]string{
                "keyA": "valA",
            },
        },
        "tags": []string{
            "red",
            "green",
            "blue",
        },
    }

    r := mapretriever.NewMapRetriever(data)

    // safe access
    name, err := r.Get("name").String()
    if err != nil {
        fmt.Printf("retrieving name failed: %v\n", err)
    } else {
        fmt.Printf("name: %s\n", name)
    }

    // unsafe access
    region := r.Get("meta", "region").Unsafe().String()
    fmt.Printf("region: %s\n", region)

    // iterate a slice
    for i, e := range r.Get("tags").Unsafe().ValueSlice() {
        fmt.Printf("tag %d: %s\n", i, e.Unsafe().String())
    }

    // string-based path (dot notation and bracket indices)
    val := r.Path("meta.nested.keyA").Unsafe().String()
    fmt.Printf("nested keyA: %s\n", val)

    firstTag := r.Path("tags[0]").Unsafe().String()
    fmt.Printf("first tag: %s\n", firstTag)
}
```

## Navigation Rules

- `Get(...)` walks through nested maps by key(s).
- `At(...)` walks through slices or arrays by index.
- `Fetch(...)` walks through a mixed path — integer arguments (including `int8`, `int16`, `int32`, `int64`) are treated as slice/array indices, everything else is treated as a map key. It is a convenience that combines `Get` and `At` so you don't have to switch between them manually.
- Negative indexes are supported in `At(...)` and `Fetch(...)`, so `At(-1)` and `Fetch(-1)` mean the last item.
- `Path(...)` accepts a string in dot/bracket notation and delegates to `Fetch`. Use dots to separate map keys and `[n]` for slice indices: `"meta.nested.keyA"` → `Fetch("meta", "nested", "keyA")`, `"tags[0]"` → `Fetch("tags", 0)`, `"grid[1][2]"` → `Fetch("grid", 1, 2)`. This is convenient when the navigation path is known as a literal string.
- `Head()` is shorthand for `At(0)`.
- `Tail()` is shorthand for `At(-1)`.


## Safe vs Unsafe Access

### Safe access

Use the normal methods when you want error handling:

```go
author, err := r.Get("author").String()
if err != nil {
    fmt.Printf("retrieving author failed: %v\n", err)
} else {
    fmt.Printf("author: %s\n", author)
}
```

### Unsafe access

Use `Unsafe()` when zero values are acceptable:

```go
age := r.Get("details", "age").Unsafe().Int()

fmt.Println(age) // 0
```

## Tracing and Debugging

When a lookup fails, `Trace()` and `Debug()` help explain where it happened.

```go
node := r.Get("details", "job", "salary")

fmt.Println(node.Trace())
fmt.Println(node.Debug())
fmt.Println(node.Error())
```

Typical output shape:

```text
source.get(details).get(job).get(salary)
source.get(details).get(job).get(salary)
        x
        |
        cannot get value for key details: key not found in map
cannot get value for key details: key not found in map
```
