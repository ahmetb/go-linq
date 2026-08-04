# go-linq v5

[![Build Status](https://github.com/ahmetb/go-linq/actions/workflows/ci.yml/badge.svg)](https://github.com/ahmetb/go-linq/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ahmetb/go-linq/v5.svg)](https://pkg.go.dev/github.com/ahmetb/go-linq/v5)

A lazy, statically typed LINQ-style query library for Go 1.27.

Version 5 uses Go 1.27 generic methods. A query keeps its element type through
filtering and ordering, and methods such as `Select`, `Join`, `GroupBy`, and
`AggregateWithSeedBy` infer their own result, key, inner, and accumulator types.
There is no reflection or `any` in production query execution.

> Go 1.27 is not yet generally available. This development branch is pinned to
> `go1.27rc2` in `go.mod` and CI.

## Installation

```sh
go get github.com/ahmetb/go-linq/v5
```

Use a Go 1.27 toolchain while the generic-method implementation is under
development:

```sh
go install golang.org/dl/go1.27rc2@latest
go1.27rc2 download
go1.27rc2 test ./...
```

## Quick start

```go
package main

import (
	"fmt"

	linq "github.com/ahmetb/go-linq/v5"
)

type Car struct {
	Year  int
	Owner string
}

func (car Car) OwnerName() string { return car.Owner }

func main() {
	cars := []Car{
		{Year: 2012, Owner: "Ada"},
		{Year: 2024, Owner: "Grace"},
		{Year: 2020, Owner: "Linus"},
	}

	owners := linq.FromSlice(cars).
		Where(func(car Car) bool { return car.Year >= 2015 }).
		Select(Car.OwnerName).
		OrderBy(func(owner string) string { return owner }).
		Results()

	fmt.Println(owners) // [Grace Linus]
}
```

`Car.OwnerName` is a method expression with type `func(Car) string`. Go 1.27
infers `Select`'s result type, so the expression changes `Query[Car]` into
`Query[string]` without a `SelectT` variant or a type assertion.

## Generic methods

The methods can introduce types independently of the receiver:

```go
type Person struct {
	ID   int
	Name string
}

type Pet struct {
	OwnerID int
	Name    string
}

func (person Person) Key() int { return person.ID }
func (pet Pet) OwnerKey() int  { return pet.OwnerID }

names := linq.FromSlice(people).Join(
	linq.FromSlice(pets),
	Person.Key,
	Pet.OwnerKey,
	func(person Person, pet Pet) string {
		return person.Name + ":" + pet.Name
	},
).Results()
```

Here `Join` infers the inner element (`Pet`), shared key (`int`), and result
(`string`) types. The same pattern powers:

- `Select`, `SelectIndexed`, `SelectMany`, and `Zip`
- `Join`, `GroupJoin`, and `GroupBy`
- `OrderBy`, `ThenBy`, and their descending variants
- `AggregateWithSeed` and `AggregateWithSeedBy`
- `ToMapBy`, `SumBy`, `AverageBy`, `MinBy`, and `MaxBy`
- keyed set and equality operations such as `DistinctBy`, `UnionBy`,
  `ExceptBy`, `IntersectBy`, `ContainsBy`, and `SequenceEqualBy`

## Constructors

All constructors preserve their source types:

```go
linq.FromSlice(values)                 // Query[T]
linq.FromMap(values)                   // Query[KeyValue[K, V]]
linq.FromChannel(ch)                   // Query[T]
linq.FromChannelWithContext(ctx, ch)   // Query[T]
linq.FromString(text)                  // Query[rune]
linq.FromIterable(collection)          // Query[T]
linq.FromSeq(sequence)                  // Query[T]
linq.Range(10, 5)                      // Query[int]
linq.Repeat(value, 5)                  // Query[T]
```

The reflection-based catch-all `From(any)` constructor was removed. Explicit
constructors give the compiler enough information to infer `T`.

## Grouping and ordering

`GroupBy` returns typed groups and preserves first-key encounter order:

```go
groups := linq.FromSlice(words).GroupBy(
	func(word string) int { return len(word) },
	func(word string) string { return word },
).Results()

for _, group := range groups {
	fmt.Println(group.Key, group.Group)
}
```

Ordering is stable. `ThenBy` and `ThenByDescending` append keys without
mutating their parent ordered query:

```go
ordered := linq.FromSlice(people).
	OrderBy(Person.LastName).
	ThenBy(Person.FirstName)
```

## Equality and set operations

Go does not allow a method to narrow `Query[T any]` to `T comparable`.
Consequently, equality-based operations use an inferred comparable key:

```go
unique := linq.FromSlice(people).DistinctBy(Person.Key)
found := unique.ContainsBy(42, Person.Key)
```

This permits non-comparable source values while making key comparability a
compile-time requirement. Set operators preserve first-occurrence order and
use true set semantics.

## Results and iteration

Materialize a query with `Results` or iterate it directly:

```go
result := linq.Range(1, 3).Results() // []int{1, 2, 3}

for value := range linq.Range(1, 3).Iterate {
	fmt.Println(value)
}
```

Queries are lazy unless an operation requires complete knowledge of the
source, such as sorting, reversing, grouping, or building a lookup for a join.
Consumer cancellation propagates through lazy pipelines.

`First`, `Last`, `Single`, and their predicate variants preserve the historical
single-return API and return the zero value of `T` when no result exists.
`Single` also returns zero when more than one result exists. Use `Any` or
`Count` when absence must be distinguished from a legitimate zero value.

## Migrating from v4

Version 5 is intentionally breaking:

- `Query` is now `Query[T]`.
- `From(any)` and reflection dispatch are removed.
- `...T` method variants are removed; the primary methods are statically typed.
- `Results` returns `[]T`; `ToMapBy` returns `map[K]V` directly.
- map queries yield `KeyValue[K,V]`, and grouping yields `Group[K,E]`.
- unconstrained equality operations use keyed `...By` forms.
- `SumBy` and `AverageBy` replace runtime numeric conversion methods.

Comparative benchmarks retain the v4 engine as test-only code. Typical generic
method pipelines reduce thousands of reflection allocations to single digits,
with improvements ranging from roughly 2× for materialization-heavy operations
to more than 50× for reflected predicate calls.

## Development

```sh
go1.27rc2 fmt ./...
go1.27rc2 vet ./...
go1.27rc2 test ./...
golangci-lint run ./...
```

`staticcheck` is temporarily disabled in `.golangci.yml` because the version
bundled with golangci-lint v2.12.2 does not terminate while analyzing exported
Go 1.27 generic methods. The remaining standard linters and `go vet` stay
enabled.
