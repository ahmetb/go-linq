# go-linq [![travis-ci status](https://api.travis-ci.org/ahmetalpbalkan/go-linq.png)](https://travis-ci.org/ahmetalpbalkan/go-linq) 

A powerful language integrated query library for Go. Querying and manipulation
operations made easy, don't repeat yourself. Inspired by Microsoft's
[LINQ](http://msdn.microsoft.com/en-us/library/bb397926.aspx).

## Installation

    $ go get github.com/ahmetalpbalkan/go-linq

then in your project 

## Quick Start

Let's find names of students over 18:

```
import . "github.com/ahmetalpbalkan/go-linq"
	
type Student struct {
    id, age int
    name string
}

var over18Names, err = From(students)
	.Where(func (s T) (bool,error){
		return s.(*Student).age >= 18, nil
	})
	.Select(func (s T) (T,error){
		return s.(*Student).name, nil
	}).Results()



```

## Documentation

## FAQ

## License

## Contributors

