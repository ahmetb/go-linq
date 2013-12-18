package linq

import (
	"errors"
	"sort"
)

type queryable struct {
	values []interface{}
	err    error
	less   func(this, that interface{}) bool
}

func (q queryable) Len() int           { return len(q.values) }
func (q queryable) Swap(i, j int)      { q.values[i], q.values[j] = q.values[j], q.values[i] }
func (q queryable) Less(i, j int) bool { return q.less(q.values[i], q.values[j]) }

var (
	ErrNilFunc         = errors.New("linq: passed evaluation function is nil")
	ErrNilInput        = errors.New("linq: nil input passed to From")
	ErrNoElement       = errors.New("linq: element satisfying the conditions does not exist")
	ErrNegativeParam   = errors.New("linq: parameter cannot be negative")
	ErrUnsupportedType = errors.New("linq: sorting this type with Order is not supported, use OrderBy")
)

func From(input []interface{}) queryable {
	var _err error
	if input == nil {
		_err = ErrNilInput
	}
	return queryable{
		values: input,
		err:    _err}
}

func (q queryable) Results() ([]interface{}, error) {
	return q.values, q.err
}

func (q queryable) Where(f func(interface{}) (bool, error)) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return r
	}
	if f == nil {
		r.err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		ok, err := f(i)
		if err != nil {
			r.err = err
			return r
		}
		if ok {
			r.values = append(r.values, i)
		}
	}
	return r
}

func (q queryable) Select(f func(interface{}) (interface{}, error)) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return r
	}
	if f == nil {
		r.err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		val, err := f(i)
		if err != nil {
			r.err = err
			return r
		}
		r.values = append(r.values, val)
	}
	return
}

func (q queryable) Distinct() (r queryable) {
	return q.distinct(nil)
}

func (q queryable) DistinctBy(f func(interface{}, interface{}) (bool, error)) (r queryable) {
	if f == nil {
		r.err = ErrNilFunc
		return
	}
	return q.distinct(f)
}

func (q queryable) distinct(f func(interface{}, interface{}) (bool, error)) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return r
	}

	if f == nil {
		// basic equality comparison using dict
		dict := make(map[interface{}]bool)
		for _, v := range q.values {
			if _, ok := dict[v]; !ok {
				dict[v] = true
			}
		}
		res := make([]interface{}, len(dict))
		i := 0
		for key, _ := range dict {
			res[i] = key
			i++
		}
		r.values = res
	} else {
		// use equality comparer and bool flags for each item
		// here we check all a[i]==a[j] i<j, practically worst case
		// for this is O(N^2) where all elements are different and best case
		// is O(N) where all elements are the same
		// pick lefthand side value of the comparison in the result
		l := len(q.values)
		results := make([]interface{}, 0)
		included := make([]bool, l)
		for i := 0; i < l; i++ {
			if included[i] {
				continue
			}
			for j := i + 1; j < l; j++ {
				equals, err := f(q.values[i], q.values[j])
				if err != nil {
					r.err = err
					return
				}
				if equals {
					included[j] = true // don't include righthand side value
				}
			}
			results = append(results, q.values[i])
		}
		r.values = results
	}
	return
}

func (q queryable) Union(in []interface{}) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if in == nil {
		r.err = ErrNilInput
		return
	}
	var set map[interface{}]bool = make(map[interface{}]bool)

	for _, v := range q.values {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	for _, v := range in {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	r.values = make([]interface{}, len(set))
	i := 0
	for k, _ := range set {
		r.values[i] = k
		i++
	}
	return
}

func (q queryable) Intersect(in []interface{}) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if in == nil {
		r.err = ErrNilInput
		return
	}
	var set map[interface{}]bool = make(map[interface{}]bool)
	var intersection map[interface{}]bool = make(map[interface{}]bool)

	for _, v := range q.values {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	for _, v := range in {
		if _, ok := set[v]; ok {
			delete(set, v)
			if _, added := intersection[v]; !added {
				intersection[v] = true
			}
		}
	}
	r.values = make([]interface{}, len(intersection))
	i := 0
	for k, _ := range intersection {
		r.values[i] = k
		i++
	}
	return
}

func (q queryable) Except(in []interface{}) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if in == nil {
		r.err = ErrNilInput
		return
	}
	var set map[interface{}]bool = make(map[interface{}]bool)

	for _, v := range q.values {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	for _, v := range in {
		delete(set, v)
	}
	r.values = make([]interface{}, len(set))
	i := 0
	for k, _ := range set {
		r.values[i] = k
		i++
	}
	return
}

func (q queryable) Count() (count int, err error) {
	return len(q.values), q.err
}

func (q queryable) CountBy(f func(interface{}) (bool, error)) (c int, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		if ok {
			c++
		}
	}
	return
}

func (q queryable) Any() (exists bool, err error) {
	return len(q.values) > 0, q.err
}

func (q queryable) AnyWith(f func(interface{}) (bool, error)) (exists bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		if ok {
			exists = true
			return
		}
	}
	return
}

func (q queryable) All(f func(interface{}) (bool, error)) (all bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	all = true // if no elements, result is true
	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		all = all && ok
	}
	return
}

func (q queryable) Single(f func(interface{}) (bool, error)) (single bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	count, e := q.CountBy(f)
	if e != nil {
		err = e
		return
	}
	single = count == 1
	return
}

func (q queryable) First() (elem interface{}, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		err = ErrNoElement
	} else {
		elem = q.values[0]
	}
	return
}

func (q queryable) FirstOrNil() (elem interface{}, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) > 0 {
		elem = q.values[0]
	}
	return
}

func (q queryable) firstBy(f func(interface{}) (bool, error)) (elem interface{}, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		if ok {
			elem = i
			found = true
			break
		}
	}
	return
}

func (q queryable) FirstBy(f func(interface{}) (bool, error)) (elem interface{}, err error) {
	var found bool
	elem, found, err = q.firstBy(f)

	if err == nil && !found {
		err = ErrNoElement
	}
	return
}

func (q queryable) FirstOrNilBy(f func(interface{}) (bool, error)) (elem interface{}, err error) {
	elem, found, err := q.firstBy(f)
	if !found {
		elem = nil
	}
	return
}

func (q queryable) Last() (elem interface{}, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		err = ErrNoElement
	} else {
		elem = q.values[len(q.values)-1]
	}
	return
}

func (q queryable) LastOrNil() (elem interface{}, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) > 0 {
		elem = q.values[len(q.values)-1]
	}
	return
}

func (q queryable) lastBy(f func(interface{}) (bool, error)) (elem interface{}, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	for i := len(q.values) - 1; i >= 0; i-- {
		item := q.values[i]
		ok, e := f(item)
		if e != nil {
			err = e
			return
		}
		if ok {
			elem = item
			found = true
			break
		}
	}
	return
}

func (q queryable) LastBy(f func(interface{}) (bool, error)) (elem interface{}, err error) {
	var found bool
	elem, found, err = q.lastBy(f)

	if err == nil && !found {
		err = ErrNoElement
	}
	return
}

func (q queryable) LastOrNilBy(f func(interface{}) (bool, error)) (elem interface{}, err error) {
	elem, found, err := q.lastBy(f)
	if !found {
		elem = nil
	}
	return
}

func (q queryable) Reverse() (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	c := len(q.values)
	j := 0
	r.values = make([]interface{}, c)
	for i := c - 1; i >= 0; i-- {
		r.values[j] = q.values[i]
		j++
	}
	return
}

func (q queryable) Take(n int) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if n < 0 {
		n = 0
	}
	if n >= len(q.values) {
		n = len(q.values)
	}
	r.values = q.values[:n]
	return
}

func (q queryable) TakeWhile(f func(interface{}) (bool, error)) (r queryable) {
	n, err := q.findWhileTerminationIndex(f)
	if err != nil {
		r.err = err
		return
	}
	return q.Take(n)
}

func (q queryable) Skip(n int) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if n < 0 {
		n = 0
	}
	if n >= len(q.values) {
		n = len(q.values)
	}
	r.values = q.values[n:]
	return
}

func (q queryable) SkipWhile(f func(interface{}) (bool, error)) (r queryable) {
	n, err := q.findWhileTerminationIndex(f)
	if err != nil {
		r.err = err
		return
	}
	return q.Skip(n)
}

func (q queryable) findWhileTerminationIndex(f func(interface{}) (bool, error)) (n int, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	n = 0
	for _, v := range q.values {
		ok, e := f(v)
		if e != nil {
			err = e
			return
		}
		if ok {
			n++
		} else {
			break
		}
	}
	return
}

//TODO document: only sorts int, string, float64
func (q queryable) Order() (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}

	if len(q.values) > 0 {
		f := q.values[0]
		if _, ints := f.(int); ints {
			vals := toInts(q.values)
			sort.Ints(vals)
			r.values = intsToInterface(vals)
		} else if _, strings := f.(string); strings {
			vals := toStrings(q.values)
			sort.Strings(vals)
			r.values = stringsToInterface(vals)
		} else if _, float64s := f.(float64); float64s {
			vals := toFloat64s(q.values)
			sort.Float64s(vals)
			r.values = float64sToInterface(vals)
		} else {
			r.err = ErrUnsupportedType
		}
	}
	return
}

func (q queryable) OrderBy(less func(this interface{}, that interface{}) bool) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if less == nil {
		r.err = ErrNilFunc
		return
	}
	r.less = less
	r.values = make([]interface{}, len(q.values))
	_ = copy(r.values, q.values)
	sort.Sort(r)
	return
}

func (q queryable) Join(innerCollection []interface{},
	outerKeySelector func(interface{}) interface{},
	innerKeySelector func(interface{}) interface{},
	resultSelector func(
		outer interface{},
		inner interface{}) interface{}) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if innerCollection == nil {
		r.err = ErrNilInput
		return
	}
	if outerKeySelector == nil || innerKeySelector == nil || resultSelector == nil {
		r.err = ErrNilFunc
		return
	}
	var outerCollection = q.values
	innerKeyLookup := make(map[interface{}]interface{})

	for _, outer := range outerCollection {
		outerKey := outerKeySelector(outer)
		for _, inner := range innerCollection {
			innerKey, ok := innerKeyLookup[inner]
			if !ok {
				innerKey = innerKeySelector(inner)
				innerKeyLookup[inner] = innerKey
			}
			if innerKey == outerKey {
				elem := resultSelector(outer, inner)
				r.values = append(r.values, elem)
			}
		}
	}
	return
}

func (q queryable) GroupJoin(innerCollection []interface{},
	outerKeySelector func(interface{}) interface{},
	innerKeySelector func(interface{}) interface{},
	resultSelector func(
		outer interface{},
		inners []interface{}) interface{}) (r queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if innerCollection == nil {
		r.err = ErrNilInput
		return
	}
	if outerKeySelector == nil || innerKeySelector == nil || resultSelector == nil {
		r.err = ErrNilFunc
		return
	}
	var outerCollection = q.values
	innerKeyLookup := make(map[interface{}]interface{})

	var results = make(map[interface{}][]interface{}) // outer --> inner...
	for _, outer := range outerCollection {
		outerKey := outerKeySelector(outer)
		bucket := make([]interface{}, 0)
		results[outer] = bucket
		for _, inner := range innerCollection {
			innerKey, ok := innerKeyLookup[inner]
			if !ok {
				innerKey = innerKeySelector(inner)
				innerKeyLookup[inner] = innerKey
			}
			if innerKey == outerKey {
				results[outer] = append(results[outer], inner)
			}
		}
	}

	r.values = make([]interface{}, len(results))
	i := 0
	for k, v := range results {
		outer := k
		inners := v
		r.values[i] = resultSelector(outer, inners)
		i++
	}
	return
}

//TODO document integer oveflows are not handled
func Range(start, count int) (q queryable) {
	if count < 0 {
		q.err = ErrNegativeParam
		return
	}
	q.values = make([]interface{}, count)
	for i := 0; i < count; i++ {
		q.values[i] = start + i
	}
	return
}
