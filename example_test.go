package linq

import (
	"fmt"
	"strings"
)

// The Query type can be iterated manually. The yield function returns false
// to stop the iteration early.
func ExampleQuery() {
	query := FromSlice([]int{1, 2, 3, 4, 5}).Where(func(i int) bool {
		return i <= 3
	})

	query.Iterate(func(item int) bool {
		fmt.Println(item)
		return true
	})

	// Output:
	// 1
	// 2
	// 3
}

func ExampleKeyValue() {
	m := map[string]bool{"foo": true}
	fmt.Println(FromMap(m).ToSlice())

	// Output:
	// [{foo true}]
}

// MyQuery shows how to extend the Query type with custom operators.
type MyQuery Query[int]

func (q MyQuery) GreaterThan(threshold int) Query[int] {
	return Query[int](q).Where(func(item int) bool {
		return item > threshold
	})
}

func ExampleMyQuery_GreaterThan() {
	result := MyQuery(Range(1, 10)).GreaterThan(5).ToSlice()
	fmt.Println(result)

	// Output:
	// [6 7 8 9 10]
}

func ExampleFromChannel() {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)

	fmt.Println(FromChannel(ch).ToSlice())

	// Output:
	// [10 20 30]
}

func ExampleFromString() {
	q := FromString("linq").Select(func(r rune) string {
		return strings.ToUpper(string(r))
	})

	fmt.Println(q.ToSlice())

	// Output:
	// [L I N Q]
}

func ExampleFromSeq() {
	seq := func(yield func(int) bool) {
		for i := 1; i <= 4; i++ {
			if !yield(i * i) {
				return
			}
		}
	}

	fmt.Println(FromSeq(seq).ToSlice())

	// Output:
	// [1 4 9 16]
}

func ExampleRange() {
	fmt.Println(Range(1, 5).ToSlice())

	// Output:
	// [1 2 3 4 5]
}

func ExampleSequence() {
	fmt.Println(Sequence(2, 10, 3).ToSlice())

	// Output:
	// [2 5 8]
}

func ExampleInfiniteSequence() {
	fmt.Println(InfiniteSequence(10, -2).Take(4).ToSlice())

	// Output:
	// [10 8 6 4]
}

func ExampleRepeat() {
	fmt.Println(Repeat("go", 3).ToSlice())

	// Output:
	// [go go go]
}

func ExampleQuery_Where() {
	fruits := []string{"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry"}

	query := FromSlice(fruits).Where(func(f string) bool {
		return len(f) > 6
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [passionfruit blueberry strawberry]
}

func ExampleQuery_WhereIndexed() {
	numbers := []int{0, 30, 20, 15, 90, 85, 40, 75}

	query := FromSlice(numbers).WhereIndexed(func(index int, number int) bool {
		return number <= index*10
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [0 20 15 40]
}

func ExampleQuery_Select() {
	squares := Range(1, 5).Select(func(x int) int {
		return x * x
	})

	fmt.Println(squares.ToSlice())

	// Output:
	// [1 4 9 16 25]
}

func ExampleQuery_SelectIndexed() {
	fruits := []string{"apple", "banana", "mango", "orange"}

	query := FromSlice(fruits).SelectIndexed(func(index int, fruit string) string {
		return fmt.Sprintf("%d:%s", index, fruit)
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [0:apple 1:banana 2:mango 3:orange]
}

func ExampleQuery_SelectMany() {
	numbers := [][]int{{1, 2, 3}, {4}, {5, 6}}

	query := FromSlice(numbers).SelectMany(func(ns []int) Query[int] {
		return FromSlice(ns)
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [1 2 3 4 5 6]
}

func ExampleQuery_SelectManyBy() {
	type pet struct {
		Name string
	}
	type person struct {
		Name string
		Pets []pet
	}

	people := []person{
		{Name: "Hines", Pets: []pet{{Name: "Barley"}, {Name: "Boots"}}},
		{Name: "Ashkenazi", Pets: []pet{{Name: "Sugar"}}},
	}

	query := FromSlice(people).SelectManyBy(
		func(p person) Query[pet] { return FromSlice(p.Pets) },
		func(p pet, owner person) string { return owner.Name + ": " + p.Name },
	)

	fmt.Println(query.ToSlice())

	// Output:
	// [Hines: Barley Hines: Boots Ashkenazi: Sugar]
}

func ExampleQuery_GroupBy() {
	query := Range(1, 9).GroupBy(
		func(i int) int { return i % 2 },
		func(i int) int { return i },
	)

	for group := range query.Iterate {
		fmt.Println(group.Key, group.Group)
	}

	// Output:
	// 1 [1 3 5 7 9]
	// 0 [2 4 6 8]
}

func ExampleQuery_Join() {
	type person struct {
		Name string
	}
	type pet struct {
		Name  string
		Owner string
	}

	people := []person{{"Magnus"}, {"Terry"}, {"Charlotte"}}
	pets := []pet{
		{"Barley", "Terry"},
		{"Boots", "Terry"},
		{"Whiskers", "Charlotte"},
		{"Daisy", "Magnus"},
	}

	query := FromSlice(people).Join(
		FromSlice(pets),
		func(p person) string { return p.Name },
		func(p pet) string { return p.Owner },
		func(owner person, p pet) string {
			return owner.Name + " - " + p.Name
		})

	for s := range query.Iterate {
		fmt.Println(s)
	}

	// Output:
	// Magnus - Daisy
	// Terry - Barley
	// Terry - Boots
	// Charlotte - Whiskers
}

func ExampleQuery_LeftJoin() {
	query := FromSlice([]string{"apple", "banana"}).LeftJoin(
		FromSlice([]string{"apricot"}),
		func(s string) byte { return s[0] },
		func(s string) byte { return s[0] },
		func(left, right string) string {
			if right == "" {
				right = "<none>"
			}
			return left + " - " + right
		},
	)

	fmt.Println(query.ToSlice())

	// Output:
	// [apple - apricot banana - <none>]
}

func ExampleQuery_RightJoin() {
	query := FromSlice([]string{"apple"}).RightJoin(
		FromSlice([]string{"apricot", "banana"}),
		func(s string) byte { return s[0] },
		func(s string) byte { return s[0] },
		func(left, right string) string {
			if left == "" {
				left = "<none>"
			}
			return left + " - " + right
		},
	)

	fmt.Println(query.ToSlice())

	// Output:
	// [apple - apricot <none> - banana]
}

func ExampleQuery_FullJoin() {
	query := FromSlice([]string{"apple", "banana"}).FullJoin(
		FromSlice([]string{"apricot", "cherry"}),
		func(s string) byte { return s[0] },
		func(s string) byte { return s[0] },
		func(left, right string) string {
			if left == "" {
				left = "<none>"
			}
			if right == "" {
				right = "<none>"
			}
			return left + " - " + right
		},
	)

	fmt.Println(query.ToSlice())

	// Output:
	// [apple - apricot banana - <none> <none> - cherry]
}

func ExampleQuery_GroupJoin() {
	fruits := []string{"apple", "banana", "apricot", "cherry", "clementine"}

	query := FromString("abc").GroupJoin(
		FromSlice(fruits),
		func(r rune) rune { return r },
		func(fruit string) rune { return rune(fruit[0]) },
		func(r rune, fruits []string) string {
			return fmt.Sprintf("%s: %v", string(r), fruits)
		})

	for s := range query.Iterate {
		fmt.Println(s)
	}

	// Output:
	// a: [apple apricot]
	// b: [banana]
	// c: [cherry clementine]
}

func ExampleQuery_OrderBy() {
	query := Range(1, 10).OrderBy(func(i int) int {
		return i % 3
	}).ThenByDescending(func(i int) int {
		return i
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [9 6 3 10 7 4 1 8 5 2]
}

func ExampleOrderedQuery_ThenBy() {
	fruits := []string{"grape", "passionfruit", "banana", "mango",
		"orange", "raspberry", "apple", "blueberry"}

	query := FromSlice(fruits).OrderBy(func(f string) int {
		return len(f)
	}).ThenBy(func(f string) string {
		return f
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [apple grape mango banana orange blueberry raspberry passionfruit]
}

func ExampleOrderedQuery_Distinct() {
	query := FromSlice([]int{4, 2, 2, 9, 4, 1}).OrderBy(func(i int) int {
		return i
	}).Distinct()

	fmt.Println(query.ToSlice())

	// Output:
	// [1 2 4 9]
}

func ExampleQuery_Distinct() {
	fmt.Println(FromSlice([]int{1, 2, 2, 3, 1}).Distinct().ToSlice())

	// Output:
	// [1 2 3]
}

func ExampleQuery_DistinctBy() {
	type product struct {
		Name string
		Code int
	}

	products := []product{
		{"apple", 9},
		{"orange", 4},
		{"apple", 9},
		{"lemon", 12},
	}

	query := FromSlice(products).DistinctBy(func(p product) int {
		return p.Code
	}).Select(func(p product) string {
		return p.Name
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [apple orange lemon]
}

func ExampleQuery_Union() {
	q1 := FromSlice([]int{1, 2, 3})
	q2 := FromSlice([]int{2, 4, 5, 1})

	fmt.Println(q1.Union(q2).ToSlice())

	// Output:
	// [1 2 3 4 5]
}

func ExampleQuery_Intersect() {
	q1 := FromSlice([]int{44, 26, 92, 30, 71, 38})
	q2 := FromSlice([]int{39, 59, 83, 47, 26, 4, 30})

	fmt.Println(q1.Intersect(q2).ToSlice())

	// Output:
	// [26 30]
}

func ExampleQuery_Except() {
	q1 := FromSlice([]float64{2.0, 2.1, 2.2, 2.3, 2.4, 2.5})
	q2 := FromSlice([]float64{2.2})

	fmt.Println(q1.Except(q2).ToSlice())

	// Output:
	// [2 2.1 2.3 2.4 2.5]
}

func ExampleQuery_Zip() {
	numbers := []int{1, 2, 3, 4, 5}
	words := []string{"one", "two", "three"}

	query := FromSlice(numbers).Zip(FromSlice(words), func(n int, w string) string {
		return fmt.Sprintf("%d=%s", n, w)
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [1=one 2=two 3=three]
}

func ExampleQuery_Zip3() {
	query := Range(1, 3).Zip3(
		FromSlice([]string{"one", "two", "three"}),
		FromSlice([]string{"I", "II", "III"}),
		func(n int, word, roman string) string {
			return fmt.Sprintf("%d=%s=%s", n, word, roman)
		},
	)

	fmt.Println(query.ToSlice())

	// Output:
	// [1=one=I 2=two=II 3=three=III]
}

func ExampleQuery_Aggregate() {
	fruits := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	longest, _ := FromSlice(fruits).Aggregate(func(longest, next string) string {
		if len(next) > len(longest) {
			return next
		}
		return longest
	})

	fmt.Println(longest)

	// Output:
	// passionfruit
}

func ExampleQuery_AggregateWithSeed() {
	fruits := []string{"apple", "mango", "orange"}

	// The accumulator type may differ from the element type.
	totalLength := FromSlice(fruits).AggregateWithSeed(0, func(acc int, fruit string) int {
		return acc + len(fruit)
	})

	fmt.Println(totalLength)

	// Output:
	// 16
}

func ExampleQuery_AggregateBy() {
	type score struct {
		id    string
		value int
	}
	scores := []score{{"0", 42}, {"1", 5}, {"1", 10}, {"0", 25}}

	totals := FromSlice(scores).AggregateBy(
		func(score score) string { return score.id },
		0,
		func(total int, score score) int { return total + score.value },
	)
	fmt.Println(totals.ToSlice())

	// Output:
	// [{0 67} {1 15}]
}

func ExampleQuery_AggregateByWithSeedSelector() {
	type sale struct {
		region string
		amount int
	}
	sales := []sale{{"EU", 1}, {"APAC", 5}, {"EU", 2}}

	totals := FromSlice(sales).AggregateByWithSeedSelector(
		func(sale sale) string { return sale.region },
		func(region string) int { return len(region) },
		func(total int, sale sale) int { return total + sale.amount },
	)
	fmt.Println(totals.ToSlice())

	// Output:
	// [{EU 5} {APAC 9}]
}

func ExampleQuery_All() {
	pets := map[string]int{"Barley": 10}

	allNamesStartWithB := FromMap(pets).All(func(kv KeyValue[string, int]) bool {
		return strings.HasPrefix(kv.Key, "B")
	})

	fmt.Println(allNamesStartWithB)

	// Output:
	// true
}

func ExampleQuery_AnyWith() {
	fmt.Println(Range(1, 10).AnyWith(func(i int) bool {
		return i > 9
	}))

	// Output:
	// true
}

func ExampleQuery_CountWith() {
	fmt.Println(Range(1, 10).CountWith(func(i int) bool {
		return i%2 == 0
	}))

	// Output:
	// 5
}

func ExampleQuery_CountBy() {
	counts := FromString("abracadabra").CountBy(func(r rune) string { return string(r) })
	fmt.Println(counts.ToSlice())

	// Output:
	// [{a 5} {b 2} {r 2} {c 1} {d 1}]
}

func ExampleQuery_First() {
	first, ok := FromSlice([]int{10, 20, 30}).First()
	fmt.Println(first, ok)

	_, ok = FromSlice([]int{}).First()
	fmt.Println(ok)

	// Output:
	// 10 true
	// false
}

func ExampleQuery_FirstWith() {
	first, _ := Range(1, 10).FirstWith(func(i int) bool {
		return i%3 == 0
	})

	fmt.Println(first)

	// Output:
	// 3
}

func ExampleQuery_Single() {
	value, ok := FromSlice([]int{42}).Single()
	fmt.Println(value, ok)

	_, ok = FromSlice([]int{1, 2}).Single()
	fmt.Println(ok)

	// Output:
	// 42 true
	// false
}

func ExampleQuery_Take() {
	fmt.Println(Range(1, 10).Take(3).ToSlice())

	// Output:
	// [1 2 3]
}

func ExampleQuery_TakeRange() {
	fmt.Println(Range(0, 10).TakeRange(PositionFromEnd(4), PositionFromEnd(1)).ToSlice())

	// Output:
	// [6 7 8]
}

func ExampleQuery_TakeWhile() {
	fruits := []string{"apple", "banana", "mango", "orange", "passionfruit", "grape"}

	query := FromSlice(fruits).TakeWhile(func(fruit string) bool {
		return fruit != "orange"
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [apple banana mango]
}

func ExampleQuery_Skip() {
	fmt.Println(Range(1, 10).Skip(7).ToSlice())

	// Output:
	// [8 9 10]
}

func ExampleQuery_Reverse() {
	fmt.Println(FromString("gnirts").Reverse().Select(func(r rune) string {
		return string(r)
	}).ToSlice())

	// Output:
	// [s t r i n g]
}

func ExampleQuery_Concat() {
	q := FromSlice([]int{1, 2}).Concat(FromSlice([]int{3, 4}))
	fmt.Println(q.ToSlice())

	// Output:
	// [1 2 3 4]
}

func ExampleQuery_Append() {
	fmt.Println(FromSlice([]int{1, 2}).Append(3).ToSlice())

	// Output:
	// [1 2 3]
}

func ExampleQuery_Prepend() {
	fmt.Println(FromSlice([]int{1, 2}).Prepend(0).ToSlice())

	// Output:
	// [0 1 2]
}

func ExampleQuery_DefaultIfEmpty() {
	fmt.Println(FromSlice([]int{}).DefaultIfEmpty(-1).ToSlice())

	// Output:
	// [-1]
}

func ExampleQuery_IndexOf() {
	fmt.Println(FromSlice([]string{"a", "b", "c"}).IndexOf(func(s string) bool {
		return s == "b"
	}))

	// Output:
	// 1
}

func ExampleQuery_ElementAt() {
	value, ok := FromSlice([]string{"a", "b", "c"}).ElementAt(PositionFromEnd(1))
	fmt.Println(value, ok)

	// Output:
	// c true
}

func ExampleIndex() {
	fmt.Println(Index(FromSlice([]string{"zero", "one"})).ToSlice())

	// Output:
	// [{0 zero} {1 one}]
}

func ExampleChunk() {
	fmt.Println(Chunk(Range(1, 5), 2).ToSlice())

	// Output:
	// [[1 2] [3 4] [5]]
}

func ExampleQuery_ToMapBy() {
	fruits := []string{"apple", "banana", "fig"}

	m := FromSlice(fruits).ToMapBy(
		func(f string) string { return f },
		func(f string) int { return len(f) },
	)

	fmt.Println(m["banana"])

	// Output:
	// 6
}

func ExampleToMap() {
	pairs := FromSlice([]string{"apple", "fig"}).Select(func(f string) KeyValue[string, int] {
		return KeyValue[string, int]{Key: f, Value: len(f)}
	})

	m := ToMap(pairs)
	fmt.Println(m["fig"])

	// Output:
	// 3
}

func ExampleQuery_ToChannel() {
	ch := make(chan int)
	go Range(1, 3).ToChannel(ch)

	for v := range ch {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleSum() {
	fmt.Println(Sum(Range(1, 10)))

	// Output:
	// 55
}

func ExampleAverage() {
	fmt.Println(Average(FromSlice([]float64{1.5, 2.5, 5.0})))

	// Output:
	// 3
}

func ExampleMin() {
	minVal, _ := Min(FromSlice([]int{3, 1, 4, 1, 5}))
	fmt.Println(minVal)

	// Output:
	// 1
}

func ExampleMax() {
	maxVal, _ := Max(FromSlice([]int{3, 1, 4, 1, 5}))
	fmt.Println(maxVal)

	// Output:
	// 5
}

func ExampleQuery_SumBy() {
	fruits := []string{"apple", "banana", "fig"}

	fmt.Println(FromSlice(fruits).SumBy(func(f string) int {
		return len(f)
	}))

	// Output:
	// 14
}

func ExampleQuery_MaxBy() {
	fruits := []string{"apple", "banana", "fig"}

	longest, _ := FromSlice(fruits).MaxBy(func(f string) int {
		return len(f)
	})

	fmt.Println(longest)

	// Output:
	// banana
}

func ExampleQuery_MinBy() {
	fruits := []string{"apple", "banana", "fig"}

	shortest, _ := FromSlice(fruits).MinBy(func(f string) int {
		return len(f)
	})

	fmt.Println(shortest)

	// Output:
	// fig
}

func ExampleQuery_Contains() {
	fmt.Println(Range(1, 10).Contains(3))

	// Output:
	// true
}

func ExampleQuery_SequenceEqual() {
	q1 := FromSlice([]int{1, 2, 3})
	q2 := Range(1, 3)

	fmt.Println(q1.SequenceEqual(q2))

	// Output:
	// true
}

func ExampleQuery_Sort() {
	type player struct {
		Name string
		Wins int
	}

	players := []player{{"mike", 3}, {"anna", 9}, {"joe", 5}}

	query := FromSlice(players).Sort(func(a, b player) bool {
		return a.Wins > b.Wins
	}).Select(func(p player) string {
		return p.Name
	})

	fmt.Println(query.ToSlice())

	// Output:
	// [anna joe mike]
}

func ExampleQuery_ForEachIndexed() {
	FromSlice([]string{"a", "b"}).ForEachIndexed(func(i int, s string) {
		fmt.Println(i, s)
	})

	// Output:
	// 0 a
	// 1 b
}
