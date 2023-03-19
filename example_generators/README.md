Example generators
==================

Toy examples showing implementation of generators using `gopkg`.

Each generator contains their own examples which can be used to tryout each generator; e.g.
```
cd enumstringergen/example_single_enum
go generate
```

* `enumstringergen` - Implements the classic [stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer) generator (with the addition of generating tests for the implementation).

* `interfacegen` - Given an interface, stamps out an empty `struct` which implements that interface with stubbed methods (along with tests).
