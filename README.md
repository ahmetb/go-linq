# go-linq [![travis-ci status](https://api.travis-ci.org/ahmetalpbalkan/go-linq.png)](https://travis-ci.org/ahmetalpbalkan/go-linq)  [![Codebot](https://codebot.io/badge/github.com/ahmetalpbalkan/go-linq.png)](http://codebot.io/doc/pkg/github.com/ahmetalpbalkan/go-linq "Codebot") [![GoDoc](https://godoc.org/github.com/ahmetalpbalkan/go-linq?status.png)](https://godoc.org/github.com/ahmetalpbalkan/go-linq) [![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/ahmetalpbalkan/go-linq/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

A powerful language integrated query library for Go. Querying and manipulation
operations made easy, don't repeat yourself. Inspired by Microsoft's
[LINQ](http://msdn.microsoft.com/en-us/library/bb397926.aspx).

## Installation

    $ go get github.com/ahmetalpbalkan/go-linq

then in your project 

## Quick Start

Let's find names of students over 18:

```go
import . "github.com/ahmetalpbalkan/go-linq"
	
type Student struct {
    id, age int
    name string
}

over18Names, err := From(students)
	.Where(func (s T) (bool,error){
		return s.(*Student).age >= 18, nil
	})
	.Select(func (s T) (T,error){
		return s.(*Student).name, nil
	}).Results()
```

## Documentation

* [GoDoc at codebot.io](http://codebot.io/doc/pkg/github.com/ahmetalpbalkan/go-linq)
* [GoDoc at godoc.org](http://godoc.org/github.com/ahmetalpbalkan/go-linq)

Here is wiki:

1. [Install & Import](https://github.com/ahmetalpbalkan/go-linq/wiki/Quickstart)
2. [Quickstart (Crash Course)](https://github.com/ahmetalpbalkan/go-linq/wiki/Quickstart)
3. [Parallel Linq](https://github.com/ahmetalpbalkan/go-linq/wiki/Parallel-LINQ)
4. [Table of Query Functions](https://github.com/ahmetalpbalkan/go-linq/wiki/Query-Functions)
5. [Remarks & Notes](https://github.com/ahmetalpbalkan/go-linq/wiki/Remarks-%26-notes)
6. [FAQ](https://github.com/ahmetalpbalkan/go-linq/wiki/FAQ)

## License

This software is distributed under Apache 2.0 License (see [LICENSE](LICENSE)
for more).

## Used By

Please edit [this](https://github.com/ahmetalpbalkan/go-linq/wiki/List-of-Users)
wiki page if you are using this library.

## Release Notes

~~~
v0.9-rc1
* many linq methods are implemented
* methods have error handling support
* type assertion limitations are unresolved
* travis-ci.org build integrated
* open sourced on github, master & dev branches
~~~