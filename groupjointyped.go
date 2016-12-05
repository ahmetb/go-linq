package linq

import "reflect"

// GroupJoinT is the typed version of GroupJoin.
//
// NOTE: GroupJoin method has better performance than GroupJoinT
//
// inner: The query to join to the outer query.
//
// outerKeySelectorFn is of a type "func(TOuter) TKey"
//
// innerKeySelectorFn is of a type "func(TInner) TKey"
//
// resultSelectorFn: is of a type "func(TOuter, inners []TInner) TResult"
//
func (q Query) GroupJoinT(inner Query, outerKeySelectorFn interface{}, innerKeySelectorFn interface{}, resultSelectorFn interface{}) Query {

	outerKeySelectorGenericFunc, err := newGenericFunc(
		"GroupJoinT", "outerKeySelectorFn", outerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	outerKeySelectorFunc := func(item interface{}) interface{} {
		return outerKeySelectorGenericFunc.Call(item)
	}

	innerKeySelectorFuncGenericFunc, err := newGenericFunc(
		"GroupJoinT", "innerKeySelectorFn", innerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	innerKeySelectorFunc := func(item interface{}) interface{} {
		return innerKeySelectorFuncGenericFunc.Call(item)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"GroupJoinT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, inners []interface{}) interface{} {
		innerSliceType := reflect.MakeSlice(resultSelectorGenericFunc.Cache.TypesIn[1], 0, 0)
		innersSlicePointer := reflect.New(innerSliceType.Type())
		From(inners).ToSlice(innersSlicePointer.Interface())
		innersTyped := reflect.Indirect(innersSlicePointer).Interface()
		return resultSelectorGenericFunc.Call(outer, innersTyped)
	}

	return q.GroupJoin(inner, outerKeySelectorFunc, innerKeySelectorFunc, resultSelectorFunc)
}
