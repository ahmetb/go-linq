# Migrating from go-linq v4 to v5

go-linq v5 is a complete rewrite on **Go 1.27 generic methods**
([proposal #77273](https://go.dev/issue/77273)). `Query` became the generic
`Query[T]`, every operator callback is fully typed, and the entire
reflection/type-erasure layer is gone.

The operator set and semantics are otherwise preserved: if you used the
`any`-based v4 API, most call sites migrate by deleting type assertions from
your callbacks; if you used the `…T` reflection API, your callbacks are
already typed and usually work as-is on the base method name.

## Types

| v4 | v5 |
|---|---|
| `Query` (wraps `iter.Seq[any]`) | `Query[T]` (wraps `iter.Seq[T]`) |
| `OrderedQuery` | `OrderedQuery[T]` |
| `KeyValue` (`Key`, `Value` are `any`) | `KeyValue[TKey comparable, TValue any]` |
| `Group` (`Key any`, `Group []any`) | `Group[TKey comparable, TElement any]` |
| `Iterable` (`Iterate() iter.Seq[any]`) | *removed* — pass your collection's iterator to `FromSeq(c.Iterate())` |
| `Comparable` (`CompareTo`) | *removed* — `OrderBy`/`ThenBy` keys must satisfy `cmp.Ordered`; use `Sort(less)` for custom comparison logic |

## Constructors

| v4 | v5 |
|---|---|
| `From(any)` | *removed* — use a typed constructor below |
| `FromSlice(s)` | unchanged, returns `Query[T]` |
| `FromMap(m)` | unchanged, returns `Query[KeyValue[TKey, TValue]]` |
| `FromChannel(ch)`, `FromChannelWithContext(ctx, ch)` | unchanged, return `Query[T]` |
| `FromString(s)` | unchanged, returns `Query[rune]` |
| `FromIterable(i)` | *removed* — use `FromSeq(i.Iterate())` |
| — | **new:** `FromSeq(seq iter.Seq[T])` adapts any range-over-func iterator |
| `Range(start, count)` | unchanged, returns `Query[int]` |
| `Repeat(v, count)` | unchanged, returns `Query[T]` |

## Operators

All `…T` reflection twins (`WhereT`, `SelectT`, `GroupByT`, …) are **removed**;
the base methods now accept typed functions directly and are 25–77× faster
than the `…T` API was.

| v4 | v5 |
|---|---|
| `Where(func(any) bool)` | `Where(func(T) bool)` — same for `WhereIndexed`, `TakeWhile`, `SkipWhile`, `AnyWith`, `CountWith`, `All`, `IndexOf`, … |
| `Select(func(any) any)` | `Select(func(T) TResult) Query[TResult]` — result type inferred |
| `SelectMany(func(any) Query)` | `SelectMany(func(T) Query[TResult]) Query[TResult]` — same for the `Indexed`/`By` variants |
| `Join(inner, outerKey, innerKey, result)` | `Join(inner Query[TInner], func(T) TKey, func(TInner) TKey, func(T, TInner) TResult) Query[TResult]` |
| `GroupJoin(...)` | `GroupJoin(inner Query[TInner], func(T) TKey, func(TInner) TKey, func(T, []TInner) TResult) Query[TResult]` |
| `GroupBy(keySel, elemSel)` | `GroupBy(func(T) TKey, func(T) TElement) Query[Group[TKey, TElement]]`; groups now come out in first-seen key order (was unspecified map order) |
| `Zip(q2, result)` | `Zip(q2 Query[TSecond], func(T, TSecond) TResult) Query[TResult]` |
| `OrderBy(func(any) any)` | `OrderBy(func(T) TKey)` with `TKey cmp.Ordered` — same for `OrderByDescending`, `ThenBy`, `ThenByDescending`. Each level may use a different key type. `bool` keys are no longer supported: map them to `int` |
| `Sort(less func(i, j any) bool)` | `Sort(less func(i, j T) bool)` |
| `Distinct()`, `Union(q2)`, `Except(q2)`, `Intersect(q2)`, `Contains(v)`, `SequenceEqual(q2)` | unchanged shape; elements of basic comparable kinds (integers, floats, complex, strings, booleans) are tracked and compared through strongly-typed sets with no boxing. Other element types fall back to boxed comparison at runtime (as in v4) — for those, the `By` variants with a `comparable` key remain the typed fast path |
| `DistinctBy(func(any) any)` | `DistinctBy(func(T) TKey)` with `TKey comparable` — same for `ExceptBy`, `IntersectBy`, and **new** `UnionBy` |
| `Aggregate(f)` | `Aggregate(func(acc, item T) T) (T, bool)` — `ok` replaces the nil result on empty input |
| `AggregateWithSeed(seed, f)` | `AggregateWithSeed(seed TAccumulate, func(TAccumulate, T) TAccumulate) TAccumulate` — the accumulator type may differ from `T` |
| `AggregateWithSeedBy(seed, f, resSel)` | `AggregateWithSeedBy(seed TAccumulate, func(TAccumulate, T) TAccumulate, func(TAccumulate) TResult) TResult` |

## Terminals

`nil`-returning terminals now use comma-ok:

| v4 | v5 |
|---|---|
| `First() any` | `First() (T, bool)` — same for `FirstWith`, `Last`, `LastWith`, `Single`, `SingleWith` |
| `Max() any`, `Min() any` | package functions `Max(q)`, `Min(q)` (`T cmp.Ordered`) returning `(T, bool)`; or chainable `MaxBy(func(T) TKey)`, `MinBy(func(T) TKey)` returning the *element* with the extreme key |
| `SumInts() int64`, `SumUInts() uint64`, `SumFloats() float64` | package function `Sum(q)` (result has the element type); or chainable `SumBy(func(T) TNumber)` |
| `Average() float64` | package function `Average(q)`; or chainable `AverageBy(func(T) TNumber)`. Still returns `NaN` on empty input |
| `ToSlice(&s)` | `ToSlice() []T` — no pointer argument, no capacity-reuse semantics |
| `Results() []any` | *removed* — use `ToSlice() []T` |
| `ToMap(&m)`, `ToMapBy(&m, keySel, valSel)` | package function `ToMap(q Query[KeyValue[TKey, TValue]]) map[TKey]TValue`; method `ToMapBy(func(T) TKey, func(T) TValue) map[TKey]TValue` |
| `ToChannel(chan any)` | `ToChannel(chan<- T)` |
| `ToChannelT(any)` | *removed* — `ToChannel` is typed now |

Package-level functions exist where the operation constrains the *element*
type itself (`Sum` needs numbers, `Min`/`Max` need ordering, `ToMap` needs
`KeyValue` elements) — Go methods cannot add constraints to the receiver's
type parameter. Wrap them around a chain, or stay fluent with the `…By`
variants.

## Custom query extensions

```go
// v4                              // v5
type MyQuery Query                 type MyQuery Query[int]

func (q MyQuery) GreaterThan(...)  func (q MyQuery) GreaterThan(...)
```

## Gotchas

* **Method values need instantiation:** `f := q.Select` does not compile for
  generic methods; write `f := q.Select[string]`. Ordinary chained calls are
  unaffected (inference works at call sites).
* **Construct `Query` with a named field:** write
  `Query[T]{Iterate: seq}`, not `Query[T]{seq}` — the struct carries
  additional unexported bookkeeping fields.
* **`OrderedQuery[T].Distinct()`** still deduplicates adjacent elements by
  boxed equality, matching v4.
* **`GroupBy` output order** is now deterministic (first appearance of each
  key). v4 yielded groups in random map order.
* **Empty vs nil:** `ToSlice()` on an empty query returns `nil` (v4 pointed
  the destination at an empty slice).
