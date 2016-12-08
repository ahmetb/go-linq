package linq

// AllT is the typed version of All.
//
// NOTE: All method has better performance than AllT
//
// predicateFn is of a type "func(TSource) bool"
//
func (q Query) AllT(predicateFn interface{}) bool {

	predicateGenericFunc, err := newGenericFunc(
		"AllT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}
	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.All(predicateFunc)
}

// AnyWithT is the typed version of AnyWith.
//
// NOTE: AnyWith method has better performance than AnyWithT
//
// predicateFn is of a type "func(TSource) bool"
//
func (q Query) AnyWithT(predicateFn interface{}) bool {

	predicateGenericFunc, err := newGenericFunc(
		"AnyWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.AnyWith(predicateFunc)
}

// CountWithT is the typed version of CountWith.
//
// NOTE: CountWith method has better performance than CountWithT
//
// predicateFn is of a type "func(TSource) bool"
//
func (q Query) CountWithT(predicateFn interface{}) int {

	predicateGenericFunc, err := newGenericFunc(
		"CountWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.CountWith(predicateFunc)
}

// FirstWithT is the typed version of FirstWith.
//
// NOTE: FirstWith method has better performance than FirstWithT
//
// predicateFn is of a type "func(TSource) bool"
//
func (q Query) FirstWithT(predicateFn interface{}) interface{} {

	predicateGenericFunc, err := newGenericFunc(
		"FirstWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.FirstWith(predicateFunc)
}

// LastWithT is the typed version of LastWith.
//
// NOTE: LastWith method has better performance than LastWithT
//
// predicateFn is of a type "func(TSource) bool"
//
func (q Query) LastWithT(predicateFn interface{}) interface{} {

	predicateGenericFunc, err := newGenericFunc(
		"LastWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.LastWith(predicateFunc)
}

// SingleWithT is the typed version of SingleWith.
//
// NOTE: SingleWith method has better performance than SingleWithT
//
// predicateFn is of a type "func(TSource) bool"
//
func (q Query) SingleWithT(predicateFn interface{}) interface{} {

	predicateGenericFunc, err := newGenericFunc(
		"SingleWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.SingleWith(predicateFunc)
}

// ToMapByT is the typed version of ToMapBy.
//
// NOTE: ToMapBy method has better performance than ToMapByT
//
// keySelectorFn is of a type "func(TSource)TKey"
//
// valueSelectorFn is of a type "func(TSource)TValue"
//
func (q Query) ToMapByT(result interface{}, keySelectorFn interface{}, valueSelectorFn interface{}) {

	keySelectorGenericFunc, err := newGenericFunc(
		"ToMapByT", "keySelectorFn", keySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	keySelectorFunc := func(item interface{}) interface{} {
		return keySelectorGenericFunc.Call(item)
	}

	valueSelectorGenericFunc, err := newGenericFunc(
		"ToMapByT", "valueSelectorFn", valueSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	valueSelectorFunc := func(item interface{}) interface{} {
		return valueSelectorGenericFunc.Call(item)
	}

	q.ToMapBy(result, keySelectorFunc, valueSelectorFunc)
}
