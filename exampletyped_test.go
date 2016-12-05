package linq

import "fmt"
import "strings"
import "time"

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
	//
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

// The following code example demonstrates how to
// use DistinctByT to return distinct elements from a slice of structs.
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
// to return the elements that appear in each of two slices of products with same Code.
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

// The following code example demonstrates how to use OrderByDescendingT to order an slice.
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
// to determine if the value in a slice of int match their position in the slice.
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
// to perform a projection over a list of sentences and rank the top 5 most used words
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
// to perform a one-to-many projection over an slice of log files and print out their contents.
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
// to perform a one-to-many projection over an array and use the index of each outer element.
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
//  to return elements from the start of a slice.
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
// to return elements from the start of a slice as long as
// a condition that uses the element's index is true.
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
