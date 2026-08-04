package linq_test

import (
	"fmt"

	linq "github.com/ahmetb/go-linq/v5"
)

type exampleCar struct {
	year  int
	owner string
}

func (car exampleCar) Owner() string { return car.owner }

func ExampleQuery_Select() {
	cars := []exampleCar{{2012, "Ada"}, {2024, "Grace"}, {2020, "Linus"}}

	owners := linq.FromSlice(cars).
		Where(func(car exampleCar) bool { return car.year >= 2015 }).
		Select(exampleCar.Owner).
		OrderBy(func(owner string) string { return owner }).
		Results()

	fmt.Println(owners)
	// Output: [Grace Linus]
}

func ExampleQuery_GroupBy() {
	groups := linq.FromSlice([]string{"one", "a", "two", "bb"}).GroupBy(
		func(word string) int { return len(word) },
		func(word string) string { return word },
	).Results()

	for _, group := range groups {
		fmt.Println(group.Key, group.Group)
	}
	// Output:
	// 3 [one two]
	// 1 [a]
	// 2 [bb]
}

func ExampleQuery_Join() {
	type person struct {
		id   int
		name string
	}
	type pet struct {
		ownerID int
		name    string
	}

	people := []person{{1, "Ada"}, {2, "Grace"}}
	pets := []pet{{2, "Compiler"}, {1, "Lambda"}}

	result := linq.FromSlice(people).Join(
		linq.FromSlice(pets),
		func(value person) int { return value.id },
		func(value pet) int { return value.ownerID },
		func(owner person, animal pet) string { return owner.name + ":" + animal.name },
	).Results()

	fmt.Println(result)
	// Output: [Ada:Lambda Grace:Compiler]
}

func ExampleQuery_Iterate() {
	for value := range linq.Range(2, 3).Iterate {
		fmt.Println(value)
	}
	// Output:
	// 2
	// 3
	// 4
}
