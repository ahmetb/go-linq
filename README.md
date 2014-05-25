# [go-linq][home] [![travis-ci status](https://api.travis-ci.org/ahmetalpbalkan/go-linq.png)](https://travis-ci.org/ahmetalpbalkan/go-linq) [![GoDoc](https://godoc.org/github.com/ahmetalpbalkan/go-linq?status.png)](https://godoc.org/github.com/ahmetalpbalkan/go-linq) [![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/ahmetalpbalkan/go-linq/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

[home]: http://ahmetalpbalkan.github.io/go-linq/

A powerful language integrated query library for Go. Inspired by Microsoft's
[LINQ](http://msdn.microsoft.com/en-us/library/bb397926.aspx).

* **No dependencies:** written in vanilla Go!
* **Tested:** 100.0% code coverage on all stable [releases](https://github.com/ahmetalpbalkan/go-linq/releases).
* **Backwards compatibility:** Your integration with the library will not be broken
  except major releases.

## Installation

    $ go get github.com/ahmetalpbalkan/go-linq

## Quick Start

**Example query:** Find names of students over 18:

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

* [GoDoc at godoc.org](http://godoc.org/github.com/ahmetalpbalkan/go-linq)
* [GoDoc at codebot.io](http://codebot.io/doc/pkg/github.com/ahmetalpbalkan/go-linq)

Here is wiki:

1. [Install & Import](https://github.com/ahmetalpbalkan/go-linq/wiki/Install-&-Import)
2. [Quickstart (Crash Course)](https://github.com/ahmetalpbalkan/go-linq/wiki/Quickstart)
3. [Parallel LINQ][plinq]
4. [Table of Query Functions](https://github.com/ahmetalpbalkan/go-linq/wiki/Query-Functions)
5. [Remarks & Notes](https://github.com/ahmetalpbalkan/go-linq/wiki/Remarks-%26-notes)
6. [FAQ](https://github.com/ahmetalpbalkan/go-linq/wiki/FAQ)

[plinq]: https://github.com/ahmetalpbalkan/go-linq/wiki/Parallel-LINQ-(PLINQ)

## License

This software is distributed under Apache 2.0 License (see [LICENSE](LICENSE)
for more).

## Disclaimer

As noted in LICENSE, this library is distributed on an "as is" basis and
author's employment association with Microsoft does not imply any sort of
warranty or official representation the company. This is purely a personal side
project developed on spare times.

## Authors

[Ahmet Alp Balkan](http://ahmetalpbalkan.com) – [@ahmetalpbalkan](https://twitter.com/ahmetalpbalkan)

## Release Notes

~~~

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
