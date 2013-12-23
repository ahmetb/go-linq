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
		val, _ := From(arr0).AsParallel().Where(alwaysFalse, false).Results()
		So(val, ShouldEqual, nil)
	})
	Convey("A previous error is reflected on the result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().Where(alwaysTrue, false).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("An error returned from f is reflected on the result", t, func() {
		_, err := From(arr0).AsParallel().Where(erroneusBinaryFunc, false).Results()
		So(err, ShouldNotEqual, nil)
	})
	Convey("Nil func passed", t, func() {
		_, err := From(arr0).AsParallel().Where(nil, false).Results()
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
			q := From(arr).AsParallel().Where(alwaysTrueDelayed, false)
			val, _ := q.Results()
			sum, _ := q.AsSequential().Sum()
			So(len(val), ShouldEqual, len(arr))
			So(sum, ShouldEqual, 499500)
		})
		Convey("Basic filtering (x mod 2)==0", func() {
			q := From(arr).AsParallel().Where(divisibleBy2Delayed, false)
			val, err := q.Results()
			So(len(val), ShouldEqual, n/2)
			sum, _ := q.AsSequential().Sum()
			So(err, ShouldEqual, nil)
			So(sum, ShouldEqual, float64(249500))
		})
	})

	Convey("Preserve order", t, func() {
		Convey("Chose all elements, as is", func() {
			val, _ := From(arr).AsParallel().Where(alwaysTrueDelayed, true).Results()
			So(val, ShouldResemble, arr)
		})
		Convey("Basic filtering (x mod 2)==0", func() {
			q := From(arr).AsParallel().Where(divisibleBy2Delayed, true)
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
