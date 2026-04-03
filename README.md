# go-map-retriever

A lightweight Go utility for reading values from nested map, slice, and array structures.

`go-map-retriever` is useful when you work with dynamic data such as decoded JSON, loosely typed config, or mixed `map[string]interface{}` payloads and want a cleaner way to navigate them.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/we1sper/go-map-retriever/mapretriever"
)

func main() {
    data := map[string]interface{}{
        "author":  "welsper",
        "country": "China",
        "emails": []string{
            "welsper@qq.com",
            "welsper@nit.com",
            "welsper@nekomi.com",
        },
    }

    r := mapretriever.NewMapRetriever(data)

    author, err := r.Get("author").String()
    if err != nil {
        fmt.Printf("retrieving author failed: %v\n", err)
    } else {
        fmt.Printf("author: %s\n", author)
    }

    country := r.Get("country").Unsafe().String()
    fmt.Printf("country: %s\n", country)

    for i, e := range r.Get("emails").Unsafe().ValueSlice() {
        fmt.Printf("email %d: %s\n", i, e.Unsafe().String())
    }
}
```

## Navigation Rules

- `Get(...)` walks through nested maps one key at a time.
- `At(...)` walks through slices or arrays by index.
- Negative indexes are supported in `At(...)`, so `At(-1)` means the last item.
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
source.get(details).get(job).get(salary).
        x
        |
        cannot get value for key details: key not found in map
cannot get value for key salary: parent is nil
```