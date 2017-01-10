package linq

import (
	"fmt"
	"strings"
	"time"
)

func ExampleKeyValue() {
	m := make(map[int]bool)
	m[10] = true

	fmt.Println(From(m).Results())

	// Output:
	// [{10 true}]
}

func ExampleKeyValue_second() {
	input := []KeyValue{
		{10, true},
	}

	m := make(map[int]bool)
	From(input).
		ToMap(&m)

	fmt.Println(m)

	// Output:
	// map[10:true]
}

// The following code example demonstrates how
// to use Range to generate a slice of values.
func ExampleRange() {
	// Generate a slice of integers from 1 to 10
	// and then select their squares.
	var squares []int
	Range(1, 10).
		SelectT(
			func(x int) int { return x * x },
		).
		ToSlice(&squares)

	for _, num := range squares {
		fmt.Println(num)
	}
	//Output:
	//1
	//4
	//9
	//16
	//25
	//36
	//49
	//64
	//81
	//100
}

// The following code example demonstrates how to use Repeat
// to generate a slice of a repeated value.
func ExampleRepeat() {
	var slice []string
	Repeat("I like programming.", 5).
		ToSlice(&slice)

	for _, str := range slice {
		fmt.Println(str)
	}
	//Output:
	//I like programming.
	//I like programming.
	//I like programming.
	//I like programming.
	//I like programming.

}

func ExampleQuery() {
	query := From([]int{1, 2, 3, 4, 5}).Where(func(i interface{}) bool {
		return i.(int) <= 3
	})

	next := query.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		fmt.Println(item)
	}

	// Output:
	// 1
	// 2
	// 3
}

// The following code example demonstrates how to use Aggregate function
func ExampleQuery_Aggregate() {
	fruits := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	// Determine which string in the slice is the longest.
	longestName := From(fruits).
		Aggregate(
			func(r interface{}, i interface{}) interface{} {
				if len(r.(string)) > len(i.(string)) {
					return r
				}
				return i
			},
		)

	fmt.Println(longestName)

	// Output:
	// passionfruit
}

// The following code example demonstrates how to use AggregateWithSeed function
func ExampleQuery_AggregateWithSeed() {
	ints := []int{4, 8, 8, 3, 9, 0, 7, 8, 2}

	// Count the even numbers in the array, using a seed value of 0.
	numEven := From(ints).
		AggregateWithSeed(0,
			func(total, next interface{}) interface{} {
				if next.(int)%2 == 0 {
					return total.(int) + 1
				}
				return total
			},
		)

	fmt.Printf("The number of even integers is: %d", numEven)
	// Output:
	// The number of even integers is: 6
}

// The following code example demonstrates how to use AggregateWithSeedBy function
func ExampleQuery_AggregateWithSeedBy() {
	input := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	// Determine whether any string in the array is longer than "banana".
	longestName := From(input).
		AggregateWithSeedBy("banana",
			func(longest interface{}, next interface{}) interface{} {
				if len(longest.(string)) > len(next.(string)) {
					return longest
				}
				return next

			},
			// Return the final result
			func(result interface{}) interface{} {
				return fmt.Sprintf("The fruit with the longest name is %s.", result)
			},
		)

	fmt.Println(longestName)
	// Output:
	// The fruit with the longest name is passionfruit.
}

// The following code example demonstrates how to
// use Distinct to return distinct elements from a slice of integers.
func ExampleOrderedQuery_Distinct() {
	ages := []int{21, 46, 46, 55, 17, 21, 55, 55}

	var distinctAges []int
	From(ages).
		OrderBy(
			func(item interface{}) interface{} { return item },
		).
		Distinct().
		ToSlice(&distinctAges)

	fmt.Println(distinctAges)
	// Output:
	// [17 21 46 55]
}

// The following code example demonstrates how to
// use DistinctBy to return distinct elements from a ordered slice of elements.
func ExampleOrderedQuery_DistinctBy() {
	type Product struct {
		Name string
		Code int
	}

	products := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	//Order and exclude duplicates.
	var noduplicates []Product
	From(products).
		OrderBy(
			func(item interface{}) interface{} { return item.(Product).Name },
		).
		DistinctBy(
			func(item interface{}) interface{} { return item.(Product).Code },
		).
		ToSlice(&noduplicates)

	for _, product := range noduplicates {
		fmt.Printf("%s %d\n", product.Name, product.Code)
	}

	// Output:
	// apple 9
	// lemon 12
	// orange 4
}

// The following code example demonstrates how to use ThenBy to perform
// a secondary ordering of the elements in a slice.
func ExampleOrderedQuery_ThenBy() {
	fruits := []string{"grape", "passionfruit", "banana", "mango", "orange", "raspberry", "apple", "blueberry"}

	// Sort the strings first by their length and then
	//alphabetically by passing the identity selector function.
	var query []string
	From(fruits).
		OrderBy(
			func(fruit interface{}) interface{} { return len(fruit.(string)) },
		).
		ThenBy(
			func(fruit interface{}) interface{} { return fruit },
		).
		ToSlice(&query)

	for _, fruit := range query {
		fmt.Println(fruit)
	}

	// Output:
	// apple
	// grape
	// mango
	// banana
	// orange
	// blueberry
	// raspberry
	// passionfruit
}

// The following code example demonstrates how to use All to determine
// whether all the elements in a slice satisfy a condition.
// Variable allStartWithB is true if all the pet names start with "B"
// or if the pets array is empty.
func ExampleQuery_All() {

	type Pet struct {
		Name string
		Age  int
	}

	pets := []Pet{
		{Name: "Barley", Age: 10},
		{Name: "Boots", Age: 4},
		{Name: "Whiskers", Age: 6},
	}

	// Determine whether all pet names
	// in the array start with 'B'.
	allStartWithB := From(pets).
		All(
			func(pet interface{}) bool { return strings.HasPrefix(pet.(Pet).Name, "B") },
		)

	fmt.Printf("All pet names start with 'B'? %t", allStartWithB)

	// Output:
	//
	//  All pet names start with 'B'? false
}

// The following code example demonstrates how to use Any to determine
// whether a slice contains any elements.
func ExampleQuery_Any() {

	numbers := []int{1, 2}
	hasElements := From(numbers).Any()

	fmt.Printf("Are there any element in the list? %t", hasElements)

	// Output:
	// Are there any element in the list? true
}

// The following code example demonstrates how to use AnyWith
// to determine whether any element in a slice satisfies a condition.
func ExampleQuery_AnyWith() {

	type Pet struct {
		Name       string
		Age        int
		Vaccinated bool
	}

	pets := []Pet{
		{Name: "Barley", Age: 8, Vaccinated: true},
		{Name: "Boots", Age: 4, Vaccinated: false},
		{Name: "Whiskers", Age: 1, Vaccinated: false},
	}

	// Determine whether any pets over age 1 are also unvaccinated.
	unvaccinated := From(pets).
		AnyWith(
			func(p interface{}) bool {
				return p.(Pet).Age > 1 && p.(Pet).Vaccinated == false
			},
		)

	fmt.Printf("Are there any unvaccinated animals over age one? %t", unvaccinated)

	// Output:
	//
	// Are there any unvaccinated animals over age one? true
}

// The following code example demonstrates how to use Append
// to include an elements in the last position of a slice.
func ExampleQuery_Append() {
	input := []int{1, 2, 3, 4}

	q := From(input).Append(5)

	last := q.Last()

	fmt.Println(last)

	// Output:
	// 5
}

//The following code example demonstrates how to use Average
//to calculate the average of a slice of values.
func ExampleQuery_Average() {
	grades := []int{78, 92, 100, 37, 81}
	average := From(grades).Average()

	fmt.Println(average)

	// Output:
	// 77.6
}

// The following code example demonstrates how to use Count
// to count the elements in an array.
func ExampleQuery_Count() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}
	numberOfFruits := From(fruits).Count()

	fmt.Println(numberOfFruits)

	// Output:
	// 6
}

// The following code example demonstrates how to use Contains
// to determine whether a slice contains a specific element.
func ExampleQuery_Contains() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	has5 := From(slice).Contains(5)

	fmt.Printf("Does the slice contains 5? %t", has5)

	// Output:
	// Does the slice contains 5? true
}

//The following code example demonstrates how to use CountWith
//to count the even numbers in an array.
func ExampleQuery_CountWith() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	evenCount := From(slice).
		CountWith(
			func(item interface{}) bool { return item.(int)%2 == 0 },
		)

	fmt.Println(evenCount)

	// Output:
	// 6
}

//The following code example demonstrates how to use Distinct
//to return distinct elements from a slice of integers.
func ExampleQuery_Distinct() {
	ages := []int{21, 46, 46, 55, 17, 21, 55, 55}

	var distinctAges []int
	From(ages).
		Distinct().
		ToSlice(&distinctAges)

	fmt.Println(distinctAges)

	// Output:
	// [21 46 55 17]
}

// The following code example demonstrates how to
// use DistinctBy to return distinct elements from a ordered slice of elements.
func ExampleQuery_DistinctBy() {
	type Product struct {
		Name string
		Code int
	}

	products := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	//Order and exclude duplicates.
	var noduplicates []Product
	From(products).
		DistinctBy(
			func(item interface{}) interface{} { return item.(Product).Code },
		).
		ToSlice(&noduplicates)

	for _, product := range noduplicates {
		fmt.Printf("%s %d\n", product.Name, product.Code)
	}

	// Output:
	// orange 4
	// apple 9
	// lemon 12

}

// The following code example demonstrates how to use the Except
// method to compare two slices of numbers and return elements
// that appear only in the first slice.
func ExampleQuery_Except() {
	numbers1 := []float32{2.0, 2.1, 2.2, 2.3, 2.4, 2.5}
	numbers2 := []float32{2.2}

	var onlyInFirstSet []float32
	From(numbers1).
		Except(From(numbers2)).
		ToSlice(&onlyInFirstSet)

	for _, number := range onlyInFirstSet {
		fmt.Println(number)
	}

	// Output:
	//2
	//2.1
	//2.3
	//2.4
	//2.5

}

// The following code example demonstrates how to use the Except
// method to compare two slices of numbers and return elements
// that appear only in the first slice.
func ExampleQuery_ExceptBy() {
	type Product struct {
		Name string
		Code int
	}

	fruits1 := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	fruits2 := []Product{
		{Name: "apple", Code: 9},
	}

	//Order and exclude duplicates.
	var except []Product
	From(fruits1).
		ExceptBy(From(fruits2),
			func(item interface{}) interface{} { return item.(Product).Code },
		).
		ToSlice(&except)

	for _, product := range except {
		fmt.Printf("%s %d\n", product.Name, product.Code)
	}

	// Output:
	// orange 4
	// lemon 12

}

// The following code example demonstrates how to use First
// to return the first element of an array.
func ExampleQuery_First() {
	numbers := []int{9, 34, 65, 92, 87, 435, 3, 54, 83, 23, 87, 435, 67, 12, 19}

	first := From(numbers).First()

	fmt.Println(first)

	// Output:
	// 9

}

//The following code example demonstrates how to use FirstWith
// to return the first element of an array that satisfies a condition.
func ExampleQuery_FirstWith() {
	numbers := []int{9, 34, 65, 92, 87, 435, 3, 54, 83, 23, 87, 435, 67, 12, 19}

	first := From(numbers).
		FirstWith(
			func(item interface{}) bool { return item.(int) > 80 },
		)

	fmt.Println(first)

	// Output:
	// 92

}

//The following code example demonstrates how to use Intersect
//to return the elements that appear in each of two slices of integers.
func ExampleQuery_Intersect() {
	id1 := []int{44, 26, 92, 30, 71, 38}
	id2 := []int{39, 59, 83, 47, 26, 4, 30}

	var both []int
	From(id1).
		Intersect(From(id2)).
		ToSlice(&both)

	for _, id := range both {
		fmt.Println(id)
	}

	// Output:
	// 26
	// 30

}

//The following code example demonstrates how to use IntersectBy
//to return the elements that appear in each of two slices of products with same Code.
func ExampleQuery_IntersectBy() {
	type Product struct {
		Name string
		Code int
	}

	store1 := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
	}

	store2 := []Product{
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	var duplicates []Product
	From(store1).
		IntersectBy(From(store2),
			func(p interface{}) interface{} { return p.(Product).Code },
		).
		ToSlice(&duplicates)

	for _, p := range duplicates {
		fmt.Println(p.Name, "", p.Code)
	}

	// Output:
	// apple  9

}

// The following code example demonstrates how to use Last
// to return the last element of an array.
func ExampleQuery_Last() {
	numbers := []int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 67, 12, 19}

	last := From(numbers).Last()

	fmt.Println(last)

	//Output:
	//19

}

// The following code example demonstrates how to use LastWith
// to return the last element of an array.
func ExampleQuery_LastWith() {
	numbers := []int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 67, 12, 19}

	last := From(numbers).
		LastWith(
			func(n interface{}) bool { return n.(int) > 80 },
		)

	fmt.Println(last)

	//Output:
	//87

}

// The following code example demonstrates how to use Max
// to determine the maximum value in a slice.
func ExampleQuery_Max() {
	numbers := []int64{4294967296, 466855135, 81125}

	last := From(numbers).Max()

	fmt.Println(last)

	//Output:
	//4294967296

}

// The following code example demonstrates how to use Min
// to determine the minimum value in a slice.
func ExampleQuery_Min() {
	grades := []int{78, 92, 99, 37, 81}

	min := From(grades).Min()

	fmt.Println(min)

	//Output:
	//37

}

// The following code example demonstrates how to use OrderByDescending
// to sort the elements of a slice in descending order by using a selector function
func ExampleQuery_OrderByDescending() {
	names := []string{"Ned", "Ben", "Susan"}

	var result []string
	From(names).
		OrderByDescending(
			func(n interface{}) interface{} { return n },
		).ToSlice(&result)

	fmt.Println(result)
	// Output:
	// [Susan Ned Ben]
}

// The following code example demonstrates how to use ThenByDescending to perform
// a secondary ordering of the elements in a slice in descending order.
func ExampleOrderedQuery_ThenByDescending() {
	fruits := []string{"apPLe", "baNanA", "apple", "APple", "orange", "BAnana", "ORANGE", "apPLE"}

	// Sort the strings first ascending by their length and
	// then descending using a custom case insensitive comparer.
	var query []string
	From(fruits).
		OrderBy(
			func(fruit interface{}) interface{} { return len(fruit.(string)) },
		).
		ThenByDescending(
			func(fruit interface{}) interface{} { return fruit.(string)[0] },
		).
		ToSlice(&query)

	for _, fruit := range query {
		fmt.Println(fruit)
	}
	// Output:
	// apPLe
	// apPLE
	// apple
	// APple
	// orange
	// baNanA
	// ORANGE
	// BAnana

}

// The following code example demonstrates how to use Concat
// to concatenate two slices.
func ExampleQuery_Concat() {
	q := From([]int{1, 2, 3}).
		Concat(From([]int{4, 5, 6}))

	fmt.Println(q.Results())

	// Output:
	// [1 2 3 4 5 6]
}

func ExampleQuery_GroupBy() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	q := From(input).GroupBy(
		func(i interface{}) interface{} { return i.(int) % 2 },
		func(i interface{}) interface{} { return i.(int) })

	fmt.Println(q.OrderBy(func(i interface{}) interface{} {
		return i.(Group).Key
	}).Results())

	// Output:
	// [{0 [2 4 6 8]} {1 [1 3 5 7 9]}]
}

// The following code example demonstrates how to use GroupJoin
// to perform a grouped join on two slices
func ExampleQuery_GroupJoin() {
	fruits := []string{
		"apple",
		"banana",
		"apricot",
		"cherry",
		"clementine",
	}

	q := FromString("abc").
		GroupJoin(From(fruits),
			func(i interface{}) interface{} { return i },
			func(i interface{}) interface{} { return []rune(i.(string))[0] },
			func(outer interface{}, inners []interface{}) interface{} {
				return KeyValue{string(outer.(rune)), inners}
			},
		)

	fmt.Println(q.Results())

	// Output:
	// [{a [apple apricot]} {b [banana]} {c [cherry clementine]}]
}

// The following code example demonstrates how to use Join
// to perform an inner join of two slices based on a common key.
func ExampleQuery_Join() {
	fruits := []string{
		"apple",
		"banana",
		"apricot",
		"cherry",
		"clementine",
	}

	q := Range(1, 10).
		Join(From(fruits),
			func(i interface{}) interface{} { return i },
			func(i interface{}) interface{} { return len(i.(string)) },
			func(outer interface{}, inner interface{}) interface{} {
				return KeyValue{outer, inner}
			},
		)

	fmt.Println(q.Results())

	// Output:
	// [{5 apple} {6 banana} {6 cherry} {7 apricot} {10 clementine}]
}

// The following code example demonstrates how to use OrderBy
// to sort the elements of a slice.
func ExampleQuery_OrderBy() {
	q := Range(1, 10).
		OrderBy(
			func(i interface{}) interface{} { return i.(int) % 2 },
		).
		ThenByDescending(
			func(i interface{}) interface{} { return i },
		)

	fmt.Println(q.Results())

	// Output:
	// [10 8 6 4 2 9 7 5 3 1]
}

// The following code example demonstrates how to use Prepend
// to include an elements in the first position of a slice.
func ExampleQuery_Prepend() {
	input := []int{2, 3, 4, 5}

	q := From(input).Prepend(1)
	first := q.First()

	fmt.Println(first)

	// Output:
	// 1
}

// The following code example demonstrates how to use Reverse
// to reverse the order of elements in a string.
func ExampleQuery_Reverse() {
	input := "apple"

	var output []rune
	From(input).
		Reverse().
		ToSlice(&output)

	fmt.Println(string(output))

	// Output:
	// elppa
}

// The following code example demonstrates how to use Select
// to project over a slice of values.
func ExampleQuery_Select() {
	squares := []int{}

	Range(1, 10).
		Select(
			func(x interface{}) interface{} { return x.(int) * x.(int) },
		).
		ToSlice(&squares)

	fmt.Println(squares)
	// Output:
	// [1 4 9 16 25 36 49 64 81 100]
}

func ExampleQuery_SelectMany() {
	input := [][]int{{1, 2, 3}, {4, 5, 6, 7}}

	q := From(input).
		SelectMany(
			func(i interface{}) Query { return From(i) },
		)

	fmt.Println(q.Results())

	// Output:
	// [1 2 3 4 5 6 7]
}

// The following code example demonstrates how to use Select
// to project over a slice of values and use the index of each element.
func ExampleQuery_SelectIndexed() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}

	result := []string{}
	From(fruits).
		SelectIndexed(
			func(index int, fruit interface{}) interface{} { return fruit.(string)[:index] },
		).
		ToSlice(&result)

	fmt.Println(result)
	// Output:
	// [ b ma ora pass grape]

}

// The following code example demonstrates how to use SelectManyByIndexed
// to perform a one-to-many projection over an array and use the index of each outer element.
func ExampleQuery_SelectManyByIndexed() {
	type Pet struct {
		Name string
	}

	type Person struct {
		Name string
		Pets []Pet
	}

	magnus := Person{
		Name: "Hedlund, Magnus",
		Pets: []Pet{{Name: "Daisy"}},
	}

	terry := Person{
		Name: "Adams, Terry",
		Pets: []Pet{{Name: "Barley"}, {Name: "Boots"}},
	}
	charlotte := Person{
		Name: "Weiss, Charlotte",
		Pets: []Pet{{Name: "Whiskers"}},
	}

	people := []Person{magnus, terry, charlotte}
	var results []string

	From(people).
		SelectManyByIndexed(
			func(index int, person interface{}) Query {
				return From(person.(Person).Pets).
					Select(func(pet interface{}) interface{} {
						return fmt.Sprintf("%d - %s", index, pet.(Pet).Name)
					})
			},
			func(indexedPet, person interface{}) interface{} {
				return fmt.Sprintf("Pet: %s, Owner: %s", indexedPet, person.(Person).Name)
			},
		).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// Pet: 0 - Daisy, Owner: Hedlund, Magnus
	// Pet: 1 - Barley, Owner: Adams, Terry
	// Pet: 1 - Boots, Owner: Adams, Terry
	// Pet: 2 - Whiskers, Owner: Weiss, Charlotte

}

// The following code example demonstrates how to use SelectManyIndexed
// to perform a one-to-many projection over an slice of log data and print out their contents.
func ExampleQuery_SelectManyIndexed() {
	type LogFile struct {
		Name  string
		Lines []string
	}

	file1 := LogFile{
		Name: "file1.log",
		Lines: []string{
			"INFO: 2013/11/05 18:11:01 main.go:44: Special Information",
			"WARNING: 2013/11/05 18:11:01 main.go:45: There is something you need to know about",
			"ERROR: 2013/11/05 18:11:01 main.go:46: Something has failed",
		},
	}

	file2 := LogFile{
		Name: "file2.log",
		Lines: []string{
			"INFO: 2013/11/05 18:11:01 main.go:46: Everything is ok",
		},
	}

	file3 := LogFile{
		Name: "file3.log",
		Lines: []string{
			"2013/11/05 18:42:26 Hello World",
		},
	}

	logFiles := []LogFile{file1, file2, file3}
	var results []string

	From(logFiles).
		SelectManyIndexedT(func(fileIndex int, file LogFile) Query {
			return From(file.Lines).
				SelectIndexedT(func(lineIndex int, line string) string {
					return fmt.Sprintf("File:[%d] - %s => line: %d - %s", fileIndex+1, file.Name, lineIndex+1, line)
				})
		}).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// File:[1] - file1.log => line: 1 - INFO: 2013/11/05 18:11:01 main.go:44: Special Information
	// File:[1] - file1.log => line: 2 - WARNING: 2013/11/05 18:11:01 main.go:45: There is something you need to know about
	// File:[1] - file1.log => line: 3 - ERROR: 2013/11/05 18:11:01 main.go:46: Something has failed
	// File:[2] - file2.log => line: 1 - INFO: 2013/11/05 18:11:01 main.go:46: Everything is ok
	// File:[3] - file3.log => line: 1 - 2013/11/05 18:42:26 Hello World

}

// The following code example demonstrates how to use SelectMany
// to perform a one-to-many projection over a slice
func ExampleQuery_SelectManyBy() {

	type Pet struct {
		Name string
	}

	type Person struct {
		Name string
		Pets []Pet
	}

	magnus := Person{
		Name: "Hedlund, Magnus",
		Pets: []Pet{{Name: "Daisy"}},
	}

	terry := Person{
		Name: "Adams, Terry",
		Pets: []Pet{{Name: "Barley"}, {Name: "Boots"}},
	}
	charlotte := Person{
		Name: "Weiss, Charlotte",
		Pets: []Pet{{Name: "Whiskers"}},
	}

	people := []Person{magnus, terry, charlotte}
	var results []string
	From(people).
		SelectManyBy(
			func(person interface{}) Query { return From(person.(Person).Pets) },
			func(pet, person interface{}) interface{} {
				return fmt.Sprintf("Owner: %s, Pet: %s", person.(Person).Name, pet.(Pet).Name)
			},
		).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// Owner: Hedlund, Magnus, Pet: Daisy
	// Owner: Adams, Terry, Pet: Barley
	// Owner: Adams, Terry, Pet: Boots
	// Owner: Weiss, Charlotte, Pet: Whiskers
}

// The following code example demonstrates how to use SequenceEqual
// to determine whether two slices are equal.
func ExampleQuery_SequenceEqual() {
	type Pet struct {
		Name string
		Age  int
	}

	pets1 := []Pet{
		{Name: "Barley", Age: 8},
		{Name: "Boots", Age: 4},
		{Name: "Whiskers", Age: 1},
		{Name: "Daisy", Age: 4},
	}

	pets2 := []Pet{
		{Name: "Barley", Age: 8},
		{Name: "Boots", Age: 4},
		{Name: "Whiskers", Age: 1},
		{Name: "Daisy", Age: 4},
	}

	equal := From(pets1).SequenceEqual(From(pets2))

	fmt.Printf("Are the lists equals? %t", equal)

	// Output:
	// Are the lists equals? true
}

// The following code example demonstrates how to use Single
// to select the only element of a slice.
func ExampleQuery_Single() {
	fruits1 := []string{"orange"}

	fruit1 := From(fruits1).Single()

	fmt.Println(fruit1)
	// Output:
	// orange
}

// The following code example demonstrates how to use SingleWith
// to select the only element of a slice that satisfies a condition.
func ExampleQuery_SingleWith() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}

	fruit := From(fruits).
		SingleWith(
			func(f interface{}) bool { return len(f.(string)) > 10 },
		)

	fmt.Println(fruit)
	// Output:
	// passionfruit
}

// The following code example demonstrates how to use Skip
// to skip a specified number of elements in a sorted array
// and return the remaining elements.
func ExampleQuery_Skip() {
	grades := []int{59, 82, 70, 56, 92, 98, 85}
	var lowerGrades []int
	From(grades).
		OrderByDescending(
			func(g interface{}) interface{} { return g },
		).
		Skip(3).
		ToSlice(&lowerGrades)

	//All grades except the top three are:
	fmt.Println(lowerGrades)
	// Output:
	// [82 70 59 56]
}

// The following code example demonstrates how to use SkipWhile
// to skip elements of an array as long as a condition is true.
func ExampleQuery_SkipWhile() {
	grades := []int{59, 82, 70, 56, 92, 98, 85}
	var lowerGrades []int
	From(grades).
		OrderByDescending(
			func(g interface{}) interface{} { return g },
		).
		SkipWhile(
			func(g interface{}) bool { return g.(int) >= 80 },
		).
		ToSlice(&lowerGrades)

	// All grades below 80:
	fmt.Println(lowerGrades)
	// Output:
	// [70 59 56]
}

// The following code example demonstrates how to use SkipWhileIndexed
// to skip elements of an array as long as a condition that depends
// on the element's index is true.
func ExampleQuery_SkipWhileIndexed() {
	amounts := []int{5000, 2500, 9000, 8000, 6500, 4000, 1500, 5500}

	var query []int
	From(amounts).
		SkipWhileIndexed(
			func(index int, amount interface{}) bool { return amount.(int) > index*1000 },
		).
		ToSlice(&query)

	fmt.Println(query)
	// Output:
	// [4000 1500 5500]

}

// The following code example demonstrates how to use Sort
// to order elements of an slice.
func ExampleQuery_Sort() {
	amounts := []int{5000, 2500, 9000, 8000, 6500, 4000, 1500, 5500}

	var query []int
	From(amounts).
		Sort(
			func(i interface{}, j interface{}) bool { return i.(int) < j.(int) },
		).
		ToSlice(&query)

	fmt.Println(query)
	// Output:
	// [1500 2500 4000 5000 5500 6500 8000 9000]

}

// The following code example demonstrates how to use SumFloats
// to sum the values of a slice.
func ExampleQuery_SumFloats() {
	numbers := []float64{43.68, 1.25, 583.7, 6.5}

	sum := From(numbers).SumFloats()

	fmt.Printf("The sum of the numbers is %f.", sum)

	// Output:
	// The sum of the numbers is 635.130000.

}

// The following code example demonstrates how to use SumInts
// to sum the values of a slice.
func ExampleQuery_SumInts() {
	numbers := []int{43, 1, 583, 6}

	sum := From(numbers).SumInts()

	fmt.Printf("The sum of the numbers is %d.", sum)

	// Output:
	// The sum of the numbers is 633.

}

// The following code example demonstrates how to use SumUInts
// to sum the values of a slice.
func ExampleQuery_SumUInts() {
	numbers := []uint{43, 1, 583, 6}

	sum := From(numbers).SumUInts()

	fmt.Printf("The sum of the numbers is %d.", sum)

	// Output:
	// The sum of the numbers is 633.

}

// The following code example demonstrates how to use Take
//  to return elements from the start of a slice.
func ExampleQuery_Take() {
	grades := []int{59, 82, 70, 56, 92, 98, 85}

	var topThreeGrades []int
	From(grades).
		OrderByDescending(
			func(grade interface{}) interface{} { return grade },
		).
		Take(3).
		ToSlice(&topThreeGrades)

	fmt.Printf("The top three grades are: %v", topThreeGrades)

	// Output:
	// The top three grades are: [98 92 85]

}

// The following code example demonstrates how to use TakeWhile
// to return elements from the start of a slice.
func ExampleQuery_TakeWhile() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}

	var query []string
	From(fruits).
		TakeWhile(
			func(fruit interface{}) bool { return fruit.(string) != "orange" },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [apple banana mango]
}

// The following code example demonstrates how to use TakeWhileIndexed
// to return elements from the start of a slice as long as
// a condition that uses the element's index is true.
func ExampleQuery_TakeWhileIndexed() {

	fruits := []string{"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry"}

	var query []string
	From(fruits).
		TakeWhileIndexed(
			func(index int, fruit interface{}) bool { return len(fruit.(string)) >= index },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [apple passionfruit banana mango orange blueberry]
}

// The following code example demonstrates how to use ToChannel
// to send a slice to a channel.
func ExampleQuery_ToChannel() {
	c := make(chan interface{})

	go func() {
		Repeat(10, 3).ToChannel(c)
	}()

	for i := range c {
		fmt.Println(i)
	}

	// Output:
	// 10
	// 10
	// 10
}

// The following code example demonstrates how to use ToMap to populate a map.
func ExampleQuery_ToMap() {
	type Product struct {
		Name string
		Code int
	}

	products := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	map1 := map[int]string{}
	From(products).
		SelectT(
			func(item Product) KeyValue { return KeyValue{Key: item.Code, Value: item.Name} },
		).
		ToMap(&map1)

	fmt.Println(map1[4])
	fmt.Println(map1[9])
	fmt.Println(map1[12])

	// Output:
	// orange
	// apple
	// lemon
}

// The following code example demonstrates how to use ToMapBy
// by using a key and value selectors to populate a map.
func ExampleQuery_ToMapBy() {
	input := [][]interface{}{{1, true}}

	result := make(map[int]bool)
	From(input).
		ToMapBy(&result,
			func(i interface{}) interface{} {
				return i.([]interface{})[0]
			},
			func(i interface{}) interface{} {
				return i.([]interface{})[1]
			},
		)

	fmt.Println(result)

	// Output:
	// map[1:true]
}

// The following code example demonstrates how to use ToSlice to populate a slice.
func ExampleQuery_ToSlice() {
	var result []int
	Range(1, 10).ToSlice(&result)

	fmt.Println(result)

	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
}

// The following code example demonstrates how to use Union
// to obtain the union of two slices of integers.
func ExampleQuery_Union() {
	q := Range(1, 10).Union(Range(6, 10))

	fmt.Println(q.Results())

	// Output:
	// [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15]
}

// The following code example demonstrates how to use Where
// to filter a slices.
func ExampleQuery_Where() {
	fruits := []string{"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry"}
	var query []string
	From(fruits).
		Where(
			func(f interface{}) bool { return len(f.(string)) > 6 },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [passionfruit blueberry strawberry]
}

// The following code example demonstrates how to use WhereIndexed
// to filter a slice based on a predicate that involves the index of each element.
func ExampleQuery_WhereIndexed() {
	numbers := []int{0, 30, 20, 15, 90, 85, 40, 75}

	var query []int
	From(numbers).
		WhereIndexed(
			func(index int, number interface{}) bool { return number.(int) <= index*10 },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [0 15 40]
}

// The following code example demonstrates how to use the Zip
// method to merge two slices.
func ExampleQuery_Zip() {
	number := []int{1, 2, 3, 4, 5}
	words := []string{"one", "two", "three"}

	q := From(number).
		Zip(From(words),
			func(a interface{}, b interface{}) interface{} { return []interface{}{a, b} },
		)

	fmt.Println(q.Results())

	// Output:
	// [[1 one] [2 two] [3 three]]
}

// The following code example demonstrates how to use ThenByDescendingT to perform
// a order in a slice of dates by year, and then by month descending.
func ExampleOrderedQuery_ThenByDescendingT() {
	dates := []time.Time{
		time.Date(2015, 3, 23, 0, 0, 0, 0, time.Local),
		time.Date(2014, 7, 11, 0, 0, 0, 0, time.Local),
		time.Date(2013, 5, 4, 0, 0, 0, 0, time.Local),
		time.Date(2015, 1, 2, 0, 0, 0, 0, time.Local),
		time.Date(2015, 7, 10, 0, 0, 0, 0, time.Local),
	}

	var orderedDates []time.Time
	From(dates).
		OrderByT(
			func(date time.Time) int {
				return date.Year()
			}).
		ThenByDescendingT(
			func(date time.Time) int { return int(date.Month()) },
		).
		ToSlice(&orderedDates)

	for _, date := range orderedDates {
		fmt.Println(date.Format("2006-Jan-02"))
	}
	// Output:
	// 2013-May-04
	// 2014-Jul-11
	// 2015-Jul-10
	// 2015-Mar-23
	// 2015-Jan-02

}

// The following code example demonstrates how to use ThenByT to perform
// a orders in a slice of dates by year, and then by day.
func ExampleOrderedQuery_ThenByT() {
	dates := []time.Time{
		time.Date(2015, 3, 23, 0, 0, 0, 0, time.Local),
		time.Date(2014, 7, 11, 0, 0, 0, 0, time.Local),
		time.Date(2013, 5, 4, 0, 0, 0, 0, time.Local),
		time.Date(2015, 1, 2, 0, 0, 0, 0, time.Local),
		time.Date(2015, 7, 10, 0, 0, 0, 0, time.Local),
	}

	var orderedDates []time.Time
	From(dates).
		OrderByT(
			func(date time.Time) int { return date.Year() },
		).
		ThenByT(
			func(date time.Time) int { return int(date.Day()) },
		).
		ToSlice(&orderedDates)

	for _, date := range orderedDates {
		fmt.Println(date.Format("2006-Jan-02"))
	}
	// Output:
	// 2013-May-04
	// 2014-Jul-11
	// 2015-Jan-02
	// 2015-Jul-10
	// 2015-Mar-23

}

// The following code example demonstrates how to reverse
// the order of words in a string using AggregateT.
func ExampleQuery_AggregateT() {
	sentence := "the quick brown fox jumps over the lazy dog"
	// Split the string into individual words.
	words := strings.Split(sentence, " ")

	// Prepend each word to the beginning of the
	// new sentence to reverse the word order.
	reversed := From(words).AggregateT(
		func(workingSentence string, next string) string { return next + " " + workingSentence },
	)

	fmt.Println(reversed)

	// Output:
	// dog lazy the over jumps fox brown quick the
}

// The following code example demonstrates how to use AggregateWithSeed function
func ExampleQuery_AggregateWithSeedT() {

	fruits := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	// Determine whether any string in the array is longer than "banana".
	longestName := From(fruits).
		AggregateWithSeedT("banana",
			func(longest, next string) string {
				if len(next) > len(longest) {
					return next
				}
				return longest
			},
		)

	fmt.Printf("The fruit with the longest name is %s.", longestName)

	// Output:
	// The fruit with the longest name is passionfruit.

}

// The following code example demonstrates how to use AggregateWithSeedByT function
func ExampleQuery_AggregateWithSeedByT() {
	input := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	// Determine whether any string in the array is longer than "banana".
	longestName := From(input).AggregateWithSeedByT("banana",
		func(longest string, next string) string {
			if len(longest) > len(next) {
				return longest
			}
			return next

		},
		// Return the final result
		func(result string) string {
			return fmt.Sprintf("The fruit with the longest name is %s.", result)
		},
	)

	fmt.Println(longestName)
	// Output:
	// The fruit with the longest name is passionfruit.
}

// The following code example demonstrates how to use AllT
// to get the students having all marks greater than 70.
func ExampleQuery_AllT() {

	type Student struct {
		Name  string
		Marks []int
	}

	students := []Student{
		{Name: "Hugo", Marks: []int{91, 88, 76, 93}},
		{Name: "Rick", Marks: []int{70, 73, 66, 90}},
		{Name: "Michael", Marks: []int{73, 80, 75, 88}},
		{Name: "Fadi", Marks: []int{82, 75, 66, 84}},
		{Name: "Peter", Marks: []int{67, 78, 70, 82}},
	}

	var approvedStudents []Student
	From(students).
		WhereT(
			func(student Student) bool {
				return From(student.Marks).
					AllT(
						func(mark int) bool { return mark > 70 },
					)
			},
		).
		ToSlice(&approvedStudents)

	//List of approved students
	for _, student := range approvedStudents {
		fmt.Println(student.Name)
	}

	// Output:
	// Hugo
	// Michael
}

// The following code example demonstrates how to use AnyWithT
// to get the students with any mark lower than 70.
func ExampleQuery_AnyWithT() {
	type Student struct {
		Name  string
		Marks []int
	}

	students := []Student{
		{Name: "Hugo", Marks: []int{91, 88, 76, 93}},
		{Name: "Rick", Marks: []int{70, 73, 66, 90}},
		{Name: "Michael", Marks: []int{73, 80, 75, 88}},
		{Name: "Fadi", Marks: []int{82, 75, 66, 84}},
		{Name: "Peter", Marks: []int{67, 78, 70, 82}},
	}

	var studentsWithAnyMarkLt70 []Student
	From(students).
		WhereT(
			func(student Student) bool {
				return From(student.Marks).
					AnyWithT(
						func(mark int) bool { return mark < 70 },
					)
			},
		).
		ToSlice(&studentsWithAnyMarkLt70)

	//List of students with any mark lower than 70
	for _, student := range studentsWithAnyMarkLt70 {
		fmt.Println(student.Name)
	}

	// Output:
	// Rick
	// Fadi
	// Peter

}

// The following code example demonstrates how to use CountWithT
// to count the elements in an slice that satisfy a condition.
func ExampleQuery_CountWithT() {
	type Pet struct {
		Name       string
		Vaccinated bool
	}

	pets := []Pet{
		{Name: "Barley", Vaccinated: true},
		{Name: "Boots", Vaccinated: false},
		{Name: "Whiskers", Vaccinated: false},
	}

	numberUnvaccinated := From(pets).
		CountWithT(
			func(p Pet) bool { return p.Vaccinated == false },
		)

	fmt.Printf("There are %d unvaccinated animals.", numberUnvaccinated)

	//Output:
	//There are 2 unvaccinated animals.
}

// The following code example demonstrates how to use DistinctByT
// to return distinct elements from a slice of structs.
func ExampleQuery_DistinctByT() {
	type Product struct {
		Name string
		Code int
	}

	products := []Product{
		{Name: "apple", Code: 9},
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
	}

	//Exclude duplicates.
	var noduplicates []Product
	From(products).
		DistinctByT(
			func(item Product) int { return item.Code },
		).
		ToSlice(&noduplicates)

	for _, product := range noduplicates {
		fmt.Printf("%s %d\n", product.Name, product.Code)
	}

	// Output:
	// apple 9
	// orange 4
	// lemon 12
}

// The following code example demonstrates how to use ExceptByT
func ExampleQuery_ExceptByT() {
	type Product struct {
		Name string
		Code int
	}

	fruits1 := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	fruits2 := []Product{
		{Name: "apple", Code: 9},
	}

	//Order and exclude duplicates.
	var except []Product
	From(fruits1).
		ExceptByT(From(fruits2),
			func(item Product) int { return item.Code },
		).
		ToSlice(&except)

	for _, product := range except {
		fmt.Printf("%s %d\n", product.Name, product.Code)
	}

	// Output:
	// orange 4
	// lemon 12

}

// The following code example demonstrates how to use FirstWithT
// to return the first element of an array that satisfies a condition.
func ExampleQuery_FirstWithT() {
	numbers := []int{9, 34, 65, 92, 87, 435, 3, 54, 83, 23, 87, 435, 67, 12, 19}

	first := From(numbers).
		FirstWithT(
			func(item int) bool { return item > 80 },
		)

	fmt.Println(first)

	// Output:
	// 92

}

// The following code example demonstrates how to use ForEach
// to output all elements of an array.
func ExampleQuery_ForEach() {
	fruits := []string{"orange", "apple", "lemon", "apple"}

	From(fruits).ForEach(func(fruit interface{}) {
		fmt.Println(fruit)
	})

	// Output:
	// orange
	// apple
	// lemon
	// apple
}

// The following code example demonstrates how to use ForEachIndexed
// to output all elements of an array with its index.
func ExampleQuery_ForEachIndexed() {
	fruits := []string{"orange", "apple", "lemon", "apple"}

	From(fruits).ForEachIndexed(func(i int, fruit interface{}) {
		fmt.Printf("%d.%s\n", i, fruit)
	})

	// Output:
	// 0.orange
	// 1.apple
	// 2.lemon
	// 3.apple
}

// The following code example demonstrates how to use ForEachT
// to output all elements of an array.
func ExampleQuery_ForEachT() {
	fruits := []string{"orange", "apple", "lemon", "apple"}

	From(fruits).ForEachT(func(fruit string) {
		fmt.Println(fruit)
	})

	// Output:
	// orange
	// apple
	// lemon
	// apple
}

// The following code example demonstrates how to use ForEachIndexedT
// to output all elements of an array with its index.
func ExampleQuery_ForEachIndexedT() {
	fruits := []string{"orange", "apple", "lemon", "apple"}

	From(fruits).ForEachIndexedT(func(i int, fruit string) {
		fmt.Printf("%d.%s\n", i, fruit)
	})

	// Output:
	// 0.orange
	// 1.apple
	// 2.lemon
	// 3.apple
}

// The following code example demonstrates how to use GroupByT
// to group the elements of a slice.
func ExampleQuery_GroupByT() {

	type Pet struct {
		Name string
		Age  int
	}
	// Create a list of pets.
	pets := []Pet{
		{Name: "Barley", Age: 8},
		{Name: "Boots", Age: 4},
		{Name: "Whiskers", Age: 1},
		{Name: "Daisy", Age: 4},
	}

	// Group the pets using Age as the key value
	// and selecting only the pet's Name for each value.
	var query []Group
	From(pets).GroupByT(
		func(p Pet) int { return p.Age },
		func(p Pet) string { return p.Name },
	).OrderByT(
		func(g Group) int { return g.Key.(int) },
	).ToSlice(&query)

	for _, petGroup := range query {
		fmt.Printf("%d\n", petGroup.Key)
		for _, petName := range petGroup.Group {
			fmt.Printf("  %s\n", petName)
		}

	}

	// Output:
	// 1
	//   Whiskers
	// 4
	//   Boots
	//   Daisy
	// 8
	//   Barley
}

// The following code example demonstrates how to use GroupJoinT
//  to perform a grouped join on two slices.
func ExampleQuery_GroupJoinT() {

	type Person struct {
		Name string
	}

	type Pet struct {
		Name  string
		Owner Person
	}

	magnus := Person{Name: "Hedlund, Magnus"}
	terry := Person{Name: "Adams, Terry"}
	charlotte := Person{Name: "Weiss, Charlotte"}

	barley := Pet{Name: "Barley", Owner: terry}
	boots := Pet{Name: "Boots", Owner: terry}
	whiskers := Pet{Name: "Whiskers", Owner: charlotte}
	daisy := Pet{Name: "Daisy", Owner: magnus}

	people := []Person{magnus, terry, charlotte}
	pets := []Pet{barley, boots, whiskers, daisy}

	// Create a slice where each element is a KeyValue
	// that contains a person's name as the key and a slice of strings
	// of names of the pets they own as a value.

	q := []KeyValue{}
	From(people).
		GroupJoinT(From(pets),
			func(p Person) Person { return p },
			func(p Pet) Person { return p.Owner },
			func(person Person, pets []Pet) KeyValue {
				var petNames []string
				From(pets).
					SelectT(
						func(pet Pet) string { return pet.Name },
					).
					ToSlice(&petNames)
				return KeyValue{person.Name, petNames}
			},
		).ToSlice(&q)

	for _, obj := range q {
		// Output the owner's name.
		fmt.Printf("%s:\n", obj.Key)
		// Output each of the owner's pet's names.
		for _, petName := range obj.Value.([]string) {
			fmt.Printf("  %s\n", petName)
		}
	}

	// Output:
	// Hedlund, Magnus:
	//   Daisy
	// Adams, Terry:
	//   Barley
	//   Boots
	// Weiss, Charlotte:
	//   Whiskers
}

// The following code example demonstrates how to use IntersectByT
// to return the elements that appear in each of two slices of products
// with same Code.
func ExampleQuery_IntersectByT() {
	type Product struct {
		Name string
		Code int
	}

	store1 := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
	}

	store2 := []Product{
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	var duplicates []Product
	From(store1).
		IntersectByT(From(store2),
			func(p Product) int { return p.Code },
		).
		ToSlice(&duplicates)

	for _, p := range duplicates {
		fmt.Println(p.Name, "", p.Code)
	}

	// Output:
	// apple  9

}

// The following code example demonstrates how to use JoinT
// to perform an inner join of two slices based on a common key.
func ExampleQuery_JoinT() {
	type Person struct {
		Name string
	}

	type Pet struct {
		Name  string
		Owner Person
	}

	magnus := Person{Name: "Hedlund, Magnus"}
	terry := Person{Name: "Adams, Terry"}
	charlotte := Person{Name: "Weiss, Charlotte"}

	barley := Pet{Name: "Barley", Owner: terry}
	boots := Pet{Name: "Boots", Owner: terry}
	whiskers := Pet{Name: "Whiskers", Owner: charlotte}
	daisy := Pet{Name: "Daisy", Owner: magnus}

	people := []Person{magnus, terry, charlotte}
	pets := []Pet{barley, boots, whiskers, daisy}

	// Create a list of Person-Pet pairs where
	// each element is an anonymous type that contains a
	// Pet's name and the name of the Person that owns the Pet.

	query := []string{}
	From(people).
		JoinT(From(pets),
			func(person Person) Person { return person },
			func(pet Pet) Person { return pet.Owner },
			func(person Person, pet Pet) string { return fmt.Sprintf("%s - %s", person.Name, pet.Name) },
		).ToSlice(&query)

	for _, line := range query {
		fmt.Println(line)
	}
	//Output:
	//Hedlund, Magnus - Daisy
	//Adams, Terry - Barley
	//Adams, Terry - Boots
	//Weiss, Charlotte - Whiskers
}

// The following code example demonstrates how to use LastWithT
// to return the last element of an array.
func ExampleQuery_LastWithT() {
	numbers := []int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 67, 12, 19}

	last := From(numbers).
		LastWithT(
			func(n int) bool { return n > 80 },
		)

	fmt.Println(last)

	//Output:
	//87

}

// The following code example demonstrates how to use OrderByDescendingT
// to order an slice.
func ExampleQuery_OrderByDescendingT() {
	type Player struct {
		Name   string
		Points int64
	}

	players := []Player{
		{Name: "Hugo", Points: 4757},
		{Name: "Rick", Points: 7365},
		{Name: "Michael", Points: 2857},
		{Name: "Fadi", Points: 85897},
		{Name: "Peter", Points: 48576},
	}

	//Order and get the top 3 players
	var top3Players []KeyValue
	From(players).
		OrderByDescendingT(
			func(p Player) int64 { return p.Points },
		).
		Take(3).
		SelectIndexedT(
			func(i int, p Player) KeyValue { return KeyValue{Key: i + 1, Value: p} },
		).
		ToSlice(&top3Players)

	for _, rank := range top3Players {
		fmt.Printf(
			"Rank: #%d - Player: %s - Points: %d\n",
			rank.Key,
			rank.Value.(Player).Name,
			rank.Value.(Player).Points,
		)

	}
	// Output:
	// Rank: #1 - Player: Fadi - Points: 85897
	// Rank: #2 - Player: Peter - Points: 48576
	// Rank: #3 - Player: Rick - Points: 7365
}

// The following code example demonstrates how to use OrderByT
// to sort the elements of a slice.
func ExampleQuery_OrderByT() {
	type Pet struct {
		Name string
		Age  int
	}
	// Create a list of pets.
	pets := []Pet{
		{Name: "Barley", Age: 8},
		{Name: "Boots", Age: 4},
		{Name: "Whiskers", Age: 1},
		{Name: "Daisy", Age: 4},
	}

	var orderedPets []Pet
	From(pets).
		OrderByT(
			func(pet Pet) int { return pet.Age },
		).
		ToSlice(&orderedPets)

	for _, pet := range orderedPets {
		fmt.Println(pet.Name, "-", pet.Age)
	}

	// Output:
	// Whiskers - 1
	// Boots - 4
	// Daisy - 4
	// Barley - 8
}

// The following code example demonstrates how to use SelectT
// to project over a slice.
func ExampleQuery_SelectT() {
	squares := []int{}

	Range(1, 10).
		SelectT(
			func(x int) int { return x * x },
		).
		ToSlice(&squares)

	fmt.Println(squares)
	// Output:
	// [1 4 9 16 25 36 49 64 81 100]
}

// The following code example demonstrates how to use SelectIndexedT
// to determine if the value in a slice of int match their position
// in the slice.
func ExampleQuery_SelectIndexedT() {
	numbers := []int{5, 4, 1, 3, 9, 8, 6, 7, 2, 0}

	var numsInPlace []KeyValue

	From(numbers).
		SelectIndexedT(
			func(index, num int) KeyValue { return KeyValue{Key: num, Value: (num == index)} },
		).
		ToSlice(&numsInPlace)

	fmt.Println("Number: In-place?")
	for _, n := range numsInPlace {
		fmt.Printf("%d: %t\n", n.Key, n.Value)
	}

	// Output:
	// Number: In-place?
	// 5: false
	// 4: false
	// 1: false
	// 3: true
	// 9: false
	// 8: false
	// 6: true
	// 7: true
	// 2: false
	// 0: false

}

// The following code example demonstrates how to use SelectManyT
// to perform a one-to-many projection over a slice
func ExampleQuery_SelectManyByT() {

	type Pet struct {
		Name string
	}

	type Person struct {
		Name string
		Pets []Pet
	}

	magnus := Person{
		Name: "Hedlund, Magnus",
		Pets: []Pet{{Name: "Daisy"}},
	}

	terry := Person{
		Name: "Adams, Terry",
		Pets: []Pet{{Name: "Barley"}, {Name: "Boots"}},
	}
	charlotte := Person{
		Name: "Weiss, Charlotte",
		Pets: []Pet{{Name: "Whiskers"}},
	}

	people := []Person{magnus, terry, charlotte}
	var results []string
	From(people).
		SelectManyByT(
			func(person Person) Query { return From(person.Pets) },
			func(pet Pet, person Person) interface{} {
				return fmt.Sprintf("Owner: %s, Pet: %s", person.Name, pet.Name)
			},
		).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// Owner: Hedlund, Magnus, Pet: Daisy
	// Owner: Adams, Terry, Pet: Barley
	// Owner: Adams, Terry, Pet: Boots
	// Owner: Weiss, Charlotte, Pet: Whiskers
}

// The following code example demonstrates how to use SelectManyT
// to perform a projection over a list of sentences and rank the
// top 5 most used words
func ExampleQuery_SelectManyT() {
	sentences := []string{
		"the quick brown fox jumps over the lazy dog",
		"pack my box with five dozen liquor jugs",
		"several fabulous dixieland jazz groups played with quick tempo",
		"back in my quaint garden jaunty zinnias vie with flaunting phlox",
		"five or six big jet planes zoomed quickly by the new tower",
		"I quickly explained that many big jobs involve few hazards",
		"The wizard quickly jinxed the gnomes before they vaporized",
	}

	var results []string
	From(sentences).
		//Split the sentences in words
		SelectManyT(func(sentence string) Query {
			return From(strings.Split(sentence, " "))
		}).
		//Grouping by word
		GroupByT(
			func(word string) string { return word },
			func(word string) string { return word },
		).
		//Ordering by word counts
		OrderByDescendingT(func(wordGroup Group) int {
			return len(wordGroup.Group)
		}).
		//Then order by word
		ThenByT(func(wordGroup Group) string {
			return wordGroup.Key.(string)
		}).
		//Take the top 5
		Take(5).
		//Project the words using the index as rank
		SelectIndexedT(func(index int, wordGroup Group) string {
			return fmt.Sprintf("Rank: #%d, Word: %s, Counts: %d", index+1, wordGroup.Key, len(wordGroup.Group))
		}).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// Rank: #1, Word: the, Counts: 4
	// Rank: #2, Word: quickly, Counts: 3
	// Rank: #3, Word: with, Counts: 3
	// Rank: #4, Word: big, Counts: 2
	// Rank: #5, Word: five, Counts: 2
}

// The following code example demonstrates how to use SelectManyIndexedT
// to perform a one-to-many projection over an slice of log files and
// print out their contents.
func ExampleQuery_SelectManyIndexedT() {
	type LogFile struct {
		Name  string
		Lines []string
	}

	file1 := LogFile{
		Name: "file1.log",
		Lines: []string{
			"INFO: 2013/11/05 18:11:01 main.go:44: Special Information",
			"WARNING: 2013/11/05 18:11:01 main.go:45: There is something you need to know about",
			"ERROR: 2013/11/05 18:11:01 main.go:46: Something has failed",
		},
	}

	file2 := LogFile{
		Name: "file2.log",
		Lines: []string{
			"INFO: 2013/11/05 18:11:01 main.go:46: Everything is ok",
		},
	}

	file3 := LogFile{
		Name: "file3.log",
		Lines: []string{
			"2013/11/05 18:42:26 Hello World",
		},
	}

	logFiles := []LogFile{file1, file2, file3}
	var results []string

	From(logFiles).
		SelectManyIndexedT(func(fileIndex int, file LogFile) Query {
			return From(file.Lines).
				SelectIndexedT(func(lineIndex int, line string) string {
					return fmt.Sprintf("File:[%d] - %s => line: %d - %s", fileIndex+1, file.Name, lineIndex+1, line)
				})
		}).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// File:[1] - file1.log => line: 1 - INFO: 2013/11/05 18:11:01 main.go:44: Special Information
	// File:[1] - file1.log => line: 2 - WARNING: 2013/11/05 18:11:01 main.go:45: There is something you need to know about
	// File:[1] - file1.log => line: 3 - ERROR: 2013/11/05 18:11:01 main.go:46: Something has failed
	// File:[2] - file2.log => line: 1 - INFO: 2013/11/05 18:11:01 main.go:46: Everything is ok
	// File:[3] - file3.log => line: 1 - 2013/11/05 18:42:26 Hello World

}

// The following code example demonstrates how to use SelectManyByIndexedT
// to perform a one-to-many projection over an array and use the index of
// each outer element.
func ExampleQuery_SelectManyByIndexedT() {
	type Pet struct {
		Name string
	}

	type Person struct {
		Name string
		Pets []Pet
	}

	magnus := Person{
		Name: "Hedlund, Magnus",
		Pets: []Pet{{Name: "Daisy"}},
	}

	terry := Person{
		Name: "Adams, Terry",
		Pets: []Pet{{Name: "Barley"}, {Name: "Boots"}},
	}
	charlotte := Person{
		Name: "Weiss, Charlotte",
		Pets: []Pet{{Name: "Whiskers"}},
	}

	people := []Person{magnus, terry, charlotte}
	var results []string

	From(people).
		SelectManyByIndexedT(
			func(index int, person Person) Query {
				return From(person.Pets).
					SelectT(func(pet Pet) string {
						return fmt.Sprintf("%d - %s", index, pet.Name)
					})
			},
			func(indexedPet string, person Person) string {
				return fmt.Sprintf("Pet: %s, Owner: %s", indexedPet, person.Name)
			},
		).
		ToSlice(&results)

	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// Pet: 0 - Daisy, Owner: Hedlund, Magnus
	// Pet: 1 - Barley, Owner: Adams, Terry
	// Pet: 1 - Boots, Owner: Adams, Terry
	// Pet: 2 - Whiskers, Owner: Weiss, Charlotte

}

//The following code example demonstrates how to use SingleWithT
// to select the only element of a slice that satisfies a condition.
func ExampleQuery_SingleWithT() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}

	fruit := From(fruits).
		SingleWithT(
			func(f string) bool { return len(f) > 10 },
		)

	fmt.Println(fruit)
	// Output:
	// passionfruit
}

// The following code example demonstrates how to use SkipWhileT
// to skip elements of an array as long as a condition is true.
func ExampleQuery_SkipWhileT() {
	grades := []int{59, 82, 70, 56, 92, 98, 85}
	var lowerGrades []int
	From(grades).
		OrderByDescendingT(
			func(g int) int { return g },
		).
		SkipWhileT(
			func(g int) bool { return g >= 80 },
		).
		ToSlice(&lowerGrades)

	//"All grades below 80:
	fmt.Println(lowerGrades)
	// Output:
	// [70 59 56]
}

// The following code example demonstrates how to use SkipWhileIndexedT
// to skip elements of an array as long as a condition that depends
// on the element's index is true.
func ExampleQuery_SkipWhileIndexedT() {
	amounts := []int{5000, 2500, 9000, 8000, 6500, 4000, 1500, 5500}

	var query []int
	From(amounts).
		SkipWhileIndexedT(
			func(index int, amount int) bool { return amount > index*1000 },
		).
		ToSlice(&query)

	fmt.Println(query)
	// Output:
	// [4000 1500 5500]

}

// The following code example demonstrates how to use SortT
// to order elements of an slice.
func ExampleQuery_SortT() {
	type Pet struct {
		Name string
		Age  int
	}
	// Create a list of pets.
	pets := []Pet{
		{Name: "Barley", Age: 8},
		{Name: "Boots", Age: 4},
		{Name: "Whiskers", Age: 1},
		{Name: "Daisy", Age: 4},
	}

	orderedPets := []Pet{}
	From(pets).
		SortT(
			func(pet1 Pet, pet2 Pet) bool { return pet1.Age < pet2.Age },
		).
		ToSlice(&orderedPets)

	for _, pet := range orderedPets {
		fmt.Println(pet.Name, "-", pet.Age)
	}

	// Output:
	// Whiskers - 1
	// Boots - 4
	// Daisy - 4
	// Barley - 8

}

// The following code example demonstrates how to use TakeWhileT
// to return elements from the start of a slice.
func ExampleQuery_TakeWhileT() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}

	var query []string
	From(fruits).
		TakeWhileT(
			func(fruit string) bool { return fruit != "orange" },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [apple banana mango]
}

// The following code example demonstrates how to use TakeWhileIndexedT
// to return elements from the start of a slice as long asa condition
// that uses the element's index is true.
func ExampleQuery_TakeWhileIndexedT() {

	fruits := []string{"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry"}

	var query []string
	From(fruits).
		TakeWhileIndexedT(
			func(index int, fruit string) bool { return len(fruit) >= index },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [apple passionfruit banana mango orange blueberry]
}

// The following code example demonstrates how to use ToMapBy
// by using a key and value selectors to populate a map.
func ExampleQuery_ToMapByT() {
	type Product struct {
		Name string
		Code int
	}

	products := []Product{
		{Name: "orange", Code: 4},
		{Name: "apple", Code: 9},
		{Name: "lemon", Code: 12},
		{Name: "apple", Code: 9},
	}

	map1 := map[int]string{}
	From(products).
		ToMapByT(&map1,
			func(item Product) int { return item.Code },
			func(item Product) string { return item.Name },
		)

	fmt.Println(map1[4])
	fmt.Println(map1[9])
	fmt.Println(map1[12])

	// Output:
	// orange
	// apple
	// lemon
}

// The following code example demonstrates how to use WhereT
// to filter a slices.
func ExampleQuery_WhereT() {
	fruits := []string{"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry"}
	var query []string
	From(fruits).
		WhereT(
			func(f string) bool { return len(f) > 6 },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [passionfruit blueberry strawberry]
}

// The following code example demonstrates how to use WhereIndexedT
// to filter a slice based on a predicate that involves the index of each element.
func ExampleQuery_WhereIndexedT() {
	numbers := []int{0, 30, 20, 15, 90, 85, 40, 75}

	var query []int
	From(numbers).
		WhereIndexedT(
			func(index int, number int) bool { return number <= index*10 },
		).
		ToSlice(&query)

	fmt.Println(query)

	// Output:
	// [0 15 40]
}

// The following code example demonstrates how to use the Zip
// method to merge two slices.
func ExampleQuery_ZipT() {
	number := []int{1, 2, 3, 4, 5}
	words := []string{"one", "two", "three"}

	q := From(number).
		ZipT(From(words),
			func(a int, b string) []interface{} { return []interface{}{a, b} },
		)

	fmt.Println(q.Results())

	// Output:
	// [[1 one] [2 two] [3 three]]
}
