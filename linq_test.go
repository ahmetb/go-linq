package linq

import "testing"
import "errors"
import "math/rand"
import "time"
import . "github.com/smartystreets/goconvey/convey"

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
		q := Queryable{
			nil,
			errors.New(errMsg)}
		_, err := q.Results()
		So(err, ShouldNotEqual, nil)
		So(err.Error(), ShouldEqual, errMsg)
	})
	Convey("Given no errors exist, non-nil results are returned", t, func() {
		q := Queryable{arr0, nil}
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

		Convey("Distinct on []Struct", func() {
			arr := []interface{}{foo{"A", 0xffff}, foo{"B", 0x7fff}, foo{"A", 0xffff}}
			res, _ := From(arr).Distinct().Results()
			So(len(res), ShouldEqual, 2)
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
	})
	Convey("All matches", t, func() {
		c, _ := From(arr0).CountBy(alwaysTrue)
		So(c, ShouldEqual, len(arr0))
	})
}

func TestAny(t *testing.T) {
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Any(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Any(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).Any(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("No matches", t, func() {
		r, _ := From(arr0).Any(alwaysFalse)
		So(r, ShouldEqual, false)
	})
	Convey("All matches", t, func() {
		r, _ := From(arr0).Any(alwaysTrue)
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
		match0 := func(i interface{}) (bool, error) {
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
	Convey("empty.LastBy is ErrNoElement", t, t, func() {
		_, err1 := From(empty).LastBy(alwaysFalse)
		So(err1, ShouldEqual, ErrNoElement)
	})

	Convey("empty.LastOrNilBy is ErrNoElement", t, t, func() {
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
