# go-linq [![GoDoc](https://godoc.org/github.com/ahmetb/go-linq?status.svg)](https://godoc.org/github.com/ahmetb/go-linq/v5) [![Build Status](https://github.com/ahmetb/go-linq/actions/workflows/ci.yml/badge.svg?branch=v5)](https://github.com/ahmetb/go-linq/actions/workflows/ci.yml) [![Coverage Status](https://coveralls.io/repos/github/ahmetb/go-linq/badge.svg?branch=v5)](https://coveralls.io/github/ahmetb/go-linq?branch=v5) [![Go Report Card](https://goreportcard.com/badge/github.com/ahmetb/go-linq)](https://goreportcard.com/report/github.com/ahmetb/go-linq)

A powerful language integrated query (LINQ) library for Go.

* **Fully type-safe:** `Query[T]` is generic, and type-changing operators like
  `Select`, `Join` and `GroupBy` are *generic methods* (a Go 1.27 language
  feature) — no `interface{}`/`any`, no type assertions, no reflection.
* **Fast:** 5–15× faster than go-linq v4 and allocation-free per element
  (see [Performance](#performance)).
* Written in vanilla Go, no dependencies!
* Complete lazy evaluation with iterator pattern based on the standard
  `iter.Seq[T]` type.
* Safe for concurrent use.
* Supports slices, maps, strings, channels, `iter.Seq[T]` iterators, and
  custom collections.

> [!NOTE]
> **Why was go-linq rewritten for Go 1.27?** Go now supports
> [generic methods](https://go.dev/issue/77273) — methods that declare their
> own type parameters, like `func (q Query[T]) Select[TResult any](...)`.
> This is the one language feature LINQ-style chaining was missing: without
> it, a method could not return a query of a *different* element type, which
> forced v1–v4 to erase everything to `any` and patch over it with runtime
> reflection. With it, the element type flows through the entire chain at
> compile time.

The same query, before and after:

<table>
<tr>
<th>go-linq v4 — type-erased</th>
<th>go-linq v5 — typed</th>
</tr>
<tr>
<td>

```go
var owners []string
FromSlice(cars).Where(func(c any) bool {
	return c.(Car).year >= 2015
}).Select(func(c any) any {
	return c.(Car).owner
}).ToSlice(&owners)
```

</td>
<td>

```go
owners := FromSlice(cars).Where(func(c Car) bool {
	return c.year >= 2015
}).Select(func(c Car) string {
	return c.owner
}).ToSlice()
```

</td>
</tr>
</table>

## Requirements

go-linq v5 requires **Go 1.27** (currently in release candidate), because it
relies on [generic methods](https://go.dev/issue/77273) for type-changing
operators. With Go 1.26 or newer installed, the toolchain listed in `go.mod`
is downloaded automatically.

For older Go versions, use [go-linq v4](https://github.com/ahmetb/go-linq/tree/master),
which offers the same operators with an `any`-based API.

## Installation

    go get github.com/ahmetb/go-linq/v5

## Quickstart

Usage is as easy as chaining methods like:

`FromSlice(slice)` `.Where(predicate)` `.Select(selector)` `.Union(data)`

Type parameters are fully inferred from your functions: you never write
explicit type arguments in a query chain.

**Example 1: Find the author who has written the most books**

```go
import . "github.com/ahmetb/go-linq/v5"

type Book struct {
	id      int
	title   string
	authors []string
}

author, ok := FromSlice(books).SelectMany( // make a flat sequence of authors
	func(book Book) Query[string] {
		return FromSlice(book.authors)
	}).GroupBy( // group by author
	func(author string) string { return author },
	func(author string) string { return author },
).MaxBy( // take the largest group
	func(group Group[string, string]) int {
		return len(group.Group)
	})
// author.Key is the author with the most books
```

**Example 2: Implement a custom method that leaves only values greater than the specified threshold**

```go
type MyQuery Query[int]

func (q MyQuery) GreaterThan(threshold int) Query[int] {
	return Query[int](q).Where(func(item int) bool {
		return item > threshold
	})
}

result := MyQuery(Range(1, 10)).GreaterThan(5).ToSlice()
```

**Example 3: "MapReduce" in a slice of string sentences to list the top 5 most used words**

```go
results := FromSlice(sentences).
	// split sentences to words
	SelectMany(func(sentence string) Query[string] {
		return FromSlice(strings.Split(sentence, " "))
	}).
	// group the words
	GroupBy(
		func(word string) string { return word },
		func(word string) string { return word },
	).
	// order by count, then by the word
	OrderByDescending(func(wordGroup Group[string, string]) int {
		return len(wordGroup.Group)
	}).
	ThenBy(func(wordGroup Group[string, string]) string {
		return wordGroup.Key
	}).
	Take(5). // take the top 5
	// project the words using the index as rank
	SelectIndexed(func(index int, wordGroup Group[string, string]) string {
		return fmt.Sprintf("Rank: #%d, Word: %s, Counts: %d",
			index+1, wordGroup.Key, len(wordGroup.Group))
	}).
	ToSlice()
```

## Manual Iteration

The `Query[T]` type exposes an `Iterate` field of type `iter.Seq[T]`, which
integrates with Go's native iteration style and the `iter`/`slices` standard
library packages.

**Example 4: Iterate over a query using the standard `for ... range` loop**

```go
q := FromSlice([]int{1, 2, 3, 4})

for v := range q.Iterate {
	fmt.Println(v)
}
```

## Data Source Constructors

Each constructor is typed for its specific input, and the element type of the
resulting query is inferred from the argument:

- `FromSlice` — creates a query from a slice.
- `FromMap` — creates a `Query[KeyValue[TKey, TValue]]` from a map.
- `FromChannel` — creates a query from a channel.
- `FromChannelWithContext` — creates a query from a channel with `Context` support.
- `FromString` — creates a `Query[rune]` from a string.
- `FromSeq` — creates a query from any standard `iter.Seq[T]` iterator,
  including custom collections that expose an iterator method.
- `Range`, `Repeat` — generate sequences.

The runtime-reflection based `From(any)` constructor from v4 has been removed:
in a fully-typed API the element type must be known at the call site.

## Performance

v5 eliminates the three taxes the type-erased v4 API paid on every element:
interface boxing, type assertions, and reflection. Per-element work in a v5
chain is just typed closure calls; the only allocations are the fixed closure
captures made when the query is constructed.

Measured on Apple M5 Pro with go1.27, 1M-element `[]int` (100k structs for
the projection case):

| Benchmark | v4 (`any` API) | v4 (`…T` reflection API) | **v5** | hand-written loop |
|---|---|---|---|---|
| `Where` → `Sum` | 23.0 ms / 999,754 allocs | 103.9 ms / 3.0M allocs | **1.5 ms / 3 allocs** | 0.45 ms / 0 allocs |
| `Where` → `ToSlice` | 16.4 ms / 1.5M allocs | — | **2.5 ms / 38 allocs** | (result slice only) |
| struct `Where` → `Select` → `ToSlice` (100k) | 3.33 ms / 226,685 allocs | 18.0 ms / 490k allocs | **0.71 ms / 30 allocs** | 26 allocs |

In short: **5–15× faster** than idiomatic v4, **25–77× faster** than the v4
`…T` reflection API, with allocations dropping from *O(n)* to *O(1)* per
query. The residual gap to a hand-written loop is the per-element closure
call, inherent to any lazy iterator.

## Migrating from v4

See [MIGRATION.md](MIGRATION.md) for a complete v4 → v5 symbol table. The
highlights:

* `Query` is now `Query[T]`; all operators take typed functions
  (`func(T) bool` instead of `func(any) bool`).
* All `…T` reflection twins (`WhereT`, `SelectT`, …) are gone — the base
  methods are now just as clean and much faster.
* Element-returning terminals (`First`, `Last`, `Single`, `Aggregate`,
  `Min`, `Max`, …) return `(T, bool)` instead of a nil-able `any`.
* `Min`, `Max`, `Sum`, `Average`, `ToMap` are package-level functions
  (their constraints depend on the element type); chainable `MinBy`,
  `MaxBy`, `SumBy`, `AverageBy`, `ToMapBy` methods are available.
* `ToSlice()` returns `[]T` instead of filling a pointer argument.

## Release Notes

```text
v5.0.0 (unreleased)
* Breaking change: COMPLETE REWRITE on Go 1.27 generic methods.
  - Query is now the generic Query[T]; operator callbacks are fully typed.
  - Type-changing operators (Select, SelectMany, Join, GroupJoin, GroupBy,
    Zip, AggregateWithSeed, ToMapBy, ...) are generic methods whose result
    types are inferred from the supplied functions.
  - Removed the entire reflection layer and all xxxT methods.
  - Removed From(any); use the typed From* constructors.
  - First/FirstWith/Last/LastWith/Single/SingleWith/Aggregate/Min/Max
    return (T, bool) instead of nil-able any.
  - Min/Max/Sum/Average/ToMap are package-level functions; added chainable
    MinBy/MaxBy/SumBy/AverageBy/ToMapBy/UnionBy methods.
  - ToSlice() returns []T; ToMap/ToMapBy return maps.
  - Removed Comparable interface; OrderBy/ThenBy keys must satisfy
    cmp.Ordered (use Sort for custom comparisons).
  - GroupBy yields groups in first-seen key order (deterministic).
  - Added FromSeq to adapt any iter.Seq[T].
  - 5-15x faster than v4; allocations drop from O(n) to O(1) per query.

v4.0.0 (2025-10-12)
* Breaking change: Migrated to standard Go iterator pattern. (thanks @kalaninja!)
* Added typed constructors: FromSlice(), FromMap(), FromChannel(),
 FromChannelWithContext(), FromString().
* Breaking change: Removed FromChannelT() in favor of FromChannel().

v3.2.0 (2020-12-29)
* Added FromChannelT().
* Added DefaultIfEmpty().

v3.1.0 (2019-07-09)
* Support for Go modules
* Added IndexOf()/IndexOfT().

v3.0.0 (2017-01-10)
* Breaking change: ToSlice() now overwrites existing slice starting
  from index 0 and grows/reslices it as needed.
* Generic methods support (thanks @cleitonmarx!)
  - Accepting parametrized functions was originally proposed in #26
  - You can now avoid type assertions and interface{}s
  - Functions with generic methods are named as "MethodNameT" and
    signature for the existing LINQ methods are unchanged.
* Added ForEach(), ForEachIndexed() and AggregateWithSeedBy().

v2.0.0 (2016-09-02)
* IMPORTANT: This release is a BREAKING CHANGE. The old version
  is archived at the 'archive/0.9' branch or the 0.9 tags.
* A COMPLETE REWRITE of go-linq with better performance and memory
  efficiency. (thanks @kalaninja!)
* API has significantly changed. Most notably:
  - linq.T removed in favor of interface{}
  - library methods no longer return errors
  - PLINQ removed for now (see channels support)
  - support for channels, custom collections and comparables

v0.9-rc4
* GroupBy()

v0.9-rc3.2
* bugfix: All() iterating over values instead of indices

v0.9-rc3.1
* bugfix: modifying result slice affects subsequent query methods

v0.9-rc3
* removed FirstOrNil, LastOrNil, ElementAtOrNil methods

v0.9-rc2.5
* slice-accepting methods accept slices of any type with reflections

v0.9-rc2
* parallel linq (plinq) implemented
* Queryable separated into Query & ParallelQuery
* fixed early termination for All

v0.9-rc1
* many linq methods are implemented
* methods have error handling support
* type assertion limitations are unresolved
* travis-ci.org build integrated
* open sourced on github, master & dev branches
```
