package linq

import "fmt"

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
	From(input).ToMap(&m)
	fmt.Println(m)

	// Output:
	// map[10:true]
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

func ExampleQuery_Aggregate() {
	input := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	result := From(input).Aggregate(func(r interface{}, i interface{}) interface{} {
		if len(r.(string)) > len(i.(string)) {
			return r
		}
		return i
	})

	fmt.Println(result)

	// Output:
	// passionfruit
}

func ExampleQuery_Concat() {
	q := From([]int{1, 2, 3}).Concat(From([]int{4, 5, 6}))
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

func ExampleQuery_GroupJoin() {
	fruits := []string{
		"apple",
		"banana",
		"apricot",
		"cherry",
		"clementine",
	}

	q := FromString("abc").GroupJoin(
		From(fruits),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return []rune(i.(string))[0] },
		func(outer interface{}, inners []interface{}) interface{} {
			return KeyValue{string(outer.(rune)), inners}
		})

	fmt.Println(q.Results())

	// Output:
	// [{a [apple apricot]} {b [banana]} {c [cherry clementine]}]
}

func ExampleQuery_Join() {
	fruits := []string{
		"apple",
		"banana",
		"apricot",
		"cherry",
		"clementine",
	}

	q := Range(1, 10).Join(From(fruits),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return len(i.(string)) },
		func(outer interface{}, inner interface{}) interface{} {
			return KeyValue{outer, inner}
		})

	fmt.Println(q.Results())

	// Output:
	// [{5 apple} {6 banana} {6 cherry} {7 apricot} {10 clementine}]
}

func ExampleQuery_OrderBy() {
	q := Range(1, 10).OrderBy(func(i interface{}) interface{} {
		return i.(int) % 2
	}).ThenByDescending(func(i interface{}) interface{} {
		return i
	})

	fmt.Println(q.Results())

	// Output:
	// [10 8 6 4 2 9 7 5 3 1]
}

func ExampleQuery_SelectMany() {
	input := [][]int{{1, 2, 3}, {4, 5, 6, 7}}

	q := From(input).SelectMany(func(i interface{}) Query {
		return From(i)
	})

	fmt.Println(q.Results())

	// Output:
	// [1 2 3 4 5 6 7]
}

func ExampleQuery_Union() {
	q := Range(1, 10).Union(Range(6, 10))

	fmt.Println(q.Results())

	// Output:
	// [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15]
}

func ExampleQuery_Zip() {
	number := []int{1, 2, 3, 4, 5}
	words := []string{"one", "two", "three"}

	q := From(number).Zip(From(words), func(a interface{}, b interface{}) interface{} {
		return []interface{}{a, b}
	})

	fmt.Println(q.Results())

	// Output:
	// [[1 one] [2 two] [3 three]]
}

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

func ExampleQuery_ToMapBy() {
	input := [][]interface{}{{1, true}}

	result := make(map[int]bool)
	From(input).ToMapBy(&result,
		func(i interface{}) interface{} {
			return i.([]interface{})[0]
		},
		func(i interface{}) interface{} {
			return i.([]interface{})[1]
		})

	fmt.Println(result)

	// Output:
	// map[1:true]
}

func ExampleQuery_ToSlice() {
	result := []int{}
	Range(1, 10).ToSlice(&result)

	fmt.Println(result)

	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
}
