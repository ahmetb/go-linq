# go-linq [![GoDoc](https://godoc.org/github.com/ahmetalpbalkan/go-linq?status.svg)](https://godoc.org/github.com/ahmetalpbalkan/go-linq) [![Build Status](https://travis-ci.org/ahmetalpbalkan/go-linq.svg?branch=master)](https://travis-ci.org/ahmetalpbalkan/go-linq) [![Coverage Status](https://coveralls.io/repos/github/ahmetalpbalkan/go-linq/badge.svg?branch=master)](https://coveralls.io/github/ahmetalpbalkan/go-linq?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/ahmetalpbalkan/go-linq)](https://goreportcard.com/report/github.com/ahmetalpbalkan/go-linq)
A powerful language integrated query (LINQ) library for Go.
* Written in vanilla Go!
* Safe for concurrent use
* Complete lazy evaluation
* Supports arrays, slices, maps, strings, channels and
custom collections (collection needs to implement Iterable interface
and element - Comparable interface)

## Installation

    $ go get github.com/ahmetalpbalkan/go-linq

## Quickstart

Usage is as easy as chaining methods like

`From(slice)` `.Where(predicate)` `.Select(selector)` `.Union(data)` 

Just keep writing.

**Example:** Find all owners of cars manufactured from 2015
```go
import . "github.com/ahmetalpbalkan/go-linq"
	
type Car struct {
    id, year int
    owner, model string
}

owners := []string{}

From(cars).Where(func(c interface{}) bool {
	return c.(Car).year >= 2015
}).Select(func(c interface{}) interface{} {
	return c.(Car).owner
}).ToSlice(&owners)
```

**Example:** Find an author that has written the most books
```go
import . "github.com/ahmetalpbalkan/go-linq"
	
type Book struct {
	id      int
	title   string
	authors []string
}

author := From(books).SelectMany( //make a flat array of authors
	func(book interface{}) Query {
		return From(book.(Book).authors)
	}).GroupBy( //group by author
	func(author interface{}) interface{} {
		return author //author as key
	}, func(author interface{}) interface{} {
		return author //author as value
	}).OrderByDescending( //sort groups by its length
	func(group interface{}) interface{} {
		return len(group.(Group).Group)
	}).Select( //get authors out of groups
	func(group interface{}) interface{} {
		return group.(Group).Key
	}).First() //take the first author
```

**More examples** can be found in [documentation](https://godoc.org/github.com/ahmetalpbalkan/go-linq)

## Release Notes
~~~

v2.0
* **Important:** This release is a breaking change
* total code rewrite with greater performance
* lots of new methods introduced

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
