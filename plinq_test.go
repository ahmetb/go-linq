package linq

import (
	c "github.com/smartystreets/goconvey/convey"
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

func TestResultsParallel(t *testing.T) {
	c.Convey("If error exists in given queryable, error is returned", t, func() {
		q := Query{
			values: nil,
			err:    errFoo}
		_, err := q.AsParallel().Results()
		c.So(err, c.ShouldEqual, errFoo)
	})
	c.Convey("Given no errors exist, non-nil results are returned", t, func() {
		q := From(arr0)
		val, err := q.AsParallel().Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(val, shouldSlicesResemble, arr0)
	})
	c.Convey("Returned result is isolated (copied) from original query source", t, func() {
		// Regression for BUG: modifying result slice effects subsequent query methods
		arr := []int{1, 2, 3, 4, 5}
		q := From(arr)
		res, _ := q.AsParallel().Results()
		res[0] = 100
		sum, _ := q.Sum()
		c.So(sum, c.ShouldEqual, 15)
	})
}

func TestWhereParallel(t *testing.T) {
	c.Convey("Chose none of the elements", t, func() {
		val, _ := From(arr0).AsParallel().Where(alwaysFalse).Results()
		c.So(len(val), c.ShouldEqual, 0)
	})
	c.Convey("A previous error is reflected on the result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().Where(alwaysTrue).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on the result", t, func() {
		_, err := From(arr0).AsParallel().Where(erroneusBinaryFunc).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Nil func passed", t, func() {
		_, err := From(arr0).AsParallel().Where(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
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
	c.Convey("Do not preserve order", t, func() {
		c.Convey("Chose all elements, as is", func() {
			q := From(arr).AsParallel().AsUnordered().Where(alwaysTrueDelayed)
			val, _ := q.Results()
			sum, _ := q.AsSequential().Sum()
			c.So(len(val), c.ShouldEqual, len(arr))
			c.So(sum, c.ShouldEqual, 499500)
		})
		c.Convey("Basic filtering (x mod 2)==0", func() {
			q := From(arr).AsParallel().AsUnordered().Where(divisibleBy2Delayed)
			val, err := q.Results()
			c.So(len(val), c.ShouldEqual, n/2)
			sum, _ := q.AsSequential().Sum()
			c.So(err, c.ShouldEqual, nil)
			c.So(sum, c.ShouldEqual, float64(249500))
		})
	})

	c.Convey("Preserve order", t, func() {
		c.Convey("Chose all elements, as is", func() {
			val, _ := From(arr).AsParallel().AsOrdered().Where(alwaysTrueDelayed).Results()
			c.So(val, shouldSlicesResemble, arr)
		})
		c.Convey("Basic filtering (x mod 2)==0", func() {
			q := From(arr).AsParallel().AsOrdered().Where(divisibleBy2Delayed)
			val, err := q.Results()
			c.So(len(val), c.ShouldEqual, n/2)
			sum, _ := q.AsSequential().Sum()
			c.So(err, c.ShouldEqual, nil)
			c.So(sum, c.ShouldEqual, float64(249500))
		})
	})

}

func TestSelectParallel(t *testing.T) {
	asIs := func(i T) (T, error) {
		return i, nil
	}
	erroneusFunc := func(i T) (T, error) {
		return nil, errFoo
	}

	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().Select(asIs).Results()
		c.So(err, c.ShouldNotEqual, nil)
	})

	c.Convey("Nil func returns error", t, func() {
		_, err := From(arr0).AsParallel().Select(nil).Results()
		c.So(err, c.ShouldEqual, ErrNilFunc)
	})

	c.Convey("Error returned from provided func", t, func() {
		val, err := From(arr0).AsParallel().Select(erroneusFunc).Results()
		c.So(err, c.ShouldNotEqual, nil)

		c.Convey("Erroneus function is in chain with as-is select", func() {
			_, err = From(arr0).AsParallel().Select(asIs).Select(erroneusFunc).Results()
			c.So(err, c.ShouldNotEqual, nil)
		})
		c.Convey("Erroneus function is in chain but not called", func() {
			val, err = From(arr0).Where(alwaysFalse).AsParallel().Select(erroneusFunc).Results()
			c.So(err, c.ShouldEqual, nil)
			c.So(len(val), c.ShouldEqual, 0)
		})

	})

	c.Convey("Select all elements as is", t, func() {
		val, err := From(arr0).AsParallel().Select(asIs).Results()
		c.So(err, c.ShouldEqual, nil)
		c.So(val, c.ShouldNotEqual, arr0)
		c.So(val, shouldSlicesResemble, arr0)
	})

	c.Convey("Pow(x,2) for i in []int", t, func() {
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
		c.So(val, shouldSlicesResemble, expected)
	})
}

func TestAnyWithParallel(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().AnyWith(alwaysTrueDelayed)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrueDelayed).AsParallel().AnyWith(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrueDelayed).AsParallel().AnyWith(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).AsParallel().AnyWith(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("No matches", t, func() {
		r, _ := From(arr0).AsParallel().AnyWith(alwaysFalseDelayed)
		c.So(r, c.ShouldEqual, false)
		r, _ = From(arr0).AsParallel().Where(alwaysFalseDelayed).Any()
		c.So(r, c.ShouldEqual, false)
	})
	c.Convey("All matches", t, func() {
		r, _ := From(arr0).AsParallel().AnyWith(alwaysTrueDelayed)
		c.So(r, c.ShouldEqual, true)
		r, _ = From(arr0).AsParallel().Where(alwaysTrueDelayed).Any()
		c.So(r, c.ShouldEqual, true)
	})
}

func TestAllParallel(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).Where(erroneusBinaryFunc).AsParallel().All(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AsParallel().All(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).Where(alwaysTrue).AsParallel().All(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		_, err = From(arr0).Where(alwaysFalse).AsParallel().All(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("Empty slice", t, func() {
		r, _ := From(empty).AsParallel().All(alwaysTrueDelayed)
		c.So(r, c.ShouldEqual, true)
	})
	c.Convey("No matches", t, func() {
		r, _ := From(arr0).AsParallel().All(alwaysFalseDelayed)
		c.So(r, c.ShouldEqual, false)
	})
	c.Convey("All matches", t, func() {
		r, _ := From(arr0).AsParallel().All(alwaysTrueDelayed)
		c.So(r, c.ShouldEqual, true)
	})
	c.Convey("Multiple matches", t, func() {
		match0 := func(i T) (bool, error) {
			return i.(int) == 0, nil
		}
		r, _ := From([]T{0, 1, 2, 2, 0}).AsParallel().All(match0)
		c.So(r, c.ShouldEqual, false)
	})
}

func TestSingleParallel(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).AsParallel().Where(erroneusBinaryFunc).Single(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).AsParallel().Where(alwaysTrueDelayed).Single(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).AsParallel().Where(alwaysTrueDelayed).Single(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		c.So(err, c.ShouldNotEqual, ErrNotSingle)
		_, err = From(arr0).AsParallel().Where(alwaysFalseDelayed).Single(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
	c.Convey("No matches", t, func() {
		_, err := From(arr0).AsParallel().Single(alwaysFalseDelayed)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
	c.Convey("All matches", t, func() {
		_, err := From(arr0).AsParallel().Single(alwaysTrueDelayed)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
	c.Convey("Only one match", t, func() {
		match := 0
		var match0 = func(i T) (bool, error) {
			return i.(int) == match, nil
		}
		r, _ := From([]T{-1, -1, 0, 1, 1}).AsParallel().Single(match0)
		c.So(r, c.ShouldEqual, match)
		_, err := From([]T{0, 1, 2, 2, 0}).AsParallel().Single(match0)
		c.So(err, c.ShouldEqual, ErrNotSingle)
	})
}

func TestCountParallel(t *testing.T) {
	c.Convey("Previous error is reflected on result", t, func() {
		_, err := From(arr0).AsParallel().Where(erroneusBinaryFunc).CountBy(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("Given a nil function, ErrNilFunc is returned", t, func() {
		_, err := From(arr0).AsParallel().Where(alwaysTrueDelayed).CountBy(nil)
		c.So(err, c.ShouldNotEqual, nil)
	})
	c.Convey("An error returned from f is reflected on Result", t, func() {
		_, err := From(arr0).AsParallel().Where(alwaysTrueDelayed).CountBy(erroneusBinaryFunc)
		c.So(err, c.ShouldNotEqual, nil)
		_, err = From(arr0).AsParallel().Where(alwaysFalseDelayed).CountBy(erroneusBinaryFunc)
		c.So(err, c.ShouldEqual, nil)
	})
	c.Convey("No matches", t, func() {
		cnt, _ := From(arr0).AsParallel().CountBy(alwaysFalseDelayed)
		c.So(cnt, c.ShouldEqual, 0)
		cnt, _ = From(arr0).AsParallel().Where(alwaysFalseDelayed).Count()
		c.So(cnt, c.ShouldEqual, 0)
	})
	c.Convey("All matches", t, func() {
		cnt, _ := From(arr0).AsParallel().CountBy(alwaysTrueDelayed)
		c.So(cnt, c.ShouldEqual, len(arr0))
		cnt, _ = From(arr0).AsParallel().Count()
		c.So(cnt, c.ShouldEqual, len(arr0))
	})
}
