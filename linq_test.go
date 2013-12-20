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

var (
	empty = []T{}
	arr0  = []T{1, 2, 3, 1, 2}
	arr1  = []T{"foo", "bar", "baz"}
	arr2  = []T{nil, "foo", 3.14, true, false}
	arr3  = []T{foo{"A", 0}, foo{"B", 1}, foo{"C", -1}}
	arr4  = []T{&foo{"C", 0xffff}, nil, &foo{"D", 0x7fff}, byte(12), nil}
)

var (
	genericError = errors.New("")
	alwaysTrue   = func(i T) (bool, error) {
		return true, nil
	}
	alwaysFalse = func(i T) (bool, error) {
		return false, nil
	}
	erroneusBinaryFunc = func(i T) (bool, error) {
		return true, genericError
	}
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
			values: nil,
			err:    errors.New(errMsg)}
		_, err := q.Results()
		So(err, ShouldNotEqual, nil)
		So(err.Error(), ShouldEqual, errMsg)
	})
	Convey("Given no errors exist, non-nil results are returned", t, func() {
		q := Queryable{values: arr0, err: nil}
		val, err := q.Results()
		So(err, ShouldEqual, nil)
		So(val, ShouldResemble, arr0)
	})
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
		divisibleBy2 := func(i T) (bool, error) {
			return i.(int)%2 == 0, nil
		}
		arr := make([]T, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		val, err := From(arr).Where(divisibleBy2).Results()
		So(err, ShouldEqual, nil)
		So(len(val), ShouldEqual, n/2)
	})
}

func TestSelect(t *testing.T) {
	asIs := func(i T) (T, error) {
		return i, nil
	}
	erroneusFunc := func(i T) (T, error) {
		return nil, genericError
	}

	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Select(asIs).Results()
		So(err, ShouldNotEqual, nil)
	})

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
		pow := func(i T) (T, error) {
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

	allSameInt := []T{1, 1, 1, 1, 1, 1, 1, 1, 1}
	allSameStruct := []T{foo{"A", -1}, foo{"A", -1}, foo{"A", -1}}
	allNil := []T{nil, nil, nil, nil, nil, nil, nil, nil, nil}

	Convey("With default equality comparer ==", t, func() {
		Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).Distinct().Results()
			So(err, ShouldNotEqual, nil)
		})
		Convey("All elements are the same", func() {
			res, _ := From(allSameInt).Distinct().Results()
			So(res, ShouldResemble, []T{allSameInt[0]})

			Convey("All elements are nil", func() {
				res, _ = From(allNil).Distinct().Results()
				So(res, ShouldResemble, []T{allNil[0]})
			})
		})
		Convey("Distinct on structs and nils", func() {
			arr := []T{foo{"A", 0xffff}, nil, foo{"B", 0x7fff}, nil, foo{"A", 0xffff}}
			res, _ := From(arr).Distinct().Results()
			So(len(res), ShouldEqual, 3)
		})
		Convey("Randomly generated integers with duplicates or more", func() {
			var arr = make([]T, 10000)
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
		fooComparer := func(i T, j T) (bool, error) {
			return i.(foo).str == j.(foo).str && i.(foo).num == j.(foo).num, nil
		}
		fooPtrComparer := func(i T, j T) (bool, error) {
			return i.(*foo).str == j.(*foo).str && i.(*foo).num == j.(*foo).num, nil
		}
		_ = fooPtrComparer

		erroneusComparer := func(i T, j T) (bool, error) {
			return false, genericError
		}

		Convey("Previous error is reflected on result", func() {
			_, err := From(allSameStruct).Where(erroneusBinaryFunc).DistinctBy(fooComparer).Results()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Provided func is nil", func() {
			_, err := From(allSameStruct).DistinctBy(nil).Results()
			So(err, ShouldEqual, ErrNilFunc)
		})
		Convey("Comparer returns error", func() {
			_, err := From(arr0).DistinctBy(erroneusComparer).Results()
			So(err, ShouldNotEqual, nil)
		})
		Convey("All elements are the same", func() {
			res, _ := From(allSameStruct).DistinctBy(fooComparer).Results()
			So(res, ShouldResemble, []T{allSameStruct[0]})
		})
		Convey("All elements are distinct", func() {
			var arr = make([]T, 100)
			for i := 0; i < len(arr); i++ {
				arr[i] = i
			}
			res, _ := From(arr).DistinctBy(func(this T, that T) (bool, error) {
				return this.(int) == that.(int), nil
			}).Results()
			So(res, ShouldResemble, arr)
		})
		Convey("Ensure leftmost appearance is returned in multiple occurrence cases", func() {
			arr := []T{&foo{"A", 0}, &foo{"B", 0}, &foo{"A", 0}, &foo{"C", 0},
				&foo{"A", 0}, &foo{"B", 0}}
			res, _ := From(arr).DistinctBy(fooPtrComparer).Results()
			So(len(res), ShouldResemble, 3)
			So(res[0], ShouldEqual, arr[0]) // A
			So(res[1], ShouldEqual, arr[1]) // B
			So(res[2], ShouldEqual, arr[3]) // C
		})
		Convey("Randomly generated integers with likely collisions", func() {
			var arr = make([]T, 10000)
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
			So(len(res), ShouldEqual, len(dict))
		})
	})
}

func TestUnion(t *testing.T) {
	uniqueArr0 := []T{1, 2, 3, 4, 5}
	uniqueArr1 := []T{"a", "b", "c"}
	allSameArr := []T{1, 1, 1, 1}
	sameStruct0 := []T{foo{"A", 0}, foo{"B", 0}}
	sameStruct1 := []T{foo{"B", 0}, foo{"A", 0}}
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(uniqueArr0).Where(erroneusBinaryFunc).Union(uniqueArr0).Results()
		So(err, ShouldNotEqual, nil)
	})
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

func TestIntersect(t *testing.T) {
	uniqueArr := []T{1, 2, 3, 4, 5}
	allSameArr := []T{1, 1, 1, 1}
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(uniqueArr).Where(erroneusBinaryFunc).Intersect(uniqueArr).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty ∩ nil", t, func() {
		_, err := From(empty).Intersect(nil).Results()
		So(err, ShouldEqual, ErrNilInput)
	})
	Convey("Empty ∩ empty", t, func() {
		res, _ := From(empty).Intersect(empty).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("Empty ∩ non-empty", t, func() {
		res, _ := From(empty).Intersect(uniqueArr).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("Non-empty ∩ empty", t, func() {
		res, _ := From(uniqueArr).Intersect(empty).Results()
		So(res, ShouldResemble, empty)
	})
	Convey("(Unique set) ∩ (itself)", t, func() {
		res, _ := From(uniqueArr).Intersect(uniqueArr).Results()
		So(res, ShouldResemble, uniqueArr)
	})
	Convey("(All same slice) ∩ (itself)", t, func() {
		res, _ := From(allSameArr).Intersect(allSameArr).Results()
		So(len(res), ShouldEqual, 1)
	})
	Convey("There is some intersection", t, func() {
		res, _ := From([]T{1, 2, 3, 4, 5}).Intersect([]T{3, 4, 5, 6, 7}).Results()
		So(res, ShouldResemble, []T{3, 4, 5})
	})
}

func TestExcept(t *testing.T) {
	uniqueArr := []T{1, 2, 3, 4, 5}
	allSameArr := []T{1, 1, 1, 1}
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(uniqueArr).Where(erroneusBinaryFunc).Except(uniqueArr).Results()
		So(err, ShouldNotEqual, nil)
	})
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
		res, _ := From([]T{1, 2, 3, 4, 5}).Except([]T{3, 4, 5, 6, 7}).Results()
		So(res, ShouldResemble, []T{1, 2})
	})
}

func TestCount(t *testing.T) {
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).CountBy(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
	})
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
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AnyWith(alwaysTrue)
		So(err, ShouldNotEqual, nil)
	})
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
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Single(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Single(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).Single(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		So(err, ShouldNotEqual, ErrNotSingle)
		_, err = From(arr0).Where(alwaysFalse).Single(erroneusBinaryFunc)
		So(err, ShouldEqual, ErrNotSingle)
	})
	Convey("No matches", t, func() {
		_, err := From(arr0).Single(alwaysFalse)
		So(err, ShouldEqual, ErrNotSingle)
	})
	Convey("All matches", t, func() {
		_, err := From(arr0).Single(alwaysTrue)
		So(err, ShouldEqual, ErrNotSingle)
	})
	Convey("Only one match", t, func() {
		match := 0
		var match0 = func(i T) (bool, error) {
			return i.(int) == match, nil
		}
		r, _ := From([]T{-1, -1, 0, 1, 1}).Single(match0)
		So(r, ShouldEqual, match)
		_, err := From([]T{0, 1, 2, 2, 0}).Single(match0)
		So(err, ShouldEqual, ErrNotSingle)
	})
}

func TestAll(t *testing.T) {
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).All(nil)
		So(err, ShouldNotEqual, nil)
	})
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
		match0 := func(i T) (bool, error) {
			return i.(int) == 0, nil
		}
		r, _ := From([]T{0, 1, 2, 2, 0}).All(match0)
		So(r, ShouldEqual, false)
	})
}

func TestElementAt_ElementAtOrNil(t *testing.T) {
	intArr := []T{1, 2, 3, 4, 5}
	Convey("empty.ElementAt is ErrNoElement", t, func() {
		_, err := From(empty).ElementAt(1)
		So(err, ShouldEqual, ErrNoElement)
	})
	Convey("empty.ElementAtOrNil is nil", t, func() {
		v, _ := From(empty).ElementAtOrNil(1)
		So(v, ShouldEqual, nil)
	})
	Convey("negative index returns is ErrNegativeParam", t, func() {
		_, err := From(empty).ElementAt(-1)
		So(err, ShouldEqual, ErrNegativeParam)
		_, err = From(empty).ElementAtOrNil(-1)
		So(err, ShouldEqual, ErrNegativeParam)
	})
	Convey("first element is returned", t, func() {
		v, _ := From(intArr).ElementAt(0)
		So(v, ShouldResemble, intArr[0])
		v, _ = From(intArr).ElementAtOrNil(0)
		So(v, ShouldResemble, intArr[0])
	})
	Convey("last element is returned", t, func() {
		v, _ := From(intArr).ElementAt(len(intArr) - 1)
		So(v, ShouldResemble, intArr[len(intArr)-1])
		v, _ = From(intArr).ElementAtOrNil(len(intArr) - 1)
		So(v, ShouldResemble, intArr[len(intArr)-1])
	})
	Convey("out of index returns ErrNoElement on non-empty slice", t, func() {
		_, err := From(intArr).ElementAt(len(intArr))
		So(err, ShouldEqual, ErrNoElement)
		_, err = From(intArr).ElementAtOrNil(len(intArr))
		So(err, ShouldEqual, nil)
	})
	Convey("previous errors are reflected", t, func() {
		_, err1 := From(arr0).Where(erroneusBinaryFunc).ElementAt(0)
		_, err2 := From(arr0).Where(erroneusBinaryFunc).ElementAtOrNil(0)
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
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
		v, _ = From(arr3).FirstOrNil()
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
	Convey("previous errors are reflected", t, func() {
		_, err1 := From(arr0).Where(erroneusBinaryFunc).FirstBy(alwaysTrue)
		_, err2 := From(arr0).Where(erroneusBinaryFunc).FirstOrNilBy(alwaysTrue)
		So(err1, ShouldNotEqual, nil)
		So(err2, ShouldNotEqual, nil)
	})
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
		v, _ = From(arr3).LastOrNil()
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
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).LastBy(alwaysTrue)
		So(err, ShouldNotEqual, nil)
	})

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
		arr := []T{1, 2, 3, 4, 5}
		rev := []T{5, 4, 3, 2, 1}
		res, _ := From(arr).Reverse().Results()
		So(res, ShouldResemble, rev)

		Convey("Slice containing nils", func() {
			arr := []T{1, nil, nil, 2, nil, 3, nil}
			rev := []T{nil, 3, nil, 2, nil, nil, 1}
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
		in := []T{1, 2, 3, 4, 5}
		res, _ := From(in).Take(3).Results()
		So(res, ShouldResemble, []T{1, 2, 3})
		Convey("Take n ≥ len(arr)", func() {
			res, _ := From(in).Take(len(in)).Results()
			So(res, ShouldResemble, res)
			res, _ = From(in).Take(len(in) + 1).Results()
			So(res, ShouldResemble, res)
		})
	})
}

func TestTakeWhile(t *testing.T) {
	Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).TakeWhile(alwaysTrue).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Nil func passed", t, func() {
		_, err := From(arr0).TakeWhile(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})
	Convey("Error returned from passed func is reflected in result", t, func() {
		_, err := From(arr0).TakeWhile(erroneusBinaryFunc).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty slice take all", t, func() {
		res, err := From(empty).TakeWhile(alwaysTrue).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})

	Convey("Take none", t, func() {
		res, _ := From(arr0).TakeWhile(alwaysFalse).Results()
		So(res, ShouldResemble, empty)
	})

	Convey("Take only first", t, func() {
		in := []T{1, 2, 3, 4, 5}
		res, err := From(in).TakeWhile(func(i T) (bool, error) { return i.(int) < 2, nil }).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, in[:1])
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
		in := []T{1, 2, 3, 4, 5}
		res, _ := From(in).Skip(3).Results()
		So(res, ShouldResemble, []T{4, 5})
		Convey("Skip n ≥ len(arr)", func() {
			res, _ := From(in).Skip(len(in)).Results()
			So(res, ShouldResemble, empty)
			res, _ = From(in).Skip(len(in) + 1).Results()
			So(res, ShouldResemble, empty)
		})
	})

	Convey("Skip & take & skip", t, func() {
		in := []T{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		res, _ := From(in).Skip(0).Skip(-1000).Skip(1).Take(1000).Take(5).Results()
		So(res, ShouldResemble, []T{1, 2, 3, 4, 5})
	})
}

func TestSkipWhile(t *testing.T) {
	Convey("Previous error is reflected in result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).SkipWhile(alwaysTrue).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Nil func passed", t, func() {
		_, err := From(arr0).SkipWhile(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})
	Convey("Error returned from passed func is reflected in result", t, func() {
		_, err := From(arr0).SkipWhile(erroneusBinaryFunc).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty slice Skip all", t, func() {
		res, err := From(empty).SkipWhile(alwaysTrue).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})

	Convey("Skip none", t, func() {
		res, _ := From(arr0).SkipWhile(alwaysFalse).Results()
		So(res, ShouldResemble, arr0)
	})

	Convey("Skip all", t, func() {
		res, _ := From(arr0).SkipWhile(alwaysTrue).Results()
		So(res, ShouldResemble, empty)
	})

	Convey("Skip only first", t, func() {
		in := []T{1, 2, 3, 4, 5}
		res, _ := From(in).SkipWhile(func(i T) (bool, error) { return i.(int) < 2, nil }).Results()
		So(res, ShouldResemble, in[1:])
	})

	Convey("SkipWhile & TakeWhile & SkipWhile", t, func() {
		in := []T{1, 2, 3, 4, 5, 6, 7, 8, 9}
		lessThanTwo := func(i T) (bool, error) { return i.(int) < 2, nil }
		lessThanSix := func(i T) (bool, error) { return i.(int) < 6, nil }
		res, _ := From(in).SkipWhile(alwaysFalse).SkipWhile(lessThanTwo).TakeWhile(lessThanSix).Results()
		So(res, ShouldResemble, []T{2, 3, 4, 5})
	})
}

func TestOrder(t *testing.T) {
	Convey("Sort ints", t, func() {
		arr := []T{6, 1, 4, 0, -1, 2}
		arrSorted := []T{-1, 0, 1, 2, 4, 6}
		unsupportedArr := []T{6, 1, 4, 0, -1, 2, ""}

		Convey("Previous error is reflected on result", func() {
			_, err := From(arr).Where(erroneusBinaryFunc).OrderInts().Results()
			So(err, ShouldNotEqual, nil)
		})

		Convey("Sort order is correct", func() {
			res, _ := From(arr).OrderInts().Results()
			So(res, ShouldResemble, arrSorted)
		})

		Convey("Sequence contain unsupported types", func() {
			_, err := From(unsupportedArr).OrderInts().Results()
			So(err, ShouldEqual, ErrTypeMismatch)
		})
	})

	Convey("Sort float64s", t, func() {
		arr := []T{1.000000001, 1.0000000001, 0.1, 0.01, 0.00001, 0.0000000000001}
		arrSorted := []T{0.0000000000001, 0.00001, 0.01, 0.1, 1.0000000001, 1.000000001}
		unsupportedArr := []T{1.000000001, "", 1.0000000001, 0.1, nil}

		Convey("Previous error is reflected on result", func() {
			_, err := From(arr).Where(erroneusBinaryFunc).OrderFloat64s().Results()
			So(err, ShouldNotEqual, nil)
		})

		Convey("Sort order is correct", func() {
			res, _ := From(arr).OrderFloat64s().Results()
			So(res, ShouldResemble, arrSorted)
		})

		Convey("Sequence contain unsupported types", func() {
			_, err := From(unsupportedArr).OrderFloat64s().Results()
			So(err, ShouldEqual, ErrTypeMismatch)
		})
	})

	Convey("Sort strings", t, func() {
		arr := []T{"c", "a", "", "aa", "b"}
		arrSorted := []T{"", "a", "aa", "b", "c"}

		unsupportedArr := []T{"", "aa", "ccc", nil}

		Convey("Previous error is reflected on result", func() {
			_, err := From(arr).Where(erroneusBinaryFunc).OrderStrings().Results()
			So(err, ShouldNotEqual, nil)
		})

		Convey("Sort order is correct", func() {
			res, _ := From(arr).OrderStrings().Results()
			So(res, ShouldResemble, arrSorted)
		})

		Convey("Sequence contain unsupported types", func() {
			_, err := From(unsupportedArr).OrderStrings().Results()
			So(err, ShouldEqual, ErrTypeMismatch)
		})
	})

}

func TestOrderBy(t *testing.T) {
	unsorted := []T{&foo{"A", 5}, &foo{"B", 1}, &foo{"C", 3}}
	sorted := []T{&foo{"B", 1}, &foo{"C", 3}, &foo{"A", 5}}
	sortByNum := func(this T, that T) bool {
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

	people := []T{magnus, terry, charlotte, ahmet}
	pets := []T{barley, boots, whiskers, daisy, sasha}

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

	Convey("Equi-join", t, func() {
		Convey("Errors from the previous of the chain are carried on", func() {
			_, err := From(people).Where(erroneusBinaryFunc).Join(pets, dummyKeySelector,
				dummyKeySelector, dummyResultSelector).Results()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Nil funcs passed", func() {
			_, err := From(people).Join(pets, nil, nil, nil).Results()
			So(err, ShouldEqual, ErrNilFunc)
		})
		Convey("Nil input passed", func() {
			_, err := From(people).Join(nil, dummyKeySelector, dummyKeySelector,
				dummyResultSelector).Results()
			So(err, ShouldEqual, ErrNilInput)
		})
		Convey("Pets & owners example join (also checks preserving the order)", func() {

			res, err := From(people).Join(pets,
				func(person T) T { return person.(Person).Name },
				func(pet T) T { return pet.(Pet).Owner.Name },
				func(outer T, inner T) T {
					return ResultPair{outer.(Person).Name, inner.(Pet).Name}
				}).Results()
			So(err, ShouldEqual, nil)
			So(res, ShouldResemble, equiJoinExpected)
		})
	})

	Convey("Group-join", t, func() {
		Convey("Errors from the previous of the chain are carried on", func() {
			_, err := From(people).Where(erroneusBinaryFunc).GroupJoin(pets, dummyKeySelector,
				dummyKeySelector, dummyGroupResultSelector).Results()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Nil funcs passed", func() {
			_, err := From(people).GroupJoin(pets, nil, nil, nil).Results()
			So(err, ShouldEqual, ErrNilFunc)
		})
		Convey("Nil input passed", func() {
			_, err := From(people).GroupJoin(nil, dummyKeySelector, dummyKeySelector,
				dummyGroupResultSelector).Results()
			So(err, ShouldEqual, ErrNilInput)
		})
		Convey("Pets & owners example join (also checks preserving the order)", func() {

			res, err := From(people).GroupJoin(pets,
				func(person T) T { return person.(Person).Name },
				func(pet T) T { return pet.(Pet).Owner.Name },
				func(outer T, inners []T) T {
					return ResultGroup{outer.(Person).Name, inners}
				}).Results()
			So(err, ShouldEqual, nil)
			So(res, ShouldResemble, groupJoinExpected)
		})
	})
}

func TestRange(t *testing.T) {
	Convey("count < 0", t, func() {
		_, err := Range(1, -1).Results()
		So(err, ShouldEqual, ErrNegativeParam)
	})
	Convey("count = 0", t, func() {
		res, err := Range(1, 0).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, empty)
	})
	Convey("range(1,10)", t, func() {
		res, err := Range(1, 10).Results()
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, []T{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})
}

var (
	intArr            = []T{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10}
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
	Convey("Previous errors are reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Sum()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty slice", t, func() {
		res, err := From(empty).Sum()
		So(err, ShouldEqual, nil)
		So(res, ShouldEqual, 0.0)
	})
	Convey("Slice of ints", t, func() {
		res, _ := From(intArr).Sum()
		So(res, ShouldEqual, intArrSumExpected)
	})
	Convey("Slice of mixed numeric types", t, func() {
		res, _ := From(mixedArr).Sum()
		So(res, ShouldAlmostEqual, mixedArrSumExpected, 0.000001) // float32 requires less tolerance than goconvey default
	})
	Convey("Slice with numeric types and nils", t, func() {
		_, err := From(mixedArrContainingNil).Sum()
		So(err, ShouldEqual, ErrNan)
	})
	Convey("Slice contains unsupported type", t, func() {
		_, err := From(mixedArrContainingUnsupported).Sum()
		So(err, ShouldNotEqual, nil)
	})
}

func TestAverage(t *testing.T) {
	Convey("Previous errors are reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).Average()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Empty slice", t, func() {
		_, err := From(empty).Average()
		So(err, ShouldEqual, ErrEmptySequence)
	})
	Convey("Slice of ints", t, func() {
		res, _ := From(intArr).Average()
		So(res, ShouldEqual, intArrAvgExpected)
	})
	Convey("Slice of mixed numeric types", t, func() {
		res, _ := From(mixedArr).Average()
		So(res, ShouldAlmostEqual, mixedArrAvgExpected, 0.000001) // float32 requires less tolerance than goconvey default
	})
	Convey("Slice with numeric types and nils", t, func() {
		_, err := From(mixedArrContainingNil).Average()
		So(err, ShouldEqual, ErrNan)
	})
	Convey("Slice contains unsupported type", t, func() {
		_, err := From(mixedArrContainingUnsupported).Average()
		So(err, ShouldNotEqual, nil)
	})
}

func TestMinMax(t *testing.T) {
	Convey("MinInt/MaxInt", t, func() {
		var (
			arr            = []T{-1, -9, 0, 9, 1}
			arrUnsupported = []T{-1, -9, 0, 9, 1, nil}
			expectedMin    = -9
			expectedMax    = 9
		)
		Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).MinInt()
			So(err, ShouldNotEqual, nil)
			_, err = From(arr0).Where(erroneusBinaryFunc).MaxInt()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Empty slice", func() {
			_, err := From(empty).MinInt()
			So(err, ShouldEqual, ErrEmptySequence)
			_, err = From(empty).MaxInt()
			So(err, ShouldEqual, ErrEmptySequence)
		})
		Convey("Sequence contains unsupported types", func() {
			_, err := From(arrUnsupported).MinInt()
			So(err, ShouldEqual, ErrTypeMismatch)
			_, err = From(arrUnsupported).MaxInt()
			So(err, ShouldEqual, ErrTypeMismatch)
		})
		Convey("Verify min/max result", func() {
			res, _ := From(arr).MinInt()
			So(res, ShouldEqual, expectedMin)
			res, _ = From(arr).MaxInt()
			So(res, ShouldEqual, expectedMax)
		})
	})
	Convey("MinUint/MaxUint", t, func() {
		var (
			arr            = []T{uint(1), uint(9), uint(100), uint(99), uint(0)}
			arrUnsupported = []T{uint(1), uint(9), uint(100), uint(99), uint(0), 0}
			expectedMin    = uint(0)
			expectedMax    = uint(100)
		)
		Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).MinUint()
			So(err, ShouldNotEqual, nil)
			_, err = From(arr0).Where(erroneusBinaryFunc).MaxUint()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Empty slice", func() {
			_, err := From(empty).MinUint()
			So(err, ShouldEqual, ErrEmptySequence)
			_, err = From(empty).MaxUint()
			So(err, ShouldEqual, ErrEmptySequence)
		})
		Convey("Sequence contains unsupported types", func() {
			_, err := From(arrUnsupported).MinUint()
			So(err, ShouldEqual, ErrTypeMismatch)
			_, err = From(arrUnsupported).MaxUint()
			So(err, ShouldEqual, ErrTypeMismatch)
		})
		Convey("Verify min/max result", func() {
			res, _ := From(arr).MinUint()
			So(res, ShouldEqual, expectedMin)
			res, _ = From(arr).MaxUint()
			So(res, ShouldEqual, expectedMax)
		})
	})
	Convey("MinFloat64/MaxFloat64", t, func() {
		var (
			arr            = []T{float64(-9), float64(-9.9), float64(0), float64(99), float64(99.9)}
			arrUnsupported = []T{float64(-9), float64(-9.9), float64(0), float64(99), float64(99.9), uint(0)}
			expectedMin    = float64(-9.9)
			expectedMax    = float64(99.9)
		)
		Convey("Previous error is reflected on result", func() {
			_, err := From(arr0).Where(erroneusBinaryFunc).MinFloat64()
			So(err, ShouldNotEqual, nil)
			_, err = From(arr0).Where(erroneusBinaryFunc).MaxFloat64()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Empty slice", func() {
			_, err := From(empty).MinFloat64()
			So(err, ShouldEqual, ErrEmptySequence)
			_, err = From(empty).MaxFloat64()
			So(err, ShouldEqual, ErrEmptySequence)
		})
		Convey("Sequence contains unsupported types", func() {
			_, err := From(arrUnsupported).MinFloat64()
			So(err, ShouldEqual, ErrTypeMismatch)
			_, err = From(arrUnsupported).MaxFloat64()
			So(err, ShouldEqual, ErrTypeMismatch)
		})
		Convey("Verify min/max result", func() {
			res, _ := From(arr).MinFloat64()
			So(res, ShouldEqual, expectedMin)
			res, _ = From(arr).MaxFloat64()
			So(res, ShouldEqual, expectedMax)
		})
	})
}
