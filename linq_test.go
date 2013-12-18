package linq

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

type foo struct {
	str string
	num int
}

var genericError = errors.New("")

var (
	empty = []interface{}{}
	arr0  = []interface{}{1, 2, 3, 1, 2}
	arr1  = []interface{}{"foo", "bar", "baz"}
	arr2  = []interface{}{nil, "foo", 3.14, true, false}
	arr3  = []interface{}{foo{"A", 0}, foo{"B", 1}, foo{"C", -1}}
	arr4  = []interface{}{&foo{"C", 0xffff}, nil, &foo{"D", 0x7fff}, byte(12), nil}
)

func TestFrom(t *testing.T) {
	Convey("When passed nil value, error returned", t, func() {
		So(From(nil).err, ShouldNotEqual, nil)
	})

	Convey("When passed non-nil value, structure should have the exact same slice at different location", t, func() {
		Convey("Empty array", func() {
			r := From(empty).values
			So(r, ShouldResemble, empty)
			So(r, ShouldNotEqual, empty) // slice copied?
		})
		Convey("Non-empty arrays", func() {

			Convey("Passed & held slices are different", func() {
				So(From(arr0).values, ShouldNotEqual, arr0)
				So(From(arr4).values, ShouldNotEqual, arr4)
			})

			Convey("Deep slice equality", func() {
				So(From(arr0).values, ShouldResemble, arr0)
				So(From(arr1).values, ShouldResemble, arr1)
				So(From(arr2).values, ShouldResemble, arr2)
				So(From(arr3).values, ShouldResemble, arr3)
				So(From(arr4).values, ShouldResemble, arr4)
			})
		})
	})
}

func TestResults(t *testing.T) {
	Convey("If error exists in given queryable, error is returned", t, func() {
		errMsg := "dummy error"
		q := queryable{
			values: nil,
			err:    errors.New(errMsg)}
		_, err := q.Results()
		So(err, ShouldNotEqual, nil)
		So(err.Error(), ShouldEqual, errMsg)
	})
	Convey("Given no errors exist, non-nil results are returned", t, func() {
		q := queryable{values: arr0, err: nil}
		val, err := q.Results()
		So(err, ShouldEqual, nil)
		So(val, ShouldResemble, arr0)
	})
}

var alwaysTrue = func(i interface{}) (bool, error) {
	return true, nil
}
var alwaysFalse = func(i interface{}) (bool, error) {
	return false, nil
}
var erroneusBinaryFunc = func(i interface{}) (bool, error) {
	return true, genericError
}

func TestWhere(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})

	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Where(alwaysTrue).Results()
		So(err, ShouldNotEqual, nil)

		Convey("Chain successful and erroneus functions", func() {
			_, err := From(arr0).Where(alwaysTrue).Where(erroneusBinaryFunc).Results()
			So(err, ShouldNotEqual, nil)

			Convey("Erroneus function is in chain but not called", func() {
				_, err := From(arr0).Where(alwaysTrue).Where(alwaysFalse).Where(erroneusBinaryFunc).Results()
				So(err, ShouldEqual, nil)
			})
		})
	})

	Convey("Chose none of the elements", t, func() {
		val, _ := From(arr0).Where(alwaysFalse).Results()
		So(val, ShouldEqual, nil)
	})

	Convey("Chose all elements, as is", t, func() {
		val, _ := From(arr0).Where(alwaysTrue).Results()
		So(val, ShouldResemble, arr0)
	})

	Convey("Basic filtering (x mod 2)==0", t, func() {
		n := 100
		divisibleBy2 := func(i interface{}) (bool, error) {
			return i.(int)%2 == 0, nil
		}
		arr := make([]interface{}, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		val, err := From(arr).Where(divisibleBy2).Results()
		So(err, ShouldEqual, nil)
		So(len(val), ShouldEqual, n/2)
	})
}

func TestSelect(t *testing.T) {
	asIs := func(i interface{}) (interface{}, error) {
		return i, nil
	}
	erroneusFunc := func(i interface{}) (interface{}, error) {
		return nil, genericError
	}

	Convey("Nil func returns error", t, func() {
		_, err := From(arr0).Select(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})

	Convey("Error returned from provided func", t, func() {
		val, err := From(arr0).Select(erroneusFunc).Results()
		So(err, ShouldNotEqual, nil)

		Convey("Erroneus function is in chain with as-is select", func() {
			_, err = From(arr0).Select(asIs).Select(erroneusFunc).Results()
			So(err, ShouldNotEqual, nil)
		})

		Convey("Erroneus function is in chain but not called", func() {
			val, err = From(arr0).Where(alwaysFalse).Select(erroneusFunc).Results()
			So(err, ShouldEqual, nil)
			So(len(val), ShouldEqual, 0)
		})

	})

	Convey("Select all elements as is", t, func() {
		val, err := From(arr0).Select(asIs).Results()
		So(err, ShouldEqual, nil)
		So(val, ShouldResemble, arr0)
	})

	Convey("Pow(x,2) for i in []int", t, func() {
		pow := func(i interface{}) (interface{}, error) {
			return i.(int) * i.(int), nil
		}
		val, err := From(arr0).Select(pow).Results()
		So(err, ShouldEqual, nil)
		arr := make([]int, len(arr0))
		for j, i := range arr0 {
			arr[j] = i.(int) * i.(int)
		}
		res := make([]int, len(val))
		for j, v := range val {
			res[j] = v.(int)
		}
		So(res, ShouldResemble, arr)
	})
}

func TestDistinct(t *testing.T) {
	Convey("Empty slice", t, func() {
		res, err := From(empty).Distinct().Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})

	allSameInt := []interface{}{1, 1, 1, 1, 1, 1, 1, 1, 1}
	allSameStruct := []interface{}{foo{"A", -1}, foo{"A", -1}, foo{"A", -1}}
	allNil := []interface{}{nil, nil, nil, nil, nil, nil, nil, nil, nil}

	Convey("With default equality comparer ==", t, func() {

		Convey("All elements are the same", func() {
			res, _ := From(allSameInt).Distinct().Results()
			So(res, ShouldResemble, []interface{}{allSameInt[0]})

			Convey("All elements are nil", func() {
				res, _ = From(allNil).Distinct().Results()
				So(res, ShouldResemble, []interface{}{allNil[0]})
			})
		})

		Convey("Distinct on structs and nils", func() {
			arr := []interface{}{foo{"A", 0xffff}, nil, foo{"B", 0x7fff}, nil, foo{"A", 0xffff}}
			res, _ := From(arr).Distinct().Results()
			So(len(res), ShouldEqual, 3)
		})

		Convey("Randomly generated integers with duplicates or more", func() {
			var arr = make([]interface{}, 10000)
			var dict = make(map[int]bool, len(arr))

			rand.Seed(time.Now().UnixNano())
			for i := 0; i < len(arr); i++ {
				r := rand.Intn(len(arr) * 4 / 5) // collision 20%
				arr[i] = r
				dict[r] = true
			}
			res, _ := From(arr).Distinct().Results()
			So(len(res), ShouldEqual, len(dict))
		})
	})

	Convey("With provided equality comparer", t, func() {
		fooComparer := func(i interface{}, j interface{}) (bool, error) {
			return i.(foo).str == j.(foo).str && i.(foo).num == j.(foo).num, nil
		}
		fooPtrComparer := func(i interface{}, j interface{}) (bool, error) {
			return i.(*foo).str == j.(*foo).str && i.(*foo).num == j.(*foo).num, nil
		}
		_ = fooPtrComparer

		erroneusComparer := func(i interface{}, j interface{}) (bool, error) {
			return false, genericError
		}

		Convey("Comparer returns error", func() {
			_, err := From(arr0).DistinctBy(erroneusComparer).Results()
			So(err, ShouldNotEqual, nil)
		})

		Convey("All elements are the same", func() {
			res, _ := From(allSameStruct).DistinctBy(fooComparer).Results()
			So(res, ShouldResemble, []interface{}{allSameStruct[0]})
		})
		Convey("All elements are distinct", func() {
			var arr = make([]interface{}, 100)
			for i := 0; i < len(arr); i++ {
				arr[i] = i
			}
			res, _ := From(arr).DistinctBy(func(this interface{}, that interface{}) (bool, error) {
				return this.(int) == that.(int), nil
			}).Results()
			So(res, ShouldResemble, arr)
		})

		Convey("Ensure leftmost appearance is returned in multiple occurrence cases", func() {
			arr := []interface{}{&foo{"A", 0}, &foo{"B", 0}, &foo{"A", 0}, &foo{"C", 0}, &foo{"A", 0}, &foo{"B", 0}}
			res, _ := From(arr).DistinctBy(fooPtrComparer).Results()
			So(len(res), ShouldResemble, 3)
			So(res[0], ShouldEqual, arr[0]) // A
			So(res[1], ShouldEqual, arr[1]) // B
			So(res[2], ShouldEqual, arr[3]) // C
		})

		Convey("Randomly generated integers with likely collisions", func() {
			var arr = make([]interface{}, 10000)
			var dict = make(map[int]bool, len(arr))
			rand.Seed(time.Now().UnixNano())
			for i := 0; i < len(arr); i++ {
				r := rand.Intn(len(arr) * 4 / 5) // collision 20%
				arr[i] = r
				dict[r] = true
			}
			res, _ := From(arr).DistinctBy(func(this interface{}, that interface{}) (bool, error) {
				return this.(int) == that.(int), nil
			}).Results()
			So(len(res), ShouldEqual, len(dict))
		})
	})
}

func TestUnion(t *testing.T) {
	uniqueArr0 := []interface{}{1, 2, 3, 4, 5}
	uniqueArr1 := []interface{}{"a", "b", "c"}
	allSameArr := []interface{}{1, 1, 1, 1}
	sameStruct0 := []interface{}{foo{"A", 0}, foo{"B", 0}}
	sameStruct1 := []interface{}{foo{"B", 0}, foo{"A", 0}}
	Convey("Empty ∪ nil", t, func() {
		_, err := From(empty).Union(nil).Results()
		So(err, ShouldEqual, ErrNilInput)
	})
	Convey("Empty ∪ empty", t, func() {
		res, _ := From(empty).Union(empty).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("Empty ∪ non-empty", t, func() {
		res, _ := From(empty).Union(uniqueArr0).Results()
		So(res, ShouldResemble, uniqueArr0)
	})
	Convey("Non-empty ∪ empty", t, func() {
		res, _ := From(uniqueArr0).Union(empty).Results()
		So(res, ShouldResemble, uniqueArr0)
	})
	Convey("(Unique slice) ∪ (itself)", t, func() {
		res, _ := From(uniqueArr0).Union(uniqueArr0).Results()
		So(res, ShouldResemble, uniqueArr0)
	})
	Convey("(All same slice) ∪ (itself)", t, func() {
		res, _ := From(allSameArr).Union(allSameArr).Results()
		So(len(res), ShouldEqual, 1)
	})
	Convey("Mixed types", t, func() {
		res, _ := From(uniqueArr0).Union(uniqueArr1).Results()
		So(len(res), ShouldEqual, len(uniqueArr0)+len(uniqueArr1))
	})
	Convey("Same-type structs", t, func() {
		res, _ := From(sameStruct0).Union(sameStruct1).Results()
		So(len(res), ShouldEqual, len(sameStruct1))
	})
}

func TestExcept(t *testing.T) {
	uniqueArr := []interface{}{1, 2, 3, 4, 5}
	allSameArr := []interface{}{1, 1, 1, 1}
	Convey("Empty ∖ nil", t, func() {
		_, err := From(empty).Except(nil).Results()
		So(err, ShouldEqual, ErrNilInput)
	})
	Convey("Empty ∖ empty", t, func() {
		res, _ := From(empty).Except(empty).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("Empty ∖ non-empty", t, func() {
		res, _ := From(empty).Except(uniqueArr).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("Non-empty ∖ empty", t, func() {
		res, _ := From(uniqueArr).Except(empty).Results()
		So(res, ShouldResemble, uniqueArr)
	})
	Convey("(Unique set) ∖ (itself)", t, func() {
		res, _ := From(uniqueArr).Except(uniqueArr).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("(All same slice) ∖ (itself)", t, func() {
		res, _ := From(allSameArr).Except(allSameArr).Results()
		So(len(res), ShouldEqual, 0)
	})
	Convey("There is some intersection", t, func() {
		res, _ := From([]interface{}{1, 2, 3, 4, 5}).Except([]interface{}{3, 4, 5, 6, 7}).Results()
		So(res, ShouldResemble, []interface{}{1, 2})
	})
}

func TestCount(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).CountBy(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).CountBy(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).CountBy(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("No matches", t, func() {
		c, _ := From(arr0).CountBy(alwaysFalse)
		So(c, ShouldEqual, 0)
		c, _ = From(arr0).Where(alwaysFalse).Count()
		So(c, ShouldEqual, 0)
	})
	Convey("All matches", t, func() {
		c, _ := From(arr0).CountBy(alwaysTrue)
		So(c, ShouldEqual, len(arr0))
		c, _ = From(arr0).Count()
		So(c, ShouldEqual, len(arr0))
	})
}

func TestAny(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AnyWith(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AnyWith(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).AnyWith(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("No matches", t, func() {
		r, _ := From(arr0).AnyWith(alwaysFalse)
		So(r, ShouldEqual, false)
		r, _ = From(arr0).Where(alwaysFalse).Any()
		So(r, ShouldEqual, false)
	})
	Convey("All matches", t, func() {
		r, _ := From(arr0).AnyWith(alwaysTrue)
		So(r, ShouldEqual, true)
		r, _ = From(arr0).Where(alwaysTrue).Any()
		So(r, ShouldEqual, true)
	})
}

func TestSingle(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Single(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Single(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).Single(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("No matches", t, func() {
		r, _ := From(arr0).Single(alwaysFalse)
		So(r, ShouldEqual, false)
	})
	Convey("All matches", t, func() {
		r, _ := From(arr0).Single(alwaysTrue)
		So(r, ShouldEqual, false)
	})
	Convey("Only one match", t, func() {
		var match0 = func(i interface{}) (bool, error) {
			return i.(int) == 0, nil
		}
		r, _ := From([]interface{}{-1, -1, 0, 1, 1}).Single(match0)
		So(r, ShouldEqual, true)
		r, _ = From([]interface{}{0, 1, 2, 2, 0}).Single(match0)
		So(r, ShouldEqual, false)
	})
}

func TestAll(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).All(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).All(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).All(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("Empty slice", t, func() {
		r, _ := From(empty).All(alwaysTrue)
		So(r, ShouldEqual, true)
	})
	Convey("No matches", t, func() {
		r, _ := From(arr0).All(alwaysFalse)
		So(r, ShouldEqual, false)
	})
	Convey("All matches", t, func() {
		r, _ := From(arr0).All(alwaysTrue)
		So(r, ShouldEqual, true)
	})
	Convey("Multiple matches", t, func() {
		match0 := func(i interface{}) (bool, error) {
			return i.(int) == 0, nil
		}
		r, _ := From([]interface{}{0, 1, 2, 2, 0}).All(match0)
		So(r, ShouldEqual, false)
	})
}

func TestFirst_FirstOrNil(t *testing.T) {
	Convey("empty.First is ErrNoElement", t, func() {
		_, err := From(empty).First()
		So(err, ShouldEqual, ErrNoElement)
	})
	Convey("empty.FirstOrNil is nil", t, func() {
		v, _ := From(empty).FirstOrNil()
		So(v, ShouldEqual, nil)
	})
	Convey("first element is returned", t, func() {
		v, _ := From(arr3).First()
		So(v, ShouldResemble, arr3[0])
	})
	Convey("previous errors are reflected", t, func() {
		_, err1 := From(arr0).Where(erroneusBinaryFunc).First()
		_, err2 := From(arr0).Where(erroneusBinaryFunc).FirstOrNil()
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})
}

func TestFirstBy_FirstOrNilBy(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err1 := From(arr0).FirstBy(nil)
		_, err2 := From(arr0).FirstBy(nil)
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})

	Convey("erroneus function reflected on result", t, func() {
		_, err1 := From(arr0).FirstBy(erroneusBinaryFunc)
		_, err2 := From(arr0).FirstBy(erroneusBinaryFunc)
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})
	Convey("empty.FirstBy is ErrNoElement", t, func() {
		_, err1 := From(empty).FirstBy(alwaysFalse)
		So(err1, ShouldEqual, ErrNoElement)
	})
	Convey("empty.FirstOrNilBy is ErrNoElement", t, func() {
		_, err1 := From(empty).FirstBy(alwaysFalse)
		So(err1, ShouldEqual, ErrNoElement)
	})
	Convey("Actual first element is returned", t, func() {
		val, _ := From(arr3).FirstBy(alwaysTrue)
		So(val, ShouldResemble, arr3[0])
		val, _ = From(arr3).FirstOrNilBy(alwaysTrue)
		So(val, ShouldResemble, arr3[0])
	})
	Convey("No matches", t, func() {
		_, err := From(arr3).FirstBy(alwaysFalse)
		So(err, ShouldEqual, ErrNoElement)
		elm, err := From(arr3).FirstOrNilBy(alwaysFalse)
		So(err, ShouldEqual, nil)
		So(elm, ShouldEqual, nil)
	})
}

func TestLast_LastOrNil(t *testing.T) {
	Convey("empty.Last is ErrNoElement", t, func() {
		_, err := From(empty).Last()
		So(err, ShouldEqual, ErrNoElement)
	})
	Convey("empty.LastOrNil is nil", t, func() {
		v, _ := From(empty).LastOrNil()
		So(v, ShouldEqual, nil)
	})
	Convey("Last element is returned", t, func() {
		v, _ := From(arr3).Last()
		So(v, ShouldResemble, arr3[len(arr3)-1])
	})
	Convey("previous errors are reflected", t, func() {
		_, err1 := From(arr0).Where(erroneusBinaryFunc).Last()
		_, err2 := From(arr0).Where(erroneusBinaryFunc).LastOrNil()
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})
}

func TestLastBy_LastOrNilBy(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err1 := From(arr0).LastBy(nil)
		_, err2 := From(arr0).LastBy(nil)
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})

	Convey("erroneus function reflected on result", t, func() {
		_, err1 := From(arr0).LastBy(erroneusBinaryFunc)
		_, err2 := From(arr0).LastBy(erroneusBinaryFunc)
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})
	Convey("empty.LastBy is ErrNoElement", t, func() {
		_, err1 := From(empty).LastBy(alwaysFalse)
		So(err1, ShouldEqual, ErrNoElement)
	})

	Convey("empty.LastOrNilBy is ErrNoElement", t, func() {
		_, err1 := From(empty).LastBy(alwaysFalse)
		So(err1, ShouldEqual, ErrNoElement)
	})

	Convey("Actual last element is returned", t, func() {
		val1, _ := From(arr3).LastBy(alwaysTrue)
		val2, _ := From(arr3).LastOrNilBy(alwaysTrue)
		So(val1, ShouldResemble, arr3[len(arr3)-1])
		So(val2, ShouldResemble, arr3[len(arr3)-1])
	})
	Convey("No matches", t, func() {
		_, err := From(arr3).LastBy(alwaysFalse)
		So(err, ShouldEqual, ErrNoElement)
		elm, err := From(arr3).LastOrNilBy(alwaysFalse)
		So(err, ShouldEqual, nil)
		So(elm, ShouldEqual, nil)
	})

}

func TestReverse(t *testing.T) {
	Convey("Previous errors are returned", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Reverse().Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Reversing empty", t, func() {
		res, err := From(empty).Reverse().Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})
	Convey("Actual reverse", t, func() {
		arr := []interface{}{1, 2, 3, 4, 5}
		rev := []interface{}{5, 4, 3, 2, 1}
		res, _ := From(arr).Reverse().Results()
		So(res, ShouldResemble, rev)

		Convey("Slice containing nils", func() {
			arr := []interface{}{1, nil, nil, 2, nil, 3, nil}
			rev := []interface{}{nil, 3, nil, 2, nil, nil, 1}
			res, _ := From(arr).Reverse().Results()
			So(res, ShouldResemble, rev)
		})
	})
}

func TestTake(t *testing.T) {
	Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Take(1).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty slice take n>0", t, func() {
		res, err := From(empty).Take(1).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})

	Convey("Take 0", t, func() {
		res, _ := From(arr0).Take(0).Results()
		So(res, ShouldResemble, empty)
	})

	Convey("Take n < 0", t, func() {
		res, err := From(arr0).Take(-1).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})

	Convey("Take n > 0", t, func() {
		in := []interface{}{1, 2, 3, 4, 5}
		res, _ := From(in).Take(3).Results()
		So(res, ShouldResemble, []interface{}{1, 2, 3})

		Convey("Take n ≥ len(arr)", func() {
			res, _ := From(in).Take(len(in)).Results()
			So(res, ShouldResemble, res)
			res, _ = From(in).Take(len(in) + 1).Results()
			So(res, ShouldResemble, res)
		})
	})
}

func TestSkip(t *testing.T) {
	Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Skip(1).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty slice Skip n>0", t, func() {
		res, err := From(empty).Skip(1).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})

	Convey("Skip 0", t, func() {
		res, _ := From(arr0).Skip(0).Results()
		So(res, ShouldResemble, arr0)
	})

	Convey("Skip n < 0", t, func() {
		res, err := From(arr0).Skip(-1).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, arr0)
	})

	Convey("Skip n > 0", t, func() {
		in := []interface{}{1, 2, 3, 4, 5}
		res, _ := From(in).Skip(3).Results()
		So(res, ShouldResemble, []interface{}{4, 5})

		Convey("Skip n ≥ len(arr)", func() {
			res, _ := From(in).Skip(len(in)).Results()
			So(res, ShouldResemble, empty)
			res, _ = From(in).Skip(len(in) + 1).Results()
			So(res, ShouldResemble, empty)
		})
	})

	Convey("Skip & take & skip", t, func() {
		in := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		res, _ := From(in).Skip(0).Skip(-1000).Skip(1).Take(1000).Take(5).Results()
		So(res, ShouldResemble, []interface{}{1, 2, 3, 4, 5})
	})
}

func TestOrder(t *testing.T) {
	Convey("Sort empty", t, func() {
		res, _ := From(empty).Order().Results()
		So(len(res), ShouldResemble, len(empty))
	})
	Convey("Sort ints", t, func() {
		in := []interface{}{6, 1, 4, 0, -1, 2}
		res, _ := From(in).Order().Results()
		So(res, ShouldResemble, []interface{}{-1, 0, 1, 2, 4, 6})
	})
	Convey("Sort float64s", t, func() {
		in := []interface{}{1.000000001, 1.0000000001, 0.1, 0.01, 0.00001, 0.0000000000001}
		res, _ := From(in).Order().Results()
		So(res, ShouldResemble, []interface{}{0.0000000000001, 0.00001, 0.01, 0.1, 1.0000000001, 1.000000001})
	})
	Convey("Sort strings", t, func() {
		in := []interface{}{"c", "a", "", "aa", "b"}
		res, _ := From(in).Order().Results()
		So(res, ShouldResemble, []interface{}{"", "a", "aa", "b", "c"})
	})
	Convey("Attempt with unsupported types", t, func() {
		in := []interface{}{true, false, true, nil, byte(10)}
		_, err := From(in).Order().Results()
		So(err, ShouldEqual, ErrUnsupportedType)
	})
}

func TestOrderBy(t *testing.T) {
	unsorted := []interface{}{&foo{"A", 5}, &foo{"B", 1}, &foo{"C", 3}}
	sorted := []interface{}{&foo{"B", 1}, &foo{"C", 3}, &foo{"A", 5}}
	sortByNum := func(this interface{}, that interface{}) bool {
		_this := this.(*foo)
		_that := that.(*foo)
		return _this.num <= _that.num
	}

	Convey("Nil comparator passed", t, func() {
		_, err := From(unsorted).OrderBy(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})
	Convey("Previous error is reflected in result", t, func() {
		_, err := From(unsorted).Where(erroneusBinaryFunc).OrderBy(sortByNum).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Sort empty", t, func() {
		res, _ := From(empty).OrderBy(sortByNum).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("Sort on structs", t, func() {
		res, _ := From(unsorted).OrderBy(sortByNum).Results()
		So(res, ShouldResemble, sorted)
	})
}

func TestJoin(t *testing.T) {
	type Person struct{ Name string }
	type Pet struct {
		Name  string
		Owner Person
	}
	type ResultPair struct{ OwnerName, PetName string }
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

	people := []interface{}{magnus, terry, charlotte, ahmet}
	pets := []interface{}{barley, boots, whiskers, daisy, sasha}

	var dummyKeySelector = func(i interface{}) interface{} { return i }
	var dummyValueSelector = func(i, j interface{}) interface{} { return nil }

	natJoinExpected := []interface{}{
		ResultPair{magnus.Name, daisy.Name},
		ResultPair{terry.Name, barley.Name},
		ResultPair{terry.Name, boots.Name},
		ResultPair{charlotte.Name, whiskers.Name}}

	Convey("Errors from the previous of the chain are carried on", t, func() {
		_, err := From(people).Where(erroneusBinaryFunc).Join(pets, dummyKeySelector, dummyKeySelector, dummyValueSelector).Results()
		So(err, ShouldNotEqual, nil)
	})

	Convey("Nil func passed", t, func() {
		_, err := From(people).Join(pets, nil, nil, nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})

	Convey("Nil input passed", t, func() {
		_, err := From(people).Join(nil, dummyKeySelector, dummyKeySelector, dummyValueSelector).Results()
		So(err, ShouldEqual, ErrNilInput)
	})

	Convey("Pets & owners example join (also checks preserving the order)", t, func() {

		res, err := From(people).Join(pets,
			func(person interface{}) interface{} { return person.(Person).Name },
			func(pet interface{}) interface{} { return pet.(Pet).Owner.Name },
			func(outer interface{}, inner interface{}) interface{} {
				return ResultPair{outer.(Person).Name, inner.(Pet).Name}
			}).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, natJoinExpected)
	})
}
