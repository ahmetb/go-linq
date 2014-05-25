package linq

import (
	"errors"
	"fmt"
	"github.com/jacobsa/oglematchers"
	c "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

type foo struct {
	str string
	num int
}

var (
	empty []interface{}
	arr0  = []int{1, 2, 3, 1, 2}
	arr1  = []string{"foo", "bar", "baz"}
	arr2  = []T{nil, "foo", 3.14, true, false}
	arr3  = []foo{foo{"A", 0}, foo{"B", 1}, foo{"C", -1}}
	arr4  = []T{&foo{"C", 0xffff}, nil, &foo{"D", 0x7fff}, byte(12), nil}
)

var (
	maxRandDelayMs = 100
	errFoo         = errors.New("")
	alwaysTrue     = func(i T) (bool, error) {
		return true, nil
	}
	alwaysFalse = func(i T) (bool, error) {
		return false, nil
	}
	erroneusBinaryFunc = func(i T) (bool, error) {
		return true, errFoo
	}
)

func shouldSlicesResemble(actual interface{}, expected ...interface{}) string {
	expectedSlice, ok := takeSliceArg(expected[0])
	if !ok {
		return "Cannot cast expected slice to []T"
	}
	actualSlice, ok := takeSliceArg(actual)
	if !ok {
		return "Cannot cast actual slice to []T"
	}

	if len(expectedSlice) != len(actualSlice) {
		return fmt.Sprintf("Expected: '%v'\nActual:   '%v'\n(Should resemble: slices have different lengths.)", expectedSlice, actualSlice)
	}
	for i := 0; i < len(expectedSlice); i++ {
		if matchError := oglematchers.DeepEquals(expectedSlice[i]).Matches(actualSlice[i]); matchError != nil {
			return fmt.Sprintf("Expected: '%v'\nActual:   '%v'\n(Element[%v] Should be equal: %v)", expectedSlice, actualSlice, i, matchError)
		}
	}
	return ""
}

func TestFrom(t *testing.T) {
	c.Convey("When passed nil value, error returned", t, func() {
		c.So(From(nil).err, c.ShouldNotEqual, nil)
	})

	c.Convey("When passed non-slice value, error returned", t, func() {
		var t, u, v T
		t = "ahoy!"
		u = foo{"A", 0}
		v = byte(12)
		c.So(From(t).err, c.ShouldEqual, ErrInvalidInput)
		c.So(From(u).err, c.ShouldEqual, ErrInvalidInput)
		c.So(From(v).err, c.ShouldEqual, ErrInvalidInput)
	})

	c.Convey("When passed non-nil value, structure should have the exact same slice at different location", t, func() {
		c.Convey("Empty array", func() {
			r := From(empty).values
			c.So(r, shouldSlicesResemble, empty)
			c.So(r, c.ShouldNotEqual, empty) // slice copied?
		})
		c.Convey("Non-empty arrays", func() {
			c.Convey("Passed & held slices are different", func() {
				c.So(From(arr0).values, c.ShouldNotEqual, arr0)
				c.So(From(arr4).values, c.ShouldNotEqual, arr4)
			})

			c.Convey("Deep slice equality", func() {
				c.So(From(arr0).values, shouldSlicesResemble, arr0)
				c.So(From(arr1).values, shouldSlicesResemble, arr1)
				c.So(From(arr2).values, shouldSlicesResemble, arr2)
				c.So(From(arr3).values, shouldSlicesResemble, arr3)
				c.So(From(arr4).values, shouldSlicesResemble, arr4)
			})
		})
	})
}

func TestResults(t *testing.T) {
	c.Convey("If error exists in given queryable, error is returned", t, func() {
		q := Query{
			values: nil,
			err:    errFoo}
		_, err := q.Results()
		c.So(err, c.ShouldEqual, errFoo)
	})
	c.Convey("Given no errors exist, non-nil results are returned", t, func() {
		q := From(arr0)
		val, err := q.Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(val, shouldSlicesResemble, arr0)
	})
	c.Convey("Returned result is isolated (copied) from original query source", t, func() {
		// Regression for BUG: modifying result slice effects subsequent query methods
		arr := []int{1, 2, 3, 4, 5}
		q := From(arr)
		res, _ := q.Results()
		res[0] = 100
		sum, _ := q.Sum()
		c.So(sum, c.ShouldEqual, 15)
	})
}

func TestWhere(t *testing.T) {
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
	})

	c.Convey("An error returned from f is reflected on the result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Where(alwaysTrue).Results()
		c.So(err, c.ShouldNotEqual, nil)

		c.Convey("Chain successful and erroneus functions", func() {
			_, err := From(arr0).Where(alwaysTrue).Where(erroneusBinaryFunc).Results()
			c.So(err, c.ShouldNotEqual, nil)

			c.Convey("Erroneus function is in chain but not called", func() {
				_, err := From(arr0).Where(alwaysTrue).Where(alwaysFalse).Where(erroneusBinaryFunc).Results()
				c.So(err, c.ShouldEqual, nil)
			})
		})
	})

	c.Convey("Chose none of the elements", t, func() {
		val, _ := From(arr0).Where(alwaysFalse).Results()
		c.So(len(val), c.ShouldEqual, 0)
	})

	c.Convey("Chose all elements, as is", t, func() {
		val, _ := From(arr0).Where(alwaysTrue).Results()
		c.So(val, shouldSlicesResemble, arr0)
	})

	c.Convey("Basic filtering (x mod 2)==0", t, func() {
		n := 100
		divisibleBy2 := func(i T) (bool, error) {
			return i.(int)%2 == 0, nil
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		val, err := From(arr).Where(divisibleBy2).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(len(val), c.ShouldEqual, n/2)
	})
}

func TestSelect(t *testing.T) {
	asIs := func(i T) (T, error) {
		return i, nil
	}
	erroneusFunc := func(i T) (T, error) {
		return nil, errFoo
	}

	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Select(asIs).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})

	c.Convey("Nil func returns error", t, func() {
		_, err := From(arr0).Select(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
	})

	c.Convey("Error returned from provided func", t, func() {
		val, err := From(arr0).Select(erroneusFunc).Results()
		c.So(err, c.ShouldNotEqual, nil)

		c.Convey("Erroneus function is in chain with as-is select", func() {
			_, err = From(arr0).Select(asIs).Select(erroneusFunc).Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Erroneus function is in chain but not called", func() {
			val, err = From(arr0).Where(alwaysFalse).Select(erroneusFunc).Results()
			c.So(err, c.ShouldEqual, nil)
			c.So(len(val), c.ShouldEqual, 0)
		})

	})

	c.Convey("Select all elements as is", t, func() {
		val, err := From(arr0).Select(asIs).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(val, shouldSlicesResemble, arr0)
	})

	c.Convey("Pow(x,2) for i in []int", t, func() {
		pow := func(i T) (T, error) {
			return i.(int) * i.(int), nil
		}
		val, err := From(arr0).Select(pow).Results()
		c.So(err, c.ShouldEqual, nil)
		arr := make([]int, len(arr0))
		for j, i := range arr0 {
			arr[j] = i * i
		}
		res := make([]int, len(val))
		for j, v := range val {
			res[j] = v.(int)
		}
		c.So(res, shouldSlicesResemble, arr)
	})
}

func TestDistinct(t *testing.T) {
	c.Convey("Empty slice", t, func() {
		res, err := From(empty).Distinct().Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})

	allSameInt := []int{1, 1, 1, 1, 1, 1, 1, 1, 1}
	allSameStruct := []foo{foo{"A", -1}, foo{"A", -1}, foo{"A", -1}}
	allNil := []T{nil, nil, nil, nil, nil, nil, nil, nil, nil}

	c.Convey("With default equality comparer ==", t, func() {
		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).Distinct().Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("All elements are the same", func() {
			res, _ := From(allSameInt).Distinct().Results()
			c.So(res, shouldSlicesResemble, []int{allSameInt[0]})

			c.Convey("All elements are nil", func() {
				res, _ = From(allNil).Distinct().Results()
				c.So(res, shouldSlicesResemble, []T{allNil[0]})
			})
		})
		c.Convey("Distinct on structs and nils", func() {
			arr := []T{foo{"A", 0xffff}, nil, foo{"B", 0x7fff}, nil, foo{"A", 0xffff}}
			res, _ := From(arr).Distinct().Results()
			c.So(len(res), c.ShouldEqual, 3)
		})
		c.Convey("Randomly generated integers with duplicates or more", func() {
			var arr = make([]int, 10000)
			var dict = make(map[int]bool, len(arr))

			rand.Seed(time.Now().UnixNano())
			for i := 0; i < len(arr); i++ {
				r := rand.Intn(len(arr) * 4 / 5) // collision 20%
				arr[i] = r
				dict[r] = true
			}
			res, _ := From(arr).Distinct().Results()
			c.So(len(res), c.ShouldEqual, len(dict))
		})
	})

	c.Convey("With provided equality comparer", t, func() {
		fooComparer := func(i T, j T) (bool, error) {
			return i.(foo).str == j.(foo).str && i.(foo).num == j.(foo).num, nil
		}
		fooPtrComparer := func(i T, j T) (bool, error) {
			return i.(*foo).str == j.(*foo).str && i.(*foo).num == j.(*foo).num, nil
		}
		_ = fooPtrComparer

		erroneusComparer := func(i T, j T) (bool, error) {
			return false, errFoo
		}

		c.Convey("Previous error is reflected on result", func() {
			_, err := From(allSameStruct).Where(erroneusBinaryFunc).DistinctBy(fooComparer).Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Provided func is nil", func() {
			_, err := From(allSameStruct).DistinctBy(nil).Results()
			c.So(err, c.ShouldEqual, ErrNilFunc)
		})
		c.Convey("Comparer returns error", func() {
			_, err := From(arr0).DistinctBy(erroneusComparer).Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("All elements are the same", func() {
			res, _ := From(allSameStruct).DistinctBy(fooComparer).Results()
			c.So(res, shouldSlicesResemble, []foo{allSameStruct[0]})
		})
		c.Convey("All elements are distinct", func() {
			var arr = make([]int, 100)
			for i := 0; i < len(arr); i++ {
				arr[i] = i
			}
			res, _ := From(arr).DistinctBy(func(this T, that T) (bool, error) {
				return this.(int) == that.(int), nil
			}).Results()
			c.So(res, shouldSlicesResemble, arr)
		})
		c.Convey("Ensure leftmost appearance is returned in multiple occurrence cases", func() {
			arr := []*foo{&foo{"A", 0}, &foo{"B", 0}, &foo{"A", 0}, &foo{"C", 0},
				&foo{"A", 0}, &foo{"B", 0}}
			res, _ := From(arr).DistinctBy(fooPtrComparer).Results()
			c.So(len(res), c.ShouldEqual, 3)
			c.So(res[0], c.ShouldEqual, arr[0]) // A
			c.So(res[1], c.ShouldEqual, arr[1]) // B
			c.So(res[2], c.ShouldEqual, arr[3]) // C
		})
		c.Convey("Randomly generated integers with likely collisions", func() {
			var arr = make([]int, 10000)
			var dict = make(map[int]bool, len(arr))
			rand.Seed(time.Now().UnixNano())
			for i := 0; i < len(arr); i++ {
				r := rand.Intn(len(arr) * 4 / 5) // collision 20%
				arr[i] = r
				dict[r] = true
			}
			res, _ := From(arr).DistinctBy(func(this T, that T) (bool, error) {
				return this.(int) == that.(int), nil
			}).Results()
			c.So(len(res), c.ShouldEqual, len(dict))
		})
	})
}

func TestUnion(t *testing.T) {
	uniqueArr0 := []int{1, 2, 3, 4, 5}
	uniqueArr1 := []string{"a", "b", "c"}
	allSameArr := []uint{1, 1, 1, 1}
	sameStruct0 := []foo{foo{"A", 0}, foo{"B", 0}}
	sameStruct1 := []foo{foo{"B", 0}, foo{"A", 0}}
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(uniqueArr0).Where(erroneusBinaryFunc).Union(uniqueArr0).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Passed non-slice value, error returned", t, func() {
		_, err := From(empty).Union("someString").Results()
		c.So(err, c.ShouldEqual, ErrInvalidInput)
	})
	c.Convey("Empty ∪ nil", t, func() {
		_, err := From(empty).Union(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilInput)
	})
	c.Convey("Empty ∪ empty", t, func() {
		res, _ := From(empty).Union(empty).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Empty ∪ non-empty", t, func() {
		res, _ := From(empty).Union(uniqueArr0).Results()
		c.So(res, shouldSlicesResemble, uniqueArr0)
	})
	c.Convey("Non-empty ∪ empty", t, func() {
		res, _ := From(uniqueArr0).Union(empty).Results()
		c.So(res, shouldSlicesResemble, uniqueArr0)
	})
	c.Convey("(Unique slice) ∪ (itself)", t, func() {
		res, _ := From(uniqueArr0).Union(uniqueArr0).Results()
		c.So(res, shouldSlicesResemble, uniqueArr0)
	})
	c.Convey("(All same slice) ∪ (itself)", t, func() {
		res, _ := From(allSameArr).Union(allSameArr).Results()
		c.So(len(res), c.ShouldEqual, 1)
	})
	c.Convey("Mixed types", t, func() {
		res, _ := From(uniqueArr0).Union(uniqueArr1).Results()
		c.So(len(res), c.ShouldEqual, len(uniqueArr0)+len(uniqueArr1))
	})
	c.Convey("Same-type structs", t, func() {
		res, _ := From(sameStruct0).Union(sameStruct1).Results()
		c.So(len(res), c.ShouldEqual, len(sameStruct1))
	})
}

func TestIntersect(t *testing.T) {
	uniqueArr := []int{1, 2, 3, 4, 5}
	allSameArr := []int{1, 1, 1, 1}
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(uniqueArr).Where(erroneusBinaryFunc).Intersect(uniqueArr).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Passed non-slice value, error returned", t, func() {
		_, err := From(empty).Intersect("someString").Results()
		c.So(err, c.ShouldEqual, ErrInvalidInput)
	})
	c.Convey("Empty ∩ nil", t, func() {
		_, err := From(empty).Intersect(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilInput)
	})
	c.Convey("Empty ∩ empty", t, func() {
		res, _ := From(empty).Intersect(empty).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Empty ∩ non-empty", t, func() {
		res, _ := From(empty).Intersect(uniqueArr).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Non-empty ∩ empty", t, func() {
		res, _ := From(uniqueArr).Intersect(empty).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("(Unique set) ∩ (itself)", t, func() {
		res, _ := From(uniqueArr).Intersect(uniqueArr).Results()
		c.So(res, shouldSlicesResemble, uniqueArr)
	})
	c.Convey("(All same slice) ∩ (itself)", t, func() {
		res, _ := From(allSameArr).Intersect(allSameArr).Results()
		c.So(len(res), c.ShouldEqual, 1)
	})
	c.Convey("There is some intersection", t, func() {
		res, _ := From([]T{1, 2, 3, 4, 5}).Intersect([]T{3, 4, 5, 6, 7}).Results()
		c.So(res, shouldSlicesResemble, []T{3, 4, 5})
	})
}

func TestExcept(t *testing.T) {
	uniqueArr := []int{1, 2, 3, 4, 5}
	allSameArr := []int{1, 1, 1, 1}
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(uniqueArr).Where(erroneusBinaryFunc).Except(uniqueArr).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Passed non-slice value, error returned", t, func() {
		_, err := From(empty).Except("someString").Results()
		c.So(err, c.ShouldEqual, ErrInvalidInput)
	})
	c.Convey("Empty ∖ nil", t, func() {
		_, err := From(empty).Except(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilInput)
	})
	c.Convey("Empty ∖ empty", t, func() {
		res, _ := From(empty).Except(empty).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Empty ∖ non-empty", t, func() {
		res, _ := From(empty).Except(uniqueArr).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Non-empty ∖ empty", t, func() {
		res, _ := From(uniqueArr).Except(empty).Results()
		c.So(res, shouldSlicesResemble, uniqueArr)
	})
	c.Convey("(Unique set) ∖ (itself)", t, func() {
		res, _ := From(uniqueArr).Except(uniqueArr).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("(All same slice) ∖ (itself)", t, func() {
		res, _ := From(allSameArr).Except(allSameArr).Results()
		c.So(len(res), c.ShouldEqual, 0)
	})
	c.Convey("There is some intersection", t, func() {
		res, _ := From([]int{1, 2, 3, 4, 5}).Except([]int{3, 4, 5, 6, 7}).Results()
		c.So(res, shouldSlicesResemble, []int{1, 2})
	})
}

func TestCount(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).CountBy(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).CountBy(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).CountBy(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).CountBy(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("No matches", t, func() {
		cnt, _ := From(arr0).CountBy(alwaysFalse)
		c.So(cnt, c.ShouldEqual, 0)
		cnt, _ = From(arr0).Where(alwaysFalse).Count()
		c.So(cnt, c.ShouldEqual, 0)
	})
	c.Convey("All matches", t, func() {
		cnt, _ := From(arr0).CountBy(alwaysTrue)
		c.So(cnt, c.ShouldEqual, len(arr0))
		cnt, _ = From(arr0).Count()
		c.So(cnt, c.ShouldEqual, len(arr0))
	})
}

func TestAny(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AnyWith(alwaysTrue)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AnyWith(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AnyWith(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).AnyWith(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("No matches", t, func() {
		r, _ := From(arr0).AnyWith(alwaysFalse)
		c.So(r, c.ShouldEqual, false)
		r, _ = From(arr0).Where(alwaysFalse).Any()
		c.So(r, c.ShouldEqual, false)
	})
	c.Convey("All matches", t, func() {
		r, _ := From(arr0).AnyWith(alwaysTrue)
		c.So(r, c.ShouldEqual, true)
		r, _ = From(arr0).Where(alwaysTrue).Any()
		c.So(r, c.ShouldEqual, true)
	})
}

func TestSingle(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Single(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Single(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Single(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		c.So(err, c.ShouldNotEqual, ErrNotSingle)
		_, err = From(arr0).Where(alwaysFalse).Single(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
	c.Convey("No matches", t, func() {
		_, err := From(arr0).Single(alwaysFalse)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
	c.Convey("All matches", t, func() {
		_, err := From(arr0).Single(alwaysTrue)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
	c.Convey("Only one match", t, func() {
		match := 0
		var match0 = func(i T) (bool, error) {
			return i.(int) == match, nil
		}
		r, _ := From([]int{-1, -1, 0, 1, 1}).Single(match0)
		c.So(r, c.ShouldEqual, match)
		_, err := From([]int{0, 1, 2, 2, 0}).Single(match0)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
}

func TestAll(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).All(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).All(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).All(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).All(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("Empty slice", t, func() {
		r, _ := From(empty).All(alwaysTrue)
		c.So(r, c.ShouldEqual, true)
	})
	c.Convey("No matches", t, func() {
		r, _ := From(arr0).All(alwaysFalse)
		c.So(r, c.ShouldEqual, false)
	})
	c.Convey("All matches", t, func() {
		r, _ := From(arr0).All(alwaysTrue)
		c.So(r, c.ShouldEqual, true)
	})
	c.Convey("Multiple matches", t, func() {
		match0 := func(i T) (bool, error) {
			return i.(int) == 0, nil
		}
		r, _ := From([]int{0, 1, 2, 2, 0}).All(match0)
		c.So(r, c.ShouldEqual, false)
	})
}

func TestElementAt(t *testing.T) {
	intArr := []int{1, 2, 3, 4, 5}
	c.Convey("empty.ElementAt(1) is not found", t, func() {
		_, ok, err := From(empty).ElementAt(1)
		c.So(ok, c.ShouldBeFalse)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("negative index returns is ErrNegativeParam", t, func() {
		_, _, err := From(empty).ElementAt(-1)
		c.So(err, c.ShouldEqual, ErrNegativeParam)
	})
	c.Convey("first element is returned", t, func() {
		v, ok, _ := From(intArr).ElementAt(0)
		c.So(ok, c.ShouldBeTrue)
		c.So(v, c.ShouldEqual, intArr[0])
	})
	c.Convey("last element is returned", t, func() {
		v, ok, _ := From(intArr).ElementAt(len(intArr) - 1)
		c.So(ok, c.ShouldBeTrue)
		c.So(v, c.ShouldEqual, intArr[len(intArr)-1])
	})
	c.Convey("out of index returns not found on non-empty slice", t, func() {
		_, ok, err := From(intArr).ElementAt(len(intArr))
		c.So(ok, c.ShouldBeFalse)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("previous errors are reflected", t, func() {
		_, _, err := From(arr0).Where(erroneusBinaryFunc).ElementAt(0)
		c.So(err, c.ShouldNotEqual, nil)
	})
}

func TestFirst(t *testing.T) {
	c.Convey("empty.First is not found", t, func() {
		_, ok, err := From(empty).First()
		c.So(err, c.ShouldEqual, nil)
		c.So(ok, c.ShouldBeFalse)
	})
	c.Convey("first element is returned", t, func() {
		v, ok, _ := From(arr3).First()
		c.So(ok, c.ShouldBeTrue)
		c.So(v, c.ShouldResemble, arr3[0])
	})
	c.Convey("previous errors are reflected", t, func() {
		_, _, err1 := From(arr0).Where(erroneusBinaryFunc).First()
		c.So(err1, c.ShouldNotEqual, nil)
	})
}

func TestFirstBy(t *testing.T) {
	c.Convey("previous errors are reflected", t, func() {
		_, _, err1 := From(arr0).Where(erroneusBinaryFunc).FirstBy(alwaysTrue)
		c.So(err1, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, _, err1 := From(arr0).FirstBy(nil)
		_, _, err2 := From(arr0).FirstBy(nil)
		c.So(err1, c.ShouldEqual, ErrNilFunc)
		c.So(err2, c.ShouldEqual, ErrNilFunc)
	})
	c.Convey("erroneus function reflected on result", t, func() {
		_, _, err1 := From(arr0).FirstBy(erroneusBinaryFunc)
		_, _, err2 := From(arr0).FirstBy(erroneusBinaryFunc)
		c.So(err1, c.ShouldNotEqual, nil)
		c.So(err2, c.ShouldNotEqual, nil)
	})
	c.Convey("empty.FirstBy is not found", t, func() {
		_, ok, err1 := From(empty).FirstBy(alwaysFalse)
		c.So(err1, c.ShouldEqual, nil)
		c.So(ok, c.ShouldBeFalse)
	})
	c.Convey("Actual first element is returned", t, func() {
		val, ok, _ := From(arr3).FirstBy(alwaysTrue)
		c.So(ok, c.ShouldBeTrue)
		c.So(val, c.ShouldResemble, arr3[0])
	})
	c.Convey("No matches", t, func() {
		_, ok, err := From(arr3).FirstBy(alwaysFalse)
		c.So(ok, c.ShouldBeFalse)
		c.So(err, c.ShouldEqual, nil)
	})
}

func TestLast(t *testing.T) {
	c.Convey("empty.Last is not found", t, func() {
		_, ok, err := From(empty).Last()
		c.So(ok, c.ShouldBeFalse)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("Last element is returned", t, func() {
		v, ok, _ := From(arr3).Last()
		c.So(ok, c.ShouldBeTrue)
		c.So(v, c.ShouldResemble, arr3[len(arr3)-1])
	})
	c.Convey("previous errors are reflected", t, func() {
		_, _, err1 := From(arr0).Where(erroneusBinaryFunc).Last()
		c.So(err1, c.ShouldNotEqual, nil)
	})
}

func TestLastBy(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, _, err := From(arr0).Where(erroneusBinaryFunc).LastBy(alwaysTrue)
		c.So(err, c.ShouldNotEqual, nil)
	})

	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, _, err1 := From(arr0).LastBy(nil)
		_, _, err2 := From(arr0).LastBy(nil)
		c.So(err1, c.ShouldEqual, ErrNilFunc)
		c.So(err2, c.ShouldEqual, ErrNilFunc)
	})

	c.Convey("erroneus function reflected on result", t, func() {
		_, _, err1 := From(arr0).LastBy(erroneusBinaryFunc)
		_, _, err2 := From(arr0).LastBy(erroneusBinaryFunc)
		c.So(err1, c.ShouldNotEqual, nil)
		c.So(err2, c.ShouldNotEqual, nil)
	})
	c.Convey("empty.LastBy is not found", t, func() {
		_, ok, err := From(empty).LastBy(alwaysFalse)
		c.So(ok, c.ShouldBeFalse)
		c.So(err, c.ShouldEqual, nil)
	})

	c.Convey("Actual last element is returned", t, func() {
		val1, ok, _ := From(arr3).LastBy(alwaysTrue)
		c.So(ok, c.ShouldBeTrue)
		c.So(val1, c.ShouldResemble, arr3[len(arr3)-1])
	})
	c.Convey("No matches", t, func() {
		_, ok, err := From(arr3).LastBy(alwaysFalse)
		c.So(ok, c.ShouldBeFalse)
		c.So(err, c.ShouldEqual, nil)
	})

}

func TestReverse(t *testing.T) {
	c.Convey("Previous errors are returned", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Reverse().Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Reversing empty", t, func() {
		res, err := From(empty).Reverse().Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Actual reverse", t, func() {
		arr := []int{1, 2, 3, 4, 5}
		rev := []int{5, 4, 3, 2, 1}
		res, _ := From(arr).Reverse().Results()
		c.So(res, shouldSlicesResemble, rev)

		c.Convey("Slice containing nils", func() {
			arr := []T{1, nil, nil, 2, nil, 3, nil}
			rev := []T{nil, 3, nil, 2, nil, nil, 1}
			res, _ := From(arr).Reverse().Results()
			c.So(res, shouldSlicesResemble, rev)
		})
	})
}

func TestTake(t *testing.T) {
	c.Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Take(1).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Empty slice take n>0", t, func() {
		res, err := From(empty).Take(1).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Take 0", t, func() {
		res, _ := From(arr0).Take(0).Results()
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Take n < 0", t, func() {
		res, err := From(arr0).Take(-1).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Take n > 0", t, func() {
		in := []int{1, 2, 3, 4, 5}
		res, _ := From(in).Take(3).Results()
		c.So(res, shouldSlicesResemble, []int{1, 2, 3})
		c.Convey("Take n ≥ len(arr)", func() {
			res, _ := From(in).Take(len(in)).Results()
			c.So(res, shouldSlicesResemble, res)
			res, _ = From(in).Take(len(in) + 1).Results()
			c.So(res, shouldSlicesResemble, res)
		})
	})
}

func TestTakeWhile(t *testing.T) {
	c.Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).TakeWhile(alwaysTrue).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Nil func passed", t, func() {
		_, err := From(arr0).TakeWhile(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
	})
	c.Convey("Error returned from passed func is reflected in result", t, func() {
		_, err := From(arr0).TakeWhile(erroneusBinaryFunc).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Empty slice take all", t, func() {
		res, err := From(empty).TakeWhile(alwaysTrue).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Take none", t, func() {
		res, _ := From(arr0).TakeWhile(alwaysFalse).Results()
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Take only first", t, func() {
		in := []int{1, 2, 3, 4, 5}
		res, err := From(in).TakeWhile(func(i T) (bool, error) { return i.(int) < 2, nil }).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, in[:1])
	})
}

func TestSkip(t *testing.T) {
	c.Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Skip(1).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Empty slice Skip n>0", t, func() {
		res, err := From(empty).Skip(1).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Skip 0", t, func() {
		res, _ := From(arr0).Skip(0).Results()
		c.So(res, shouldSlicesResemble, arr0)
	})

	c.Convey("Skip n < 0", t, func() {
		res, err := From(arr0).Skip(-1).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, arr0)
	})

	c.Convey("Skip n > 0", t, func() {
		in := []int{1, 2, 3, 4, 5}
		res, _ := From(in).Skip(3).Results()
		c.So(res, shouldSlicesResemble, []int{4, 5})
		c.Convey("Skip n ≥ len(arr)", func() {
			res, _ := From(in).Skip(len(in)).Results()
			c.So(res, shouldSlicesResemble, empty)
			res, _ = From(in).Skip(len(in) + 1).Results()
			c.So(res, shouldSlicesResemble, empty)
		})
	})

	c.Convey("Skip & take & skip", t, func() {
		in := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		res, _ := From(in).Skip(0).Skip(-1000).Skip(1).Take(1000).Take(5).Results()
		c.So(res, shouldSlicesResemble, []int{1, 2, 3, 4, 5})
	})
}

func TestSkipWhile(t *testing.T) {
	c.Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).SkipWhile(alwaysTrue).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Nil func passed", t, func() {
		_, err := From(arr0).SkipWhile(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
	})
	c.Convey("Error returned from passed func is reflected in result", t, func() {
		_, err := From(arr0).SkipWhile(erroneusBinaryFunc).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Empty slice Skip all", t, func() {
		res, err := From(empty).SkipWhile(alwaysTrue).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Skip none", t, func() {
		res, _ := From(arr0).SkipWhile(alwaysFalse).Results()
		c.So(res, shouldSlicesResemble, arr0)
	})

	c.Convey("Skip all", t, func() {
		res, _ := From(arr0).SkipWhile(alwaysTrue).Results()
		c.So(res, shouldSlicesResemble, empty)
	})

	c.Convey("Skip only first", t, func() {
		in := []int{1, 2, 3, 4, 5}
		res, _ := From(in).SkipWhile(func(i T) (bool, error) { return i.(int) < 2, nil }).Results()
		c.So(res, shouldSlicesResemble, in[1:])
	})

	c.Convey("SkipWhile & TakeWhile & SkipWhile", t, func() {
		in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		lessThanTwo := func(i T) (bool, error) { return i.(int) < 2, nil }
		lessThanSix := func(i T) (bool, error) { return i.(int) < 6, nil }
		res, _ := From(in).SkipWhile(alwaysFalse).SkipWhile(lessThanTwo).TakeWhile(lessThanSix).Results()
		c.So(res, shouldSlicesResemble, []int{2, 3, 4, 5})
	})
}

func TestOrder(t *testing.T) {
	c.Convey("Sort ints", t, func() {
		arr := []int{6, 1, 4, 0, -1, 2}
		arrSorted := []int{-1, 0, 1, 2, 4, 6}
		unsupportedArr := []T{6, 1, 4, 0, -1, 2, ""}

		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr).Where(erroneusBinaryFunc).OrderInts().Results()
			c.So(err, c.ShouldNotEqual, nil)
		})

		c.Convey("Sort order is correct", func() {
			res, _ := From(arr).OrderInts().Results()
			c.So(res, shouldSlicesResemble, arrSorted)
		})

		c.Convey("Sequence contain unsupported types", func() {
			_, err := From(unsupportedArr).OrderInts().Results()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
		})
	})

	c.Convey("Sort float64s", t, func() {
		arr := []float64{1.000000001, 1.0000000001, 0.1, 0.01, 0.00001, 0.0000000000001}
		arrSorted := []float64{0.0000000000001, 0.00001, 0.01, 0.1, 1.0000000001, 1.000000001}
		unsupportedArr := []T{1.000000001, "", 1.0000000001, 0.1, nil}

		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr).Where(erroneusBinaryFunc).OrderFloat64s().Results()
			c.So(err, c.ShouldNotEqual, nil)
		})

		c.Convey("Sort order is correct", func() {
			res, _ := From(arr).OrderFloat64s().Results()
			c.So(res, shouldSlicesResemble, arrSorted)
		})

		c.Convey("Sequence contain unsupported types", func() {
			_, err := From(unsupportedArr).OrderFloat64s().Results()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
		})
	})

	c.Convey("Sort strings", t, func() {
		arr := []string{"c", "a", "", "aa", "b"}
		arrSorted := []string{"", "a", "aa", "b", "c"}

		unsupportedArr := []T{"", "aa", "ccc", nil}

		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr).Where(erroneusBinaryFunc).OrderStrings().Results()
			c.So(err, c.ShouldNotEqual, nil)
		})

		c.Convey("Sort order is correct", func() {
			res, _ := From(arr).OrderStrings().Results()
			c.So(res, shouldSlicesResemble, arrSorted)
		})

		c.Convey("Sequence contain unsupported types", func() {
			_, err := From(unsupportedArr).OrderStrings().Results()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
		})
	})

}

func TestOrderBy(t *testing.T) {
	unsorted := []*foo{&foo{"A", 5}, &foo{"B", 1}, &foo{"C", 3}}
	sorted := []*foo{&foo{"B", 1}, &foo{"C", 3}, &foo{"A", 5}}
	sortByNum := func(this T, that T) bool {
		_this := this.(*foo)
		_that := that.(*foo)
		return _this.num <= _that.num
	}

	c.Convey("Nil comparator passed", t, func() {
		_, err := From(unsorted).OrderBy(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
	})
	c.Convey("Previous error is reflected in result", t, func() {
		_, err := From(unsorted).Where(erroneusBinaryFunc).OrderBy(sortByNum).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Sort empty", t, func() {
		res, _ := From(empty).OrderBy(sortByNum).Results()
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("Sort on structs", t, func() {
		res, _ := From(unsorted).OrderBy(sortByNum).Results()
		c.So(res, shouldSlicesResemble, sorted)
	})
}

func TestJoins(t *testing.T) {
	type Person struct{ Name string }
	type Pet struct {
		Name  string
		Owner Person
	}
	type ResultPair struct{ OwnerName, PetName string }
	type ResultGroup struct {
		OwnerName string
		Pets      []T
	}
	magnus := Person{Name: "Hedlund, Magnus"}
	terry := Person{Name: "Adams, Terry"}
	charlotte := Person{Name: "Weiss, Charlotte"}
	ahmet := Person{Name: "Balkan, Ahmet"}
	bob := Person{Name: "Marley, Bob"}

	barley := Pet{Name: "Barley", Owner: terry}
	boots := Pet{Name: "Boots", Owner: terry}
	whiskers := Pet{Name: "Whiskers", Owner: charlotte}
	daisy := Pet{Name: "Daisy", Owner: magnus}
	sasha := Pet{Name: "Sasha", Owner: bob}

	people := []Person{magnus, terry, charlotte, ahmet}
	pets := []Pet{barley, boots, whiskers, daisy, sasha}

	var dummyKeySelector = func(i T) T { return i }
	var dummyResultSelector = func(i, j T) T { return nil }
	var dummyGroupResultSelector = func(outer T, inner []T) T { return nil }

	equiJoinExpected := []T{
		ResultPair{magnus.Name, daisy.Name},
		ResultPair{terry.Name, barley.Name},
		ResultPair{terry.Name, boots.Name},
		ResultPair{charlotte.Name, whiskers.Name}}

	groupJoinExpected := []T{
		ResultGroup{magnus.Name, []T{daisy}},
		ResultGroup{terry.Name, []T{barley, boots}},
		ResultGroup{charlotte.Name, []T{whiskers}},
		ResultGroup{ahmet.Name, []T{}}}

	c.Convey("Equi-join", t, func() {
		c.Convey("Errors from the previous of the chain are carried on", func() {
			_, err := From(people).Where(erroneusBinaryFunc).Join(pets, dummyKeySelector,
				dummyKeySelector, dummyResultSelector).Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Passed non-slice value, error returned", func() {
			_, err := From(empty).Join("someString", nil, nil, nil).Results()
			c.So(err, c.ShouldEqual, ErrInvalidInput)
		})
		c.Convey("Nil funcs passed", func() {
			_, err := From(people).Join(pets, nil, nil, nil).Results()
			c.So(err, c.ShouldEqual, ErrNilFunc)
		})
		c.Convey("Nil input passed", func() {
			_, err := From(people).Join(nil, dummyKeySelector, dummyKeySelector,
				dummyResultSelector).Results()
			c.So(err, c.ShouldEqual, ErrNilInput)
		})
		c.Convey("Pets & owners example join (also checks preserving the order)", func() {

			res, err := From(people).Join(pets,
				func(person T) T { return person.(Person).Name },
				func(pet T) T { return pet.(Pet).Owner.Name },
				func(outer T, inner T) T {
					return ResultPair{outer.(Person).Name, inner.(Pet).Name}
				}).Results()
			c.So(err, c.ShouldEqual, nil)
			c.So(res, shouldSlicesResemble, equiJoinExpected)
		})
	})

	c.Convey("Group-join", t, func() {
		c.Convey("Errors from the previous of the chain are carried on", func() {
			_, err := From(people).Where(erroneusBinaryFunc).GroupJoin(pets, dummyKeySelector,
				dummyKeySelector, dummyGroupResultSelector).Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Passed non-slice value, error returned", func() {
			_, err := From(empty).GroupJoin("someString", nil, nil, nil).Results()
			c.So(err, c.ShouldEqual, ErrInvalidInput)
		})
		c.Convey("Nil funcs passed", func() {
			_, err := From(people).GroupJoin(pets, nil, nil, nil).Results()
			c.So(err, c.ShouldEqual, ErrNilFunc)
		})
		c.Convey("Nil input passed", func() {
			_, err := From(people).GroupJoin(nil, dummyKeySelector, dummyKeySelector,
				dummyGroupResultSelector).Results()
			c.So(err, c.ShouldEqual, ErrNilInput)
		})
		c.Convey("Pets & owners example join (also checks preserving the order)", func() {

			res, err := From(people).GroupJoin(pets,
				func(person T) T { return person.(Person).Name },
				func(pet T) T { return pet.(Pet).Owner.Name },
				func(outer T, inners []T) T {
					return ResultGroup{outer.(Person).Name, inners}
				}).Results()
			c.So(err, c.ShouldEqual, nil)
			c.So(res, shouldSlicesResemble, groupJoinExpected)
		})
	})
}

func TestGroupBy(t *testing.T) {
	c.Convey("Group empty slice into empty map", t, func() {
		res, err := From(empty).GroupBy(func(foo T) T { return foo }, func(foo T) T { return foo })
		c.So(err, c.ShouldEqual, nil)
		c.So(res, c.ShouldResemble, make(map[T][]T))
	})

	type Pet struct {
		Name  string
		Owner string
	}

	barley := Pet{Name: "Barley", Owner: "Damon Zhao"}
	boots := Pet{Name: "Boots", Owner: "Damon Zhao"}
	whiskers := Pet{Name: "Whiskers", Owner: "A-limon"}
	daisy := Pet{Name: "Daisy", Owner: "A-limon"}
	sasha := Pet{Name: "Sasha", Owner: "Bob"}

	pets := []Pet{barley, boots, whiskers, daisy, sasha}

	groupByExpected := map[T][]T{
		"Damon Zhao": []T{barley.Name, boots.Name},
		"A-limon":    []T{whiskers.Name, daisy.Name},
		"Bob":        []T{sasha.Name},
	}

	c.Convey("Pets group by owner", t, func() {
		res, err := From(pets).GroupBy(func(pet T) T { return pet.(Pet).Owner }, func(pet T) T { return pet.(Pet).Name })
		c.So(err, c.ShouldEqual, nil)
		c.So(res, c.ShouldResemble, groupByExpected)
	})
}

func TestRange(t *testing.T) {
	c.Convey("count < 0", t, func() {
		_, err := Range(1, -1).Results()
		c.So(err, c.ShouldEqual, ErrNegativeParam)
	})
	c.Convey("count = 0", t, func() {
		res, err := Range(1, 0).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, empty)
	})
	c.Convey("range(1,10)", t, func() {
		res, err := Range(1, 10).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, shouldSlicesResemble, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})
}

var (
	intArr            = []int{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10}
	intArrSumExpected = -55
	intArrAvgExpected = float64(intArrSumExpected) / float64(len(intArr))
	mixedArr          = []T{
		0, int(1), int8(2), int16(3), int32(4), int64(5),
		uint(6), uint8(7), uint16(8), uint32(9), uint64(10),
		float32(11.11), float64(12.12)}
	mixedArrSumExpected           = float64(78.23)
	mixedArrAvgExpected           = float64(mixedArrSumExpected) / float64(len(mixedArr))
	mixedArrContainingNil         = []T{1, 2, nil, float64(3), uint(4)}
	mixedArrContainingUnsupported = []T{1, 2, foo{"", 0}}
)

func TestSum(t *testing.T) {
	c.Convey("Previous errors are reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Sum()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Empty slice", t, func() {
		res, err := From(empty).Sum()
		c.So(err, c.ShouldEqual, nil)
		c.So(res, c.ShouldEqual, 0.0)
	})
	c.Convey("Slice of ints", t, func() {
		res, _ := From(intArr).Sum()
		c.So(res, c.ShouldEqual, intArrSumExpected)
	})
	c.Convey("Slice of mixed numeric types", t, func() {
		res, _ := From(mixedArr).Sum()
		c.So(res, c.ShouldAlmostEqual, mixedArrSumExpected, 0.000001) // float32 requires less tolerance than goconvey default
	})
	c.Convey("Slice with numeric types and nils", t, func() {
		_, err := From(mixedArrContainingNil).Sum()
		c.So(err, c.ShouldEqual, ErrNan)
	})
	c.Convey("Slice contains unsupported type", t, func() {
		_, err := From(mixedArrContainingUnsupported).Sum()
		c.So(err, c.ShouldNotEqual, nil)
	})
}

func TestAverage(t *testing.T) {
	c.Convey("Previous errors are reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Average()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Empty slice", t, func() {
		_, err := From(empty).Average()
		c.So(err, c.ShouldEqual, ErrEmptySequence)
	})
	c.Convey("Slice of ints", t, func() {
		res, _ := From(intArr).Average()
		c.So(res, c.ShouldEqual, intArrAvgExpected)
	})
	c.Convey("Slice of mixed numeric types", t, func() {
		res, _ := From(mixedArr).Average()
		c.So(res, c.ShouldAlmostEqual, mixedArrAvgExpected, 0.000001) // float32 requires less tolerance than goconvey default
	})
	c.Convey("Slice with numeric types and nils", t, func() {
		_, err := From(mixedArrContainingNil).Average()
		c.So(err, c.ShouldEqual, ErrNan)
	})
	c.Convey("Slice contains unsupported type", t, func() {
		_, err := From(mixedArrContainingUnsupported).Average()
		c.So(err, c.ShouldNotEqual, nil)
	})
}

func TestMinMax(t *testing.T) {
	c.Convey("MinInt/MaxInt", t, func() {
		var (
			arr            = []int{-1, -9, 0, 9, 1}
			arrUnsupported = []T{-1, -9, 0, 9, 1, nil}
			expectedMin    = -9
			expectedMax    = 9
		)
		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).MinInt()
			c.So(err, c.ShouldNotEqual, nil)
			_, err = From(arr0).Where(erroneusBinaryFunc).MaxInt()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Empty slice", func() {
			_, err := From(empty).MinInt()
			c.So(err, c.ShouldEqual, ErrEmptySequence)
			_, err = From(empty).MaxInt()
			c.So(err, c.ShouldEqual, ErrEmptySequence)
		})
		c.Convey("Sequence contains unsupported types", func() {
			_, err := From(arrUnsupported).MinInt()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
			_, err = From(arrUnsupported).MaxInt()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
		})
		c.Convey("Verify min/max result", func() {
			res, _ := From(arr).MinInt()
			c.So(res, c.ShouldEqual, expectedMin)
			res, _ = From(arr).MaxInt()
			c.So(res, c.ShouldEqual, expectedMax)
		})
	})
	c.Convey("MinUint/MaxUint", t, func() {
		var (
			arr            = []uint{uint(1), uint(9), uint(100), uint(99), uint(0)}
			arrUnsupported = []T{uint(1), uint(9), uint(100), uint(99), uint(0), 0}
			expectedMin    = uint(0)
			expectedMax    = uint(100)
		)
		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).MinUint()
			c.So(err, c.ShouldNotEqual, nil)
			_, err = From(arr0).Where(erroneusBinaryFunc).MaxUint()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Empty slice", func() {
			_, err := From(empty).MinUint()
			c.So(err, c.ShouldEqual, ErrEmptySequence)
			_, err = From(empty).MaxUint()
			c.So(err, c.ShouldEqual, ErrEmptySequence)
		})
		c.Convey("Sequence contains unsupported types", func() {
			_, err := From(arrUnsupported).MinUint()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
			_, err = From(arrUnsupported).MaxUint()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
		})
		c.Convey("Verify min/max result", func() {
			res, _ := From(arr).MinUint()
			c.So(res, c.ShouldEqual, expectedMin)
			res, _ = From(arr).MaxUint()
			c.So(res, c.ShouldEqual, expectedMax)
		})
	})
	c.Convey("MinFloat64/MaxFloat64", t, func() {
		var (
			arr            = []float64{float64(-9), float64(-9.9), float64(0), float64(99), float64(99.9)}
			arrUnsupported = []T{float64(-9), float64(-9.9), float64(0), float64(99), float64(99.9), uint(0)}
			expectedMin    = float64(-9.9)
			expectedMax    = float64(99.9)
		)
		c.Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).MinFloat64()
			c.So(err, c.ShouldNotEqual, nil)
			_, err = From(arr0).Where(erroneusBinaryFunc).MaxFloat64()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Empty slice", func() {
			_, err := From(empty).MinFloat64()
			c.So(err, c.ShouldEqual, ErrEmptySequence)
			_, err = From(empty).MaxFloat64()
			c.So(err, c.ShouldEqual, ErrEmptySequence)
		})
		c.Convey("Sequence contains unsupported types", func() {
			_, err := From(arrUnsupported).MinFloat64()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
			_, err = From(arrUnsupported).MaxFloat64()
			c.So(err, c.ShouldEqual, ErrTypeMismatch)
		})
		c.Convey("Verify min/max result", func() {
			res, _ := From(arr).MinFloat64()
			c.So(res, c.ShouldEqual, expectedMin)
			res, _ = From(arr).MaxFloat64()
			c.So(res, c.ShouldEqual, expectedMax)
		})
	})
}
