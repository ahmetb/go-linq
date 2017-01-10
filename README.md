# go-linq [![GoDoc](https://godoc.org/github.com/ahmetalpbalkan/go-linq?status.svg)](https://godoc.org/github.com/ahmetalpbalkan/go-linq) [![Build Status](https://travis-ci.org/ahmetalpbalkan/go-linq.svg?branch=master)](https://travis-ci.org/ahmetalpbalkan/go-linq) [![Coverage Status](https://coveralls.io/repos/github/ahmetalpbalkan/go-linq/badge.svg?branch=master)](https://coveralls.io/github/ahmetalpbalkan/go-linq?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/ahmetalpbalkan/go-linq)](https://goreportcard.com/report/github.com/ahmetalpbalkan/go-linq)
A powerful language integrated query (LINQ) library for Go.
* Written in vanilla Go, no dependencies!
* Complete lazy evaluation with iterator pattern
* Safe for concurrent use
* Supports generic functions to make your code cleaner and free of type assertions
* Supports arrays, slices, maps, strings, channels and custom collections

## Installation

    $ go get github.com/ahmetalpbalkan/go-linq

We recommend using a dependency manager (e.g. [govendor][govendor] or
[godep][godep]) to maintain a local copy of this package in your project.

[govendor]: https://github.com/kardianos/govendor
[godep]: https://github.com/tools/godep/

> :warning: :warning: `go-linq` has recently introduced _breaking API changes_
> with v2.0.0. See [release notes](#release-notes) for details. v2.0.0 comes with
> a refined interface, dramatically increased performance and memory efficiency,
> and new features such as lazy evaluation ([read more](http://kalan.rocks/2016/07/16/manipulating-data-with-iterators-in-go/)).
>
> The old version is still available in `archive/0.9` branch and tagged as `0.9`
> as well. If you are using `go-linq`, please vendor a copy of it in your
> source tree to avoid getting broken by upstream changes.

## Quickstart

Usage is as easy as chaining methods like:

`From(slice)` `.Where(predicate)` `.Select(selector)` `.Union(data)` 

**Example 1: Find all owners of cars manufactured after 2015**

```go
import . "github.com/ahmetalpbalkan/go-linq"
	
type Car struct {
    year int
    owner, model string
}

...


var owners []string

From(cars).Where(func(c interface{}) bool {
	return c.(Car).year >= 2015
}).Select(func(c interface{}) interface{} {
	return c.(Car).owner
}).ToSlice(&owners)
```

Or, you can use generic functions, like `WhereT` and `SelectT` to simplify your code
(at a performance penalty):

```go
var owners []string

From(cars).WhereT(func(c Car) bool {
	return c.year >= 2015
}).SelectT(func(c Car) string {
	return c.owner
}).ToSlice(&owners)	
```

**Example 2: Find the author who has written the most books**

```go
import . "github.com/ahmetalpbalkan/go-linq"
	
type Book struct {
	id      int
	title   string
	authors []string
}

author := From(books).SelectMany( // make a flat array of authors
	func(book interface{}) Query {
		return From(book.(Book).authors)
	}).GroupBy( // group by author
	func(author interface{}) interface{} {
		return author // author as key
	}, func(author interface{}) interface{} {
		return author // author as value
	}).OrderByDescending( // sort groups by its length
	func(group interface{}) interface{} {
		return len(group.(Group).Group)
	}).Select( // get authors out of groups
	func(group interface{}) interface{} {
		return group.(Group).Key
	}).First() // take the first author
```

**Example 3: Implement a custom method that leaves only values greater than the specified threshold**

```go
type MyQuery Query

func (q MyQuery) GreaterThan(threshold int) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if item.(int) > threshold {
						return
					}
				}

				return
			}
		},
	}
}

result := MyQuery(Range(1,10)).GreaterThan(5).Results()
```

## Generic Functions

Although Go doesn't implement generics, with some reflection tricks, you can use go-linq without
typing `interface{}`s and type assertions. This will introduce a performance penalty (5x-10x slower)
but will yield in a cleaner and more readable code.

Methods with `T` suffix (such as `WhereT`) accept functions with generic types. So instead of

    .Select(func(v interface{}) interface{} {...})

you can type:

    .SelectT(func(v YourType) YourOtherType {...})

This will make your code free of `interface{}` and type assertions.

**Example 4: "MapReduce" in a slice of string sentences to list the top 5 most used words using generic functions**

```go
var results []string

From(sentences).
	// split sentences to words
	SelectManyT(func(sentence string) Query {
		return From(strings.Split(sentence, " "))
	}).
	// group the words
	GroupByT( 
		func(word string) string { return word },
		func(word string) string { return word },
	).
	// order by count
	OrderByDescendingT(func(wordGroup Group) int {
		return len(wordGroup.Group)
	}).
	// order by the word
	ThenByT(func(wordGroup Group) string {
		return wordGroup.Key.(string)
	}).
	Take(5).  // take the top 5
	// project the words using the index as rank
	SelectIndexedT(func(index int, wordGroup Group) string {
		return fmt.Sprintf("Rank: #%d, Word: %s, Counts: %d", index+1, wordGroup.Key, len(wordGroup.Group))
	}).
	ToSlice(&results)
```

**More examples** can be found in the [documentation](https://godoc.org/github.com/ahmetalpbalkan/go-linq).

## Release Notes

~~~
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
~~~
