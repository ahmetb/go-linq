package linq

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

var (
	alwaysTrueDelayed = func(i T) (bool, error) {
		time.Sleep(time.Duration(rand.Intn(maxRandDelayMs)) * time.Millisecond)
		return true, nil
	}
	alwaysFalseDelayed = func(i T) (bool, error) {
		time.Sleep(time.Duration(rand.Intn(maxRandDelayMs)) * time.Millisecond)
		return false, nil
	}
)

func TestWhereParallel(t *testing.T) {
	Convey("Chose none of the elements", t, func() {
		val, _ := From(arr0).AsParallel().Where(alwaysFalse).Results()
		So(val, ShouldEqual, nil)
	})
	Convey("A previous error is reflected on the result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().Where(alwaysTrue).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on the result", t, func() {
		_, err := From(arr0).AsParallel().Where(erroneusBinaryFunc).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Nil func passed", t, func() {
		_, err := From(arr0).AsParallel().Where(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})

	n := 1000
	arr := make([]T, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}
	divisibleBy2Delayed := func(i T) (bool, error) {
		time.Sleep(time.Duration(rand.Intn(maxRandDelayMs)) * time.Millisecond)
		return i.(int)%2 == 0, nil
	}
	Convey("Do not preserve order", t, func() {
		Convey("Chose all elements, as is", func() {
			q := From(arr).AsParallel().AsUnordered().Where(alwaysTrueDelayed)
			val, _ := q.Results()
			sum, _ := q.AsSequential().Sum()
			So(len(val), ShouldEqual, len(arr))
			So(sum, ShouldEqual, 499500)
		})
		Convey("Basic filtering (x mod 2)==0", func() {
			q := From(arr).AsParallel().AsUnordered().Where(divisibleBy2Delayed)
			val, err := q.Results()
			So(len(val), ShouldEqual, n/2)
			sum, _ := q.AsSequential().Sum()
			So(err, ShouldEqual, nil)
			So(sum, ShouldEqual, float64(249500))
		})
	})

	Convey("Preserve order", t, func() {
		Convey("Chose all elements, as is", func() {
			val, _ := From(arr).AsParallel().AsOrdered().Where(alwaysTrueDelayed).Results()
			So(val, ShouldResemble, arr)
		})
		Convey("Basic filtering (x mod 2)==0", func() {
			q := From(arr).AsParallel().AsOrdered().Where(divisibleBy2Delayed)
			val, err := q.Results()
			So(len(val), ShouldEqual, n/2)
			sum, _ := q.AsSequential().Sum()
			So(err, ShouldEqual, nil)
			So(sum, ShouldEqual, float64(249500))
		})
	})

}

func TestSelectParallel(t *testing.T) {
	asIs := func(i T) (T, error) {
		return i, nil
	}
	erroneusFunc := func(i T) (T, error) {
		return nil, genericError
	}

	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().Select(asIs).Results()
		So(err, ShouldNotEqual, nil)
	})

	Convey("Nil func returns error", t, func() {
		_, err := From(arr0).AsParallel().Select(nil).Results()
		So(err, ShouldEqual, ErrNilFunc)
	})

	Convey("Error returned from provided func", t, func() {
		val, err := From(arr0).AsParallel().Select(erroneusFunc).Results()
		So(err, ShouldNotEqual, nil)

		Convey("Erroneus function is in chain with as-is select", func() {
			_, err = From(arr0).AsParallel().Select(asIs).Select(erroneusFunc).Results()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Erroneus function is in chain but not called", func() {
			val, err = From(arr0).Where(alwaysFalse).AsParallel().Select(erroneusFunc).Results()
			So(err, ShouldEqual, nil)
			So(len(val), ShouldEqual, 0)
		})

	})

	Convey("Select all elements as is", t, func() {
		val, err := From(arr0).AsParallel().Select(asIs).Results()
		So(err, ShouldEqual, nil)
		So(val, ShouldResemble, arr0)
	})

	Convey("Pow(x,2) for i in []int", t, func() {
		n := 10
		arr := make([]T, n)
		expected := make([]T, n)
		for i := 0; i < n; i++ {
			r := rand.Intn(999)
			arr[i] = r * r
		}
		for j, i := range arr {
			expected[j] = i.(int) * i.(int)
		}

		slowPow := func(i T) (T, error) {
			time.Sleep(time.Duration(maxRandDelayMs) * time.Millisecond)
			return i.(int) * i.(int), nil
		}
		val, _ := From(arr).AsParallel().Select(slowPow).Results()
		So(val, ShouldResemble, expected)
	})
}

func TestAnyWithParallel(t *testing.T) {
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().AnyWith(alwaysTrueDelayed)
		So(err, ShouldNotEqual, nil)
	})
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrueDelayed).AsParallel().AnyWith(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrueDelayed).AsParallel().AnyWith(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).AsParallel().AnyWith(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("No matches", t, func() {
		r, _ := From(arr0).AsParallel().AnyWith(alwaysFalseDelayed)
		So(r, ShouldEqual, false)
	})
	Convey("All matches", t, func() {
		r, _ := From(arr0).AsParallel().AnyWith(alwaysTrueDelayed)
		So(r, ShouldEqual, true)
	})
}

func TestAllParallel(t *testing.T) {
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().All(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AsParallel().All(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AsParallel().All(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).AsParallel().All(erroneusBinaryFunc)
		So(err, ShouldEqual, nil)
	})
	Convey("Empty slice", t, func() {
		r, _ := From(empty).AsParallel().All(alwaysTrueDelayed)
		So(r, ShouldEqual, true)
	})
	Convey("No matches", t, func() {
		r, _ := From(arr0).AsParallel().All(alwaysFalseDelayed)
		So(r, ShouldEqual, false)
	})
	Convey("All matches", t, func() {
		r, _ := From(arr0).AsParallel().All(alwaysTrueDelayed)
		So(r, ShouldEqual, true)
	})
	Convey("Multiple matches", t, func() {
		match0 := func(i T) (bool, error) {
			return i.(int) == 0, nil
		}
		r, _ := From([]T{0, 1, 2, 2, 0}).AsParallel().All(match0)
		So(r, ShouldEqual, false)
	})
}

func TestSingleParallel(t *testing.T) {
	Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().Single(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AsParallel().Single(nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AsParallel().Single(erroneusBinaryFunc)
		So(err, ShouldNotEqual, nil)
		So(err, ShouldNotEqual, ErrNotSingle)
		_, err = From(arr0).Where(alwaysFalse).AsParallel().Single(erroneusBinaryFunc)
		So(err, ShouldEqual, ErrNotSingle)
	})
	Convey("No matches", t, func() {
		_, err := From(arr0).AsParallel().Single(alwaysFalse)
		So(err, ShouldEqual, ErrNotSingle)
	})
	Convey("All matches", t, func() {
		_, err := From(arr0).AsParallel().Single(alwaysTrue)
		So(err, ShouldEqual, ErrNotSingle)
	})
	Convey("Only one match", t, func() {
		match := 0
		var match0 = func(i T) (bool, error) {
			return i.(int) == match, nil
		}
		r, _ := From([]T{-1, -1, 0, 1, 1}).AsParallel().Single(match0)
		So(r, ShouldEqual, match)
		_, err := From([]T{0, 1, 2, 2, 0}).AsParallel().Single(match0)
		So(err, ShouldEqual, ErrNotSingle)
	})
}
