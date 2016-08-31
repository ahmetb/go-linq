# go-linq [![GoDoc](https://godoc.org/gopkg.in/ahmetalpbalkan/go-linq.v2?status.svg)](https://godoc.org/gopkg.in/ahmetalpbalkan/go-linq.v2) [![Build Status](https://travis-ci.org/ahmetalpbalkan/go-linq.svg?branch=v2.0)](https://travis-ci.org/ahmetalpbalkan/go-linq) [![Coverage Status](https://coveralls.io/repos/github/ahmetalpbalkan/go-linq/badge.svg?branch=v2.0)](https://coveralls.io/github/ahmetalpbalkan/go-linq?branch=v2.0)
A powerful language integrated query (LINQ) library for Go.
* Safe for concurrent use
* Complete lazy evaluation
* Supports arrays, slices, maps, strings, channels and
custom collections (collection needs to implement Iterable interface
and element - Comparable interface)
* Parallel LINQ (PLINQ) *(comming soon)*

## Installation

    $ go get gopkg.in/ahmetalpbalkan/go-linq.v2

## Quickstart

Usage is as easy as chaining methods like

`From(slice)` `.Where(predicate)` `.Select(selector)` `.Union(data)` 

Just keep writing.

**Example:** Find all owners of cars manufactured from 2015
```go
import . "gopkg.in/ahmetalpbalkan/go-linq.v2"
	
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
import . "gopkg.in/ahmetalpbalkan/go-linq.v2"
	
type Book struct {
	id      int
	title   string
	authors []string
}

author := From(books).SelectMany(func(b interface{}) Query {
		return From(b.(Book).authors)
	}).GroupBy(func(a interface{}) interface{} {
		return a
	}, func(a interface{}) interface{} {
		return a
	}).OrderByDescending(func(g interface{}) interface{} {
		return len(g.(Group).Group)
	}).Select(func(g interface{}) interface{} {
		return g.(Group).Key
	}).First()
```

**More examples** can be found in [documentation](https://godoc.org/gopkg.in/ahmetalpbalkan/go-linq.v2)